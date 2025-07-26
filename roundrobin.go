package botf

import (
	"errors"
	"net/url"
	"sync"
)

type RoundRobinBalancer struct {
	conns []Connection
	sync.Mutex

	// to keep track of current connection index as we move to the next list
	idx int
}

func NewRoundRobinBalancer(urls ...string) (*RoundRobinBalancer, error) {
	b := &RoundRobinBalancer{
		conns: []Connection{},
	}

	if len(urls) == 0 {
		return nil, errors.New("need at least 1 url to create balancer instance")
	}

	for _, u := range urls {
		parsedUrl, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		b.conns = append(b.conns, NewConnection(parsedUrl))
	}

	return b, nil
}

func (r *RoundRobinBalancer) Connections() []Connection {
	return r.conns
}

func (r *RoundRobinBalancer) Get() Connection {
	r.Lock()
	defer r.Unlock()

	candidate := r.conns[r.idx]
	r.idx = (r.idx + 1) % len(r.conns)

	if candidate == nil {
		return nil
	}

	return candidate
}
