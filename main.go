package main

import (
	genimg "MidAI/cap/image"
	gentext "MidAI/cap/text"
	"fmt"
)

func main() {
	var choice int
	// Start the conversation
	fmt.Printf("select your Generative AI type:\n1. Text\n2. Image\n")
	fmt.Scan(&choice)
	if choice == 1 {
		gentext.Prompt()
	} else if choice == 2 {
		genimg.Prompt()
	}
}
