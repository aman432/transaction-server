// Package routes provides route handlers for account-related endpoints.
package routes

import (
	"github.com/gin-gonic/gin"
	"transaction-server/internal/account"
	"transaction-server/internal/dto"
)

// Accounts represents the route handler for account-related endpoints.
type Accounts struct {
	server account.IServer
}

// NewAccountsRoute creates a new Accounts route handler.
func NewAccountsRoute(server account.IServer) *Accounts {
	return &Accounts{
		server: server,
	}
}

// Create handles the creation of a new account.
// swagger:operation POST /accounts Create
//
// Creates a new account.
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - in: body
//     name: body
//     description: The account to be created.
//     required: true
//     schema:
//     "$ref": "#/definitions/CreateAccountRequest"
//
// responses:
//
//	'200':
//	  description: Account created successfully.
//	  schema:
//	    "$ref": "#/definitions/CreateAccountResponse"
//	'400':
//	  description: Bad request. Error response returned.
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
//	'500':
//	  description: Internal server error. Error response returned.
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
func (a *Accounts) Create(ctx *gin.Context) {
	var createRequest dto.CreateAccountRequest

	// Call BindJSON to bind the received JSON to request.
	if err := ctx.BindJSON(&createRequest); err != nil {
		ctx.JSON(400, dto.ErrorResponse{Message: err.Error()})
		return
	}

	response := a.server.Create(ctx, &createRequest)
	SendResponse(ctx, response)
}

// Get retrieves an account by ID.
// swagger:operation GET /accounts/{accountId} Get
//
// Retrieves an account by ID.
// ---
// produces:
// - application/json
// parameters:
//   - name: accountId
//     in: path
//     description: The ID of the account to retrieve.
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Account retrieved successfully.
//	  schema:
//	    "$ref": "#/definitions/GetAccountResponse"
//	'400':
//	  description: Bad request. Error response returned.
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
//	'500':
//	  description: Internal server error. Error response returned.
//	  schema:
//	    "$ref": "#/definitions/ErrorResponse"
func (a *Accounts) Get(ctx *gin.Context) {
	id := ctx.Param("accountId")
	response := a.server.Get(ctx, id)
	SendResponse(ctx, response)
}
