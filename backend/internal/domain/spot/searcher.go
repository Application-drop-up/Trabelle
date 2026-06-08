package spot

import "context"

type Searcher interface {
	Search(ctx context.Context, query string) ([]*Spot, error)
}
