package botf

import "net/url"

type Balancer interface {
	Get() *url.URL
	Connections() []url.URL
}
