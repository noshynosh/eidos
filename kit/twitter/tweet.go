package twitter

import "context"

type TweetTool struct{}

func (t *TweetTool) Description() string {
	return "TweetTool allows you to tweet a message to Twitter."
}

func (t *TweetTool) Manual() string {
	return `
		To use this tool you must provide a message to tweet. The message will be the entire content of your response.
		Anything you write will be tweeted! Make sure to only include the message you want to tweet.
	`
}

func (t *TweetTool) Use(ctx context.Context, req any) (string, error) {
	return "", nil
}
