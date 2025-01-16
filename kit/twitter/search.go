package twitter

import "context"

type SearchTool struct{}

func (s *SearchTool) Description() string {
	return "Search for tweets"
}

func (s *SearchTool) Use(ctx context.Context, args ...any) (string, error) {
	return "", nil
}
