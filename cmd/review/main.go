package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mauriceLC92/review"
)

const (
	HELP_REVIEW = "review - Initiate a new review or check when the next is due."
	HELP_LIST   = "list - List previous reviews which have been filled out."
	HELP_HELP   = "help - Display the CLI commands available to you."
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "review":
			runReview()
		case "list":
			listReviews()
		case "help":
			generateHelpMenu()
		default:
			fmt.Println("command not recognised")
		}
	} else {
		fmt.Println("no commands given")
	}

}

func generateHelpMenu() {
	fmt.Fprint(os.Stdout, strings.Join([]string{HELP_REVIEW, HELP_LIST, HELP_HELP}, "\n"))
}

func runReview() {
	store, err := review.OpenJSONStore("reviews.json")
	if err != nil {
		fmt.Println("error opening file of reviews", err)
	}

	r, ok := store.GetLatestReview()
	if !ok {
		fmt.Println("You have not done a review yet! Let's get you started")
	}
	if (r.Due() && !r.Answered()) || (r.CreatedToday() && !r.Answered()) {
		r.Review(os.Stdout, os.Stdin)
		review.SaveTo(r, "reviews.json")
		fmt.Printf("Thanks for the review! See you on %v for the next one!", r.NextDueDate())
	} else {
		fmt.Printf("You review is not due until %v. See you then! \n", r.NextDueDate())
	}
}

func listReviews() {
	store, err := review.OpenJSONStore("reviews.json")
	if err != nil {
		fmt.Println("error opening file of reviews", err)
	}
	reviews := store.GetAll()
	for _, r := range reviews {
		fmt.Println(r)
	}
}
