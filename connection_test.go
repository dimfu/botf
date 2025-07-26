package botf

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestConnection(t *testing.T) {
	var count int
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	url, _ := url.Parse(s.URL)
	conn := NewConnection(url)

	if conn == nil {
		t.Error("expected conn instance to be created")
	}

	isBroken := conn.IsBroken()

	if isBroken {
		t.Error("expected connection to not be broken")
	}

	if count != 1 {
		t.Error("expected count to be 1 because of initial check")
	}
}

func TestInternalServerError(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer s.Close()

	url, _ := url.Parse(s.URL)
	conn := NewConnection(url)

	if conn == nil {
		t.Error("expected conn instance to be created")
	}

	isBroken := conn.IsBroken()

	if !isBroken {
		t.Error("expected connection to be broken")
	}
}

func TestBrokenUrl(t *testing.T) {
	url, _ := url.Parse("localhost:6969")
	conn := NewConnection(url)

	if conn == nil {
		t.Error("expected conn instance to be created")
	}

	isBroken := conn.IsBroken()

	if !isBroken {
		t.Error("expected connection to be broken")
	}
}
