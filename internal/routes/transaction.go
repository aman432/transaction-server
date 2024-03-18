package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"transaction-server/internal/dto"
	"transaction-server/internal/transaction"
)

// Transactions represents the route handler for transaction-related endpoints.
type Transactions struct {
	server transaction.IServer
}

// NewTransactionsRoute creates a new Transactions route handler.
func NewTransactionsRoute(server transaction.IServer) *Transactions {
	return &Transactions{
		server: server,
	}
}

// Create handles the creation of a new transaction.
// swagger:operation POST /transactions Create
//
// Creates a new transaction.
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - in: body
//     name: body
//     description: The transaction to be created.
//     required: true
//     schema:
//     "$ref": "#/definitions/CreateTransactionRequest"
//
// responses:
//
//	'200':
//	  description: Transaction created successfully.
//	  schema:
//	    "$ref": "#/definitions/CreateTransactionResponse"
//	'4xx':
//	  description: Client Error
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
//	'5xx':
//	  description: Server Error
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
func (a *Transactions) Create(ctx *gin.Context) {
	var createRequest dto.CreateTransactionRequest

	if err := ctx.BindJSON(&createRequest); err != nil {
		// Handle client error
		ctx.JSON(http.StatusBadRequest, dto.GetErrorResponse("BadRequest", "Invalid request payload"))
		return
	}
	response := a.server.Create(ctx, &createRequest)
	SendResponse(ctx, response)
}

// Get retrieves a transaction by ID.
// swagger:operation GET /transactions/{transactionId} Get
//
// Retrieves a transaction by ID.
// ---
// produces:
// - application/json
// parameters:
//   - name: transactionId
//     in: path
//     description: The ID of the transaction to retrieve.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Transaction retrieved successfully.
//	  schema:
//	    "$ref": "#/definitions/GetTransactionResponse"
//	'4xx':
//	  description: Client Error
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
//	'5xx':
//	  description: Server Error
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
func (a *Transactions) Get(ctx *gin.Context) {
	id := ctx.Param("transactionId")
	response := a.server.Get(ctx, id)
	SendResponse(ctx, response)
}

// List retrieves a list of transactions.
// swagger:operation POST /transactions/list List
//
// Retrieves a list of transactions.
// ---
// produces:
// - application/json
// parameters:
//   - in: body
//     name: body
//     description: The transaction to be created.
//     required: true
//     schema:
//     "$ref": "#/definitions/ListTransactionRequest"
//
// responses:
//
//	'200':
//	  description: Transactions retrieved successfully.
//	  schema:
//	    "$ref": "#/definitions/ListTransactionResponse"
//	'4xx':
//	  description: Client Error
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
//	'5xx':
//	  description: Server Error
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
func (a *Transactions) List(ctx *gin.Context) {
	var listRequest dto.ListTransactionRequest

	if err := ctx.BindJSON(&listRequest); err != nil {
		// Handle client error
		ctx.JSON(http.StatusBadRequest, dto.GetErrorResponse("BadRequest", "Invalid request payload"))
		return
	}
	response := a.server.List(ctx, &listRequest)
	SendResponse(ctx, response)
}
