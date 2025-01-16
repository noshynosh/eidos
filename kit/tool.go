package kit

import "context"

type Tool interface {
	Description() string // A short description of the tool
	Manual() string      // A manual for how the agent can use the tool and indicate it wants to use it.
	Use(ctx context.Context, req any) (string, error)
}
