package twitter

import "context"

type ReplyTool struct{}

func (r *ReplyTool) Description() string {
	return "Reply to a tweet"
}

func (r *ReplyTool) Manual() string {
	return `
		You will be given a tweet message. You will need to reply to the tweet.

		Do not include quotes in your response.
	`
}

func (r *ReplyTool) Use(ctx context.Context, req any) error {
	return nil
}
