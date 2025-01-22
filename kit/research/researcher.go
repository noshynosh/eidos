package research

import "context"

type Researcher struct{}

func (r *Researcher) Description() string {
	return "Research a topic"
}

func (r *Researcher) Manual() string {
	return `
		You will be given a topic to research. You will need to research the topic and return a summary of the research.

		Do not include quotes in your response.
	`
}

func (r *Researcher) Use(ctx context.Context, req any) error {
	return nil
}
