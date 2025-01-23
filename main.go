package main

import (
	"context"
	"fmt"
	"log"

	"github.com/noshynosh/eidos/gent"
)

func main() {
	agent, err := gent.NewAgent()
	if err != nil {
		log.Fatalf("Error creating agent: %v", err)
	}

	if err := agent.Run(context.Background()); err != nil {
		fmt.Println("Error running agent:", err)
	}
}
