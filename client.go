package dehashed

import (
	"time"

	"github.com/juju/ratelimit"
)

const (
	baseURL = "https://api.dehashed.com"
	// Reload rate limit of dehashed
	rateLimitReloadInterval = time.Millisecond * 250
	rateLimitReloadAmount   = 5
	rateLimitMaxAmount      = 5
)

// Client is a dehashed client
type Client struct {
	email  string
	apiKey string

	rateLimitBucket *ratelimit.Bucket
}

// New Creates a new dehashed client
func New(email, apiKey string) *Client {
	bucket := ratelimit.NewBucketWithQuantum(rateLimitReloadInterval, rateLimitMaxAmount, rateLimitReloadAmount)

	return &Client{
		email:           email,
		apiKey:          apiKey,
		rateLimitBucket: bucket,
	}
}
