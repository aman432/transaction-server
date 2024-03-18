package db

import (
	"context"
	"errors"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

const updatedAtField = "updated_at"

// Repoer represents Repo family.
type Repoer interface {
	FindByID(ctx context.Context, receiver IModel, id string) error
	FindMany(ctx context.Context, models interface{}, req FindManyRequester) error
	FindManyWithFilters(ctx context.Context, models interface{}, req FindManyWithFiltersRequester) error
	FindByKey(ctx context.Context, model interface{}, key string, value string) error
	FindByConditions(ctx context.Context, model interface{}, conditions []clause.Expression) error
	Create(ctx context.Context, receiver IModel) error
	CreateWithAssociations(ctx context.Context, receiver IModel) error
	SaveWithAssociations(ctx context.Context, receiver IModel) error
	Update(ctx context.Context, receiver IModel, selectiveList ...string) error
	Delete(ctx context.Context, receiver IModel) error
	AssociationFind(ctx context.Context, model interface{}, association string, a interface{}) error
	Transaction(ctx context.Context, fc func(ctx context.Context) error) error
	DBInstance(ctx context.Context) *gorm.DB
}

type Repo struct {
	Db *DB
}

func NewRepo(db *DB) Repoer {
	return &Repo{Db: db}
}

// DBInstance returns gorm instance.
// If replicas are specified, for Query, Row callback, will use replicas, unless Write mode specified.
// For Raw callback, statements are considered read-only and will use replicas if the SQL starts with SELECT.
func (r *Repo) DBInstance(ctx context.Context) *gorm.DB {
	return r.Db.Instance(ctx)
}

// FindByID fetches the record which matches the ID provided from the entity defined by receiver
// and the result will be loaded into receiver
func (r *Repo) FindByID(ctx context.Context, receiver IModel, id string) error {
	q := r.DBInstance(ctx).Where("id = ?", id).First(receiver)

	return q.Error
}

// Create inserts a new record in the entity defined by the receiver
// all data filled in the receiver will inserted
func (r *Repo) Create(ctx context.Context, receiver IModel) error {
	if err := receiver.BeforeCreate(r.DBInstance(ctx)); err != nil {
		return err
	}

	if err := receiver.SetDefaults(); err != nil {
		return err
	}

	if err := receiver.Validate(); err != nil {
		return err
	}

	q := r.DBInstance(ctx).Create(receiver)

	return q.Error
}

// CreateInBatches insert the value in batches into database
func (r *Repo) CreateInBatches(ctx context.Context, receivers interface{}, batchSize int) error {
	q := r.DBInstance(ctx).CreateInBatches(receivers, batchSize)

	return q.Error
}

// Update will update the given receiver model with respect to primary key / id available in it.
// If selective list is non empty, only those fields which are present in the list will be updated.
// Note: When using selectiveList `updated_at` field need not be passed in the list.
func (r *Repo) Update(ctx context.Context, receiver IModel, selectiveList ...string) error {
	if len(selectiveList) > 0 {
		selectiveList = append(selectiveList, updatedAtField)
	}
	return r.updateSelective(ctx, receiver, selectiveList...)
}

// updateSelective will update the given receiver model with respect to primary key / id available in it.
// If selective list is non empty, only those fields which are present in the list will be updated.
// Note: When using selectiveList `updated_at` field also needs to be explicitly passed in the selectiveList.
func (r *Repo) updateSelective(ctx context.Context, receiver IModel, selectiveList ...string) error {
	q := r.DBInstance(ctx).Model(receiver)

	if len(selectiveList) > 0 {
		q = q.Select(selectiveList)
	}

	q = q.Updates(receiver)

	return q.Error
}

// Delete deletes the given model
// Soft or hard delete of model depends on the models implementation
// if the model composites SoftDeletableModel then it'll be soft deleted
func (r *Repo) Delete(ctx context.Context, receiver IModel) error {
	q := r.DBInstance(ctx).Delete(receiver)

	return q.Error
}

func (r *Repo) CreateWithAssociations(ctx context.Context, receiver IModel) error {
	q := r.DBInstance(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(receiver)

	return q.Error
}

func (r *Repo) SaveWithAssociations(ctx context.Context, receiver IModel) error {
	q := r.DBInstance(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(receiver)

	return q.Error
}

// FindMany builds query per request and loads multiple into models.
func (r *Repo) FindMany(ctx context.Context, models interface{}, req FindManyRequester) error {
	q := r.ApplyListRequest(req, r.DBInstance(ctx))

	q = q.Order("created_at DESC").Order("name desc")

	return q.Find(models).Error
}

// FindManyWithFilters builds query per request with filters and loads multiple into models.
func (r *Repo) FindManyWithFilters(ctx context.Context, models interface{}, req FindManyWithFiltersRequester) error {
	q := r.ApplyListRequest(req, r.DBInstance(ctx))
	q = q.Clauses(req.GetConditions()...).Order("created_at DESC")

	return q.Find(models).Error
}

// FindByKey builds query By the key and the value.
func (r *Repo) FindByKey(ctx context.Context, model interface{}, key string, value string) error {
	if key == "" || value == "" {
		return errors.New("key/value must not be empty")
	}

	q := r.DBInstance(ctx).Clauses([]clause.Expression{
		clause.Eq{
			Column: clause.Column{
				Name: key,
			},
			Value: value,
		},
	}...).First(model)

	return q.Error
}

func (r *Repo) FindByConditions(ctx context.Context, model interface{}, conditions []clause.Expression) error {
	if len(conditions) == 0 {
		return errors.New("key/value must not be empty")
	}

	db := r.DBInstance(ctx)
	db = db.Clauses(conditions...)

	return db.First(model).Error
}

func (r *Repo) AssociationFind(ctx context.Context, model interface{}, association string, a interface{}) error {
	if err := r.DBInstance(ctx).Model(model).Association(association).Find(a); err != nil {
		return err
	}
	return nil
}

func (r *Repo) AssociationDelete(ctx context.Context, model interface{}, association string, a ...interface{}) error {
	if err := r.DBInstance(ctx).Model(model).Association(association).Delete(a...); err != nil {
		return err
	}
	return nil
}

// Transaction will manage the execution inside a transactions
// adds the txn db in the context for downstream use case
func (r *Repo) Transaction(ctx context.Context, fc func(ctx context.Context) error) error {
	var err = r.DBInstance(ctx).Transaction(func(tx *gorm.DB) error {

		// This will ensure that when db.Instance(context) we return the txn on the context
		// & all repo queries are done on this txn. Refer usage in test.
		if err := fc(context.WithValue(ctx, ContextKeyDatabase, tx)); err != nil {
			return err
		}

		return tx.Error
	})

	return err
}

// ApplyListRequest applyListRequest applies request to given query builder and returns new query builder.
func (r *Repo) ApplyListRequest(req FindManyRequester, q *gorm.DB) *gorm.DB {
	if req.GetLimit() != 0 {
		q = q.Limit(int(req.GetLimit()))
	} else {
		q = q.Limit(10)
	}
	if req.GetOffset() != 0 {
		q = q.Offset(int(req.GetOffset()))
	} else {
		q = q.Offset(0)
	}

	return q
}

// FindManyRequester is the interface that wraps basic request attribute getters.
type FindManyRequester interface {
	GetLimit() uint32
	GetOffset() uint32
}

// FindManyRequest is a request to find many of a model type.
type FindManyRequest struct {
	Limit  uint32
	Offset uint32
}

// GetLimit returns limit.
func (r FindManyRequest) GetLimit() uint32 {
	return r.Limit
}

// GetOffset returns offset.
func (r FindManyRequest) GetOffset() uint32 {
	return r.Offset
}

// FindManyWithFiltersRequester is the interface that wraps basic request attribute getters.
type FindManyWithFiltersRequester interface {
	FindManyRequester
	GetConditions() []clause.Expression
}

// FindManyWithConditionsRequest is a request to find many of a model type.
type FindManyWithConditionsRequest struct {
	FindManyRequest
	Conditions []clause.Expression
}

// GetConditions returns list of additional filters on entity.
func (r FindManyWithConditionsRequest) GetConditions() []clause.Expression {
	return r.Conditions
}
