package kibela

import (
	"context"
	"net/http"

	"github.com/kuroponzu/kibela-cli/internal/config"
	"github.com/shurcooL/graphql"
)

// Client is a Kibela GraphQL client.
type Client struct {
	gql    *graphql.Client
	config *config.Config
}

// tokenTransport adds authorization header to requests.
type tokenTransport struct {
	token string
	base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return t.base.RoundTrip(req)
}

// NewClient creates a new Kibela client.
func NewClient(cfg *config.Config) *Client {
	httpClient := &http.Client{
		Transport: &tokenTransport{
			token: cfg.Token,
			base:  http.DefaultTransport,
		},
	}

	gqlClient := graphql.NewClient(cfg.Endpoint(), httpClient)

	return &Client{
		gql:    gqlClient,
		config: cfg,
	}
}

// Query executes a GraphQL query.
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.gql.Query(ctx, q, variables)
}

// Mutate executes a GraphQL mutation.
func (c *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	return c.gql.Mutate(ctx, m, variables)
}
