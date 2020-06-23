package dehashed

import (
	"context"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	c := New(os.Getenv("email"), os.Getenv("apikey"))
	results, err := c.Search(context.Background(), &SearchParams{
		Query: "email:test@test.com",
	})
	if err != nil {
		t.Error(err)
		return
	}

	if len(results.Entries) == 0 {
		t.Errorf("Not enough results")
	}
}
