package botf

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewRoundRobinConnections(t *testing.T) {
	balancer, err := NewRoundRobinBalancer("localhost:8080", "localhost:8081")
	if err != nil {
		t.Error(err)
	}
	if len(balancer.Connections()) != 2 {
		t.Errorf("expected 2 connections got %d", len(balancer.Connections()))
	}
}

func TestBalancer(t *testing.T) {
	visited := []int{}
	s1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visited = append(visited, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer s1.Close()

	s2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		visited = append(visited, 2)
		w.WriteHeader(http.StatusOK)
	}))
	defer s2.Close()

	balancer, err := NewRoundRobinBalancer(s1.URL, s2.URL)
	if err != nil {
		t.Error(err)
	}

	client := NewClient(balancer)
	client.Get(s1.URL)
	client.Get(s1.URL)
	client.Get(s1.URL)
	expected := []int{1, 2, 1, 2, 1}
	if !reflect.DeepEqual(visited, expected) {
		t.Errorf("expecting visited to be %v got %v", expected, visited)
	}
}
