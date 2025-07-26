package botf

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

const DEFAULT_CHECK_INTERVAL = time.Second * 10

type Connection interface {
	URL() *url.URL
	IsBroken() bool
}

type HTTPConnection struct {
	sync.Mutex
	Url                 *url.URL
	healthCheckInterval time.Duration
	healthCheckStop     chan bool
	broken              bool
}

func NewConnection(url *url.URL) *HTTPConnection {
	conn := &HTTPConnection{
		Url:                 url,
		healthCheckInterval: DEFAULT_CHECK_INTERVAL,
		healthCheckStop:     make(chan bool),
	}
	// check on initial creation
	conn.checker()
	go conn.healthCheck()
	return conn
}

func (h *HTTPConnection) healthCheck() {
	ticker := time.NewTicker(h.healthCheckInterval)
	for {
		select {
		case <-ticker.C:
			if err := h.checker(); err != nil {
				return
			}
		case <-h.healthCheckStop:
			return
		}
	}
}

func (h *HTTPConnection) HealthCheckDuration(duration time.Duration) *HTTPConnection {
	h.Lock()
	defer h.Unlock()
	h.healthCheckInterval = duration
	h.healthCheckStop <- true
	h.broken = false
	go h.healthCheck()
	return h
}

func (h *HTTPConnection) checker() error {
	req, err := http.NewRequest(http.MethodGet, h.Url.String(), nil)
	if err != nil {
		h.broken = true
		return err
	}

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err == nil {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			h.broken = false
		} else {
			h.broken = true
		}
	} else {
		h.broken = true
	}
	return nil
}

func (h *HTTPConnection) Close() {
	h.Lock()
	defer h.Unlock()
	h.healthCheckStop <- true
	return
}

func (h *HTTPConnection) IsBroken() bool {
	return h.broken
}

func (h *HTTPConnection) URL() *url.URL {
	return h.Url
}
