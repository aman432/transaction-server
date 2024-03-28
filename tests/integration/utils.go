package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"transaction-server/internal/dto"
)

func getTransactionResponse(value []byte) dto.Transaction {
	transaction := new(dto.CreateTransactionResponse)
	_ = json.Unmarshal(value, &transaction)
	return *transaction.Transaction
}

func getAccountsResponse(value []byte) dto.Account {
	account := new(dto.CreateAccountResponse)
	_ = json.Unmarshal(value, &account)
	return *account.Account
}

func marshalJson(value interface{}) []byte {
	jsonValue, _ := json.Marshal(value)
	return jsonValue
}

func makeAPICall(t *testing.T, requestBody []byte, url string, method string) []byte {
	// Create a new HTTP request targeting the API endpoint
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the HTTP request
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	return body
}
