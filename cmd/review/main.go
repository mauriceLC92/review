package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mauriceLC92/review"
)

func main() {

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "review":
			runReview()
		case "list":
			listReviews()
		default:
			fmt.Println("command not recognised")
		}
	} else {
		fmt.Println("no commands given")
	}

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
