package botf

import (
	"errors"
	"net/url"
	"sync"
)

type RoundRobinBalancer struct {
	conns []*url.URL
	sync.Mutex

	// to keep track of current connection index as we move to the next list
	idx int
}

func NewRoundRobinBalancer(urls ...string) (*RoundRobinBalancer, error) {
	b := &RoundRobinBalancer{
		conns: []*url.URL{},
	}

	if len(urls) == 0 {
		return nil, errors.New("need at least 1 url to create balancer instance")
	}

	for _, u := range urls {
		parsedUrl, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		b.conns = append(b.conns, parsedUrl)
	}

	return b, nil
}

// Returns all of the connections
func (r *RoundRobinBalancer) Connections() []url.URL {
	var urls []url.URL
	for _, u := range r.conns {
		urls = append(urls, *u)
	}
	return urls
}

// Returns the connection from the balancer that are selected using round robin arrangement
func (r *RoundRobinBalancer) Get() *url.URL {
	r.Lock()
	defer r.Unlock()

	candidate := r.conns[r.idx]
	r.idx = (r.idx + 1) % len(r.conns)

	if candidate == nil {
		return nil
	}

	return candidate
}
