package handler

import (
	"net/http"
	"testing"
)

func TestHealth_ReturnsOk(t *testing.T) {
	result, err := Health(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
}

func TestGetCategories_ReturnsData(t *testing.T) {
	result, err := GetCategories(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
}

func TestErrBadRequest_StatusCode(t *testing.T) {
	err := errBadRequest("test")
	if err.StatusCode() != http.StatusBadRequest {
		t.Errorf("got %d, want %d", err.StatusCode(), http.StatusBadRequest)
	}
	if err.Error() != "test" {
		t.Errorf("got %q, want %q", err.Error(), "test")
	}
}

func TestErrEntityTooLarge_StatusCode(t *testing.T) {
	err := errEntityTooLarge("too big")
	if err.StatusCode() != http.StatusRequestEntityTooLarge {
		t.Errorf("got %d, want %d", err.StatusCode(), http.StatusRequestEntityTooLarge)
	}
}

func TestErrUnprocessable_StatusCode(t *testing.T) {
	err := errUnprocessable("bad pdf")
	if err.StatusCode() != http.StatusUnprocessableEntity {
		t.Errorf("got %d, want %d", err.StatusCode(), http.StatusUnprocessableEntity)
	}
}
