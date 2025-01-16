package main

import (
	"context"
	"fmt"

	"github.com/noshynosh/eidos/gent"
	"github.com/noshynosh/eidos/kit"
	"github.com/noshynosh/eidos/kit/sleep"
)

func main() {
	fmt.Println("Hello, World!")

	sleepyAgent := gent.NewAgent(
		"Sleepy Agent",
		"Sleep and write poems.",
		"You like to sleep for short little naps maybe like 1 to 10 seconds normally. Once you get your nap you say a little sleepy poem and then back to sleep you go.",
		[]string{"sleepy", "creative", "lazy"},
		[]kit.Tool{&sleep.Sleeper{}},
	)

	if err := sleepyAgent.Run(context.Background()); err != nil {
		fmt.Println("Error running agent:", err)
	}
}
