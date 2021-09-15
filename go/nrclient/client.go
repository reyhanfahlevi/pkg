package nrclient

import (
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

const (
	graphqlURL = `https://api.newrelic.com/graphql`

	getListOfUserUnderAccount = `https://user-management.service.newrelic.com/accounts/%v/users`

	createUserUnderAccount = `https://rpm.newrelic.com/user_management/accounts/%v/users/new`
)

// Option nr client option
type Option struct {
	NRLoginCookies string
	APIKey         string
	Timeout        time.Duration
}

// NRClient client
type NRClient struct {
	c            *httpclient.Client
	loginCookies string
	apiKey       string
}

// New instantiate nr client
func New(opt Option) *NRClient {
	if opt.Timeout == 0 {
		opt.Timeout = time.Second * 10
	}

	client := httpclient.NewClient(httpclient.WithHTTPTimeout(opt.Timeout))

	return &NRClient{
		loginCookies: opt.NRLoginCookies,
		apiKey:       opt.APIKey,
		c:            client,
	}
}
