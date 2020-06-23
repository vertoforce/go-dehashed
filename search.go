package dehashed

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SearchResponse is the response after sending a search api call
type SearchResponse struct {
	Balance int // API Balance remaining
	Entries []struct {
		ID             string
		Email          string
		IPAddress      string `json:"ip_address"`
		Username       string
		Password       string
		HashedPassword string `json:"hashed_password"`
		HashType       string `json:"hash_type"`
		Name           string
		VIN            string
		Address        string
		Phone          string
		DatabaseName   string `json:"database_name"`
	}
	Success bool
	Took    string
	Total   int // Total results
}

// SearchParams are the parameters you can send for a search
//
// Note that the dehashed api requests:
//
// "DO NOT USE MORE THAN ONE OPTIONS IN THE SAME QUERY, IT WILL SUBSTANTIALLY SLOW DOWN THE SEARCH AND RESULT IN A BAN IF DONE ON PURPOSE."
//
type SearchParams struct {
	// DeHashed is a powerful search tool, we support true regex and wildcard, whilst allowing the end-user to search billions records within a few Micro-Seconds.
	// If you don't care to get fancy, simply wrap your search string with "" (quotation marks) to find exact results that match your query (Example: "test"). It will automatically search all fields for the string. Otherwise, if you want to get fancy and use DeHashed's true power, you can explore this section further.
	//
	// See https://www.dehashed.com/docs for examples of queries
	Query string

	// Page of results to fetch.  Default is 1.
	// Results are limited to 5000 results per query
	Page int
}

// Search performs a dehashed search.
func (c *Client) Search(ctx context.Context, params *SearchParams) (*SearchResponse, error) {
	urlValues := url.Values{}
	urlValues.Add("query", params.Query)
	if params.Page != 0 {
		urlValues.Add("page", fmt.Sprintf("%d", params.Page))
	}
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/search?%s", baseURL, urlValues.Encode()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.email, c.apiKey)))))

	// Wait on rate limit
	c.rateLimitBucket.Wait(1)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 400:
		// Sleep the rate limit reload interval and try again
		time.Sleep(rateLimitReloadInterval + time.Millisecond*5)
		return c.Search(ctx, params)
	case 401:
		return nil, fmt.Errorf("invalid API credentials")
	case 404:
		return nil, fmt.Errorf("404 not found")
	case 302:
		return nil, fmt.Errorf("invalid/missing query")
	case 200:
	default:
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	io.Copy(os.Stdout, resp.Body)

	// Decode
	ret := &SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
