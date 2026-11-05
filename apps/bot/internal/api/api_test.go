package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSendsTheAdminToken(t *testing.T) {
	var gotToken, gotStatus string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("X-Admin-Token")
		gotStatus = r.URL.Query().Get("status")
		_, _ = w.Write([]byte(`{"accounts":[{"id":"1","email":"a@example.com","name":"Ana","status":"pending"}]}`))
	}))
	t.Cleanup(server.Close)

	accounts, err := New(server.URL, "secret").List(context.Background(), "pending", "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if gotToken != "secret" {
		t.Errorf("X-Admin-Token = %q, want secret", gotToken)
	}
	if gotStatus != "pending" {
		t.Errorf("status query = %q, want pending", gotStatus)
	}
	if len(accounts) != 1 || accounts[0].Email != "a@example.com" {
		t.Errorf("List() = %v, want one pending account", accounts)
	}
}

func TestSetStatusSendsTheStatus(t *testing.T) {
	var gotBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buffer := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buffer)
		gotBody = string(buffer)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	if err := New(server.URL, "secret").SetStatus(context.Background(), "id-1", "active"); err != nil {
		t.Fatalf("SetStatus() error = %v", err)
	}
	if gotBody != `{"status":"active"}` {
		t.Errorf("body = %s, want the status payload", gotBody)
	}
}

func TestNonSuccessBecomesAStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	_, err := New(server.URL, "secret").List(context.Background(), "", "x@example.com")

	var status StatusError
	if !errors.As(err, &status) || status.Code != http.StatusNotFound {
		t.Fatalf("List() error = %v, want StatusError 404", err)
	}
}
