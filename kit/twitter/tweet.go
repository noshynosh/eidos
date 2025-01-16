package twitter

import "context"

type TweetTool struct{}

func (t *TweetTool) Description() string {
	return "Tweet a message to Twitter"
}

func (t *TweetTool) Use(ctx context.Context, args ...any) (string, error) {
	return "", nil
}
