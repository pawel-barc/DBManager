package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)
func init() {
	os.Setenv("JWT_SECRET", "test_secret")
}


func TestResister_InvalidJson(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("{invalid json")),)
	w := httptest.NewRecorder()
	Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("attendu 400, recu %d", w.Code)
	}
}


func TestRegister_InvaliPassword(t *testing.T) {
	body := []byte(`{"username": "testuser", "email": "test@email.gov", "password": "wrongpassword" }`)

	req :=httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	Register (w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Attendu 400, recu %d", w.Code)
	}
}