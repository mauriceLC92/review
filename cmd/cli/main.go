package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mauriceLC92/review"
)

var questions = []string{"What should we do next?", "What is the meaning of life?"}

func main() {

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "ask":
			askQuestions(questions)
		case "time":
			const formatString = "02-01-2006"
			myTime, _ := time.Parse(formatString, "20-06-2023")
			fmt.Println(myTime.UnixMilli())
		default:
			fmt.Println("command not recognised")
		}
	} else {
		fmt.Println("no commands given")
	}

}

// write a CLI that uses the AskTo function and have it ask me some questions
func askQuestions(questions []string) {
	for _, q := range questions {
		ans := review.AskTo(os.Stdout, os.Stdin, fmt.Sprintln(q))
		fmt.Println(ans)
	}
}
