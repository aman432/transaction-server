package routes

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"transaction-server/internal/common"
	"transaction-server/internal/registry"
)

func RegisterRoutes(ctx context.Context) *gin.Engine {

	apiRegistry := registry.NewRegistry(ctx)

	accountsRoute := NewAccountsRoute(apiRegistry.GetAccountsServer())
	transactionsRoute := NewTransactionsRoute(apiRegistry.GetTransactionsServer())

	router := gin.Default()
	router.GET("/accounts/:accountId", accountsRoute.Get)
	router.POST("/accounts", accountsRoute.Create)

	router.GET("/transactions/:transactionId", transactionsRoute.Get)
	router.POST("/transactions", transactionsRoute.Create)
	router.POST("/transactions/list", transactionsRoute.List)

	return router
}

func SendResponse(ctx *gin.Context, response interface{}) {
	// Convert interface{} to map[string]interface{}
	responseMap := convertToMap(response)
	if responseMap["error"] != nil {
		writeErrorCodeSpecificStatus(ctx, responseMap["error"])
		return
	}
	ctx.IndentedJSON(200, response)
}

func convertToMap(response interface{}) map[string]interface{} {
	responseBytes, _ := json.Marshal(response)
	var result map[string]interface{}
	_ = json.Unmarshal(responseBytes, &result)
	return result
}

func writeErrorCodeSpecificStatus(ctx *gin.Context, errorResponse interface{}) {
	errorMap := convertToMap(errorResponse)
	errorCode := errorMap["code"].(string)
	switch errorCode {
	case common.ErrValidationFailed:
		ctx.IndentedJSON(400, errorResponse)
	case common.ErrNotFoundFailed:
		ctx.IndentedJSON(404, errorResponse)
	default:
		ctx.IndentedJSON(500, errorResponse)
	}
}
