package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mauriceLC92/review"
)

var questions = []string{"What should we do next?", "What is the meaning of life?"}

func main() {

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "ask":
			askQuestions(questions)
		case "review":
			runReview()
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

func runReview() {
	reviews, err := review.Parse("reviews.json")
	if err != nil {
		fmt.Println("error opening file of reviews", err)
	}

	r, ok := review.CheckV2(reviews)
	if !ok {
		fmt.Println("You have not done a review yet! Let's get you started")
		// Call function to ask the initial questions
	}
	if r.Due() {
		fmt.Println("is due :)")
		// r.Review()
	} else {
		fmt.Printf("You review is not due until %v. See you then!", r.NextDueDate())
	}

	// check if the review is due

}
