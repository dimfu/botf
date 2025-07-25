package botf

import "net/http"

type Transport struct {
	base     http.RoundTripper
	balancer Balancer
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	conn := t.balancer.Get()

	clonedReq := req.Clone(req.Context())
	clonedReq.URL = conn
	clonedReq.Host = conn.Host
	clonedReq.Header.Set("Host", conn.Host)

	res, err := http.DefaultTransport.RoundTrip(clonedReq)

	if err != nil {
		return nil, err
	}

	return res, nil
}
