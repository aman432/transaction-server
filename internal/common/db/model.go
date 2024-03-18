package db

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
	"transaction-server/internal/common/db/datatype"
)

type Model struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type IModel interface {
	TableName() string
	EntityName() string
	GetID() string
	Validate() error
	SetDefaults() error
	BeforeCreate(tx *gorm.DB) error
}

// Validate validates base Model.
func (m *Model) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.By(datatype.IsUUID)),
		validation.Field(&m.CreatedAt, validation.By(datatype.IsTimestamp)),
		validation.Field(&m.UpdatedAt, validation.By(datatype.IsTimestamp)),
	)
}

// GetID gets identifier of entity.
func (m *Model) GetID() string {
	return m.ID
}

// GetCreatedAt gets created time of entity.
func (m *Model) GetCreatedAt() int64 {
	return m.CreatedAt
}

// GetUpdatedAt gets last updated time of entity.
func (m *Model) GetUpdatedAt() int64 {
	return m.UpdatedAt
}

// BeforeCreate sets new id.
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		uuid := uuid.New().String()
		uuid = strings.Replace(uuid, "-", "", -1)
		m.ID = uuid[len(uuid)-14:]
	}
	if m.CreatedAt == 0 {
		m.CreatedAt = time.Now().Unix()
	}
	if m.UpdatedAt == 0 {
		m.UpdatedAt = time.Now().Unix()
	}
	return nil
}
