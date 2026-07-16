package port

import "context"

type WikiPage struct {
	PageID  int
	Title   string
	Slug    string
	Extract string
}

type WikipediaClient interface {
	FetchRandom(ctx context.Context) (*WikiPage, error)
	FetchByPageID(ctx context.Context, pageID int) (*WikiPage, error)
}
