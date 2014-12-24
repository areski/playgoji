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

func TestGetActors(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	GetActors(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Response body received is not ok")
	}
}

func TestNewActor(t *testing.T) {
	actor := Actor{Id: 1, Name: "Shagrath", Age: 50}

	data, err := json.Marshal(actor)
	if err != nil {
		t.Fatalf("Error marshaling actor")
	}

	request, _ := http.NewRequest("POST", "/actors", bytes.NewReader(data))
	response := httptest.NewRecorder()

	NewActor(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("Error")
	}
}

func TestGetActor(t *testing.T) {
	params := map[string]string{
		"id": "1",
	}
	request, _ := http.NewRequest("GET", "/actors/"+params["id"], nil)
	response := httptest.NewRecorder()
	context := web.C{URLParams: params}

	GetActor(context, response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Error")
	}
}

func TestUpdateActor(t *testing.T) {
	params := map[string]string{
		"id": "1",
	}
	actor := Actor{Name: "Dek", Age: 30}

	data, err := json.Marshal(actor)
	if err != nil {
		t.Fatalf("Error marshaling actor")
	}

	request, _ := http.NewRequest("PUT", "/actors/"+params["id"], bytes.NewReader(data))
	response := httptest.NewRecorder()
	context := web.C{URLParams: params}

	UpdateActor(context, response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("Error")
	}
}

func TestDeleteActor(t *testing.T) {
	params := map[string]string{
		"id": "1",
	}
	request, _ := http.NewRequest("DELETE", "/actors/"+params["id"], nil)
	response := httptest.NewRecorder()
	context := web.C{URLParams: params}

	DeleteActor(context, response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("Error")
	}
}
