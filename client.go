package botf

import "net/http"

func NewClient(balancer Balancer) *http.Client {
	return &http.Client{
		Transport: &Transport{balancer: balancer},
	}
}
