package main

import (
	"bytes"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	initDb()
}

func TestGetCustomers(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	GetCustomers(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Response body received is not ok")
	}
}

func TestNewCustomer(t *testing.T) {
	customer := Customer{Id: 1, Name: "Shagrath", Age: 50}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Error marshaling customer")
	}

	request, _ := http.NewRequest("POST", "/customers", bytes.NewReader(data))
	response := httptest.NewRecorder()

	NewCustomer(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("Error")
	}
}

func TestGetCustomer(t *testing.T) {
	params := map[string]string{
		"id": "1",
	}
	request, _ := http.NewRequest("GET", "/customers/"+params["id"], nil)
	response := httptest.NewRecorder()
	context := web.C{URLParams: params}

	GetCustomer(context, response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Error")
	}
}

func TestUpdateCustomer(t *testing.T) {
	params := map[string]string{
		"id": "1",
	}
	customer := Customer{Name: "Dek", Age: 30}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("Error marshaling customer")
	}

	request, _ := http.NewRequest("PUT", "/customers/"+params["id"], bytes.NewReader(data))
	response := httptest.NewRecorder()
	context := web.C{URLParams: params}

	UpdateCustomer(context, response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("Error")
	}
}

func TestDeleteCustomer(t *testing.T) {
	params := map[string]string{
		"id": "1",
	}
	request, _ := http.NewRequest("DELETE", "/customers/"+params["id"], nil)
	response := httptest.NewRecorder()
	context := web.C{URLParams: params}

	DeleteCustomer(context, response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("Error")
	}
}
