package integration

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"transaction-server/internal/dto"
)

func TestCreateTransactionAPI(t *testing.T) {
	// Create an account
	account := []byte(`{
    	"account":{
        	"name":"Account aman",
        	"document_number":"12121231232323"
    	}
	}`)

	accountResponse := makeAPICall(t, account, "http://localhost:9040/accounts", "POST")
	accountId := getAccountsResponse(accountResponse).ID
	require.NotEqual(t, "", accountId)

	// Define the request body
	transaction1 := marshalJson(dto.CreateTransactionRequest{Transaction: &dto.Transaction{
		AccountID:     accountId,
		OperationType: "Withdraw",
		Amount:        50.0,
		EventDate:     "2021-09-01T00:00:00Z",
	}})

	transaction2 := marshalJson(dto.CreateTransactionRequest{Transaction: &dto.Transaction{
		AccountID:     accountId,
		OperationType: "Withdraw",
		Amount:        23.5,
		EventDate:     "2021-09-01T00:00:00Z",
	}})

	transaction3 := marshalJson(dto.CreateTransactionRequest{Transaction: &dto.Transaction{
		AccountID:     accountId,
		OperationType: "Withdraw",
		Amount:        18.7,
		EventDate:     "2021-09-01T00:00:00Z",
	}})

	transaction4 := marshalJson(dto.CreateTransactionRequest{Transaction: &dto.Transaction{
		AccountID:     accountId,
		OperationType: "Purchase_With_Installment",
		Amount:        60.0,
		EventDate:     "2021-09-01T00:00:00Z",
	}})

	transaction5 := marshalJson(dto.CreateTransactionRequest{Transaction: &dto.Transaction{
		AccountID:     accountId,
		OperationType: "Purchase_With_Installment",
		Amount:        100.0,
		EventDate:     "2021-09-01T00:00:00Z",
	}})

	// -ve amount transactions no change in balance
	// create transaction with amount 50
	resp1 := makeAPICall(t, transaction1, "http://localhost:9040/transactions", "POST")
	require.Equal(t, -50.0, getTransactionResponse(resp1).Balance)

	// create transaction with amount 23.5
	resp2 := makeAPICall(t, transaction2, "http://localhost:9040/transactions", "POST")
	require.Equal(t, -23.5, getTransactionResponse(resp2).Balance)

	// create transaction with amount 18.7
	resp3 := makeAPICall(t, transaction3, "http://localhost:9040/transactions", "POST")
	require.Equal(t, -18.7, getTransactionResponse(resp3).Balance)

	// +ve amount transactions
	// Create transaction with amount 60
	resp4 := makeAPICall(t, transaction4, "http://localhost:9040/transactions", "POST")
	require.Equal(t, 0.0, getTransactionResponse(resp4).Balance)

	// Create transaction with amount 100
	resp5 := makeAPICall(t, transaction5, "http://localhost:9040/transactions", "POST")
	require.Equal(t, 67.8, getTransactionResponse(resp5).Balance)

	fmt.Println("All test cases passed!!!!")
}
