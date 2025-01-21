package main

import (
	"context"
	"fmt"

	"github.com/noshynosh/eidos/gent"
)

func main() {
	agent := gent.NewAgent()

	if err := agent.Run(context.Background()); err != nil {
		fmt.Println("Error running agent:", err)
	}
}
