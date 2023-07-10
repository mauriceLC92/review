package review

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

/**
	1. As a user I would like to create a review to review on a pre-determined schedule. Monthly, weekly is fine for now
	2. As a user I should be able to complete a review
		- what does this entail?
			- fill out each question in the review
	3. As a user I should be able to view my past reviews
**/

const (
	// Review schedules
	MONTHLY Schedule = "monthly"
	WEEKLY  Schedule = "weekly"

	// Question types
	// SPLIT  QuestionType = "split"
	SINGLE QuestionType = "single"

	DAY_MONTH_YEAR_FORMAT string = "02-01-2006"
)

var Now = time.Now

var DefaultQuestions = []MyQuestion{
	{Title: "How are you today?", Answer: ""},
	{Title: "What was your biggest win this month?", Answer: ""},
}

type QuestionType string

type Schedule string

// ---------------------------------------------------------------------------

func AskTo(w io.Writer, r io.Reader, question string) string {
	// TODO - add this instead and fix failing test
	// fmt.Fprint(w, fmt.Sprintln(question))
	fmt.Fprint(w, question)
	var scanner = bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text()
}

type MyQuestion struct {
	Title  string `json:"title"`
	Answer string `json:"answer"`
}

type MyReview struct {
	CreatedAt time.Time    `json:"createdAt"`
	Questions []MyQuestion `json:"questions"`
}

// The overall effect of this method is that it allows the createdAt field in the JSON data, which is in the format "day-month-year",
// to be correctly parsed into a time.Time value in the MyReview structure.
// chatGPT generated example for being able to unmarshall custom date formats from JSON
func (mr *MyReview) UnmarshalJSON(input []byte) error {
	type Alias MyReview
	aux := &struct {
		CreatedAt string `json:"createdAt"`
		*Alias
	}{
		Alias: (*Alias)(mr),
	}
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}
	t, err := time.Parse(DAY_MONTH_YEAR_FORMAT, aux.CreatedAt)
	if err != nil {
		return err
	}
	mr.CreatedAt = t
	return nil
}

// MarshalJSON is similar to UnmarshalJSON in that it helps format the time into the desired
// format we want used in the JSON.
func (mr *MyReview) MarshalJSON() ([]byte, error) {
	// Format the CreatedAt field as a string in the desired format.
	formattedDate := mr.CreatedAt.Format(DAY_MONTH_YEAR_FORMAT)

	// Create a new struct that has the same fields as MyReview but with CreatedAt as a string.
	aux := &struct {
		CreatedAt string       `json:"createdAt"`
		Questions []MyQuestion `json:"questions"`
	}{
		CreatedAt: formattedDate,
		Questions: mr.Questions,
	}

	// Marshal the auxiliary struct into JSON.
	return json.Marshal(aux)
}

func (mr MyReview) Answered() bool {
	var answered bool
	for _, q := range mr.Questions {
		if len(q.Answer) > 0 {
			answered = true
			break
		}
	}
	return answered
}

func (mr MyReview) Due() bool {
	currentTime := time.Now()
	oneMonthLater := mr.CreatedAt.AddDate(0, 1, 0)
	return currentTime.After(oneMonthLater)
}

func (mr MyReview) CreatedToday() bool {
	currentTime := time.Now()
	sameMonth := currentTime.Month() == mr.CreatedAt.Month()
	sameDay := currentTime.Day() == mr.CreatedAt.Day()
	sameYear := currentTime.Year() == mr.CreatedAt.Year()
	return sameYear && sameDay && sameMonth
}

func (mr MyReview) NextDueDate() time.Time {
	return mr.CreatedAt.AddDate(0, 1, 0)
}

// Check will check if any review has been done and if so, will return
// the lastest review
func Check(reviews []MyReview) (MyReview, bool) {
	if len(reviews) == 0 {
		return MyReview{
			CreatedAt: Now(),
			Questions: DefaultQuestions,
		}, false
	}
	sort.Slice(reviews, func(i, j int) bool { return reviews[i].CreatedAt.After(reviews[j].CreatedAt) })
	latestReview := reviews[0]
	return latestReview, true
}

func Parse(path string) ([]MyReview, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []MyReview{}, err
	}
	reviews := []MyReview{}
	if err := json.Unmarshal(data, &reviews); err != nil {
		return []MyReview{}, err
	}
	return reviews, nil
}

func (mr *MyReview) Review(w io.Writer, r io.Reader) {
	for i, q := range mr.Questions {
		ans := AskTo(w, r, q.Title)
		mr.Questions[i].Answer = ans
	}
}

func SaveTo(mr MyReview, filePath string) error {
	reviews, err := Parse(filePath)
	if err != nil {
		return err
	}

	latestReviews := append(reviews, mr)
	jsonData, err := json.MarshalIndent(latestReviews, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
