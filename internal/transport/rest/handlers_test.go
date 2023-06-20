package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const url = "http://localhost:8080/"

func TestHandler_GetOrder_MethodGet(t *testing.T) {

	handler := Handler{}

	request := httptest.NewRequest(http.MethodGet, url, nil)

	w := httptest.NewRecorder()

	handler.GetOrder(w, request)

	response := w.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected: OK\nresult: %v", response.StatusCode)
	}
}
