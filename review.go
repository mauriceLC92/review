package review

import (
	"bufio"
	"encoding/json"
	"errors"
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

type QuestionType string

type Question struct {
	ID          string
	Title       string
	Description string
	Answer      string
	Type        QuestionType
}

type Questions []Question

func (qs Questions) AllComplete() (bool, error) {
	if len(qs) == 0 {
		return false, errors.New("no questions provided")
	}
	var complete = true
	for _, q := range qs {
		if q.Answer == "" {
			complete = false
		}
	}
	return complete, nil
}

type Schedule string

type Review struct {
	ID        string
	CreatedAt int64 // unix timestamp
	UpdatedAt int64 // unix timestamp
	Complete  bool
	Questions Questions
	Schedule
}

// review create
func Create(r Review) Review {
	return Review{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Complete:  false,
		Questions: r.Questions,
		Schedule:  r.Schedule,
	}
}

func (r *Review) Finish() error {
	questions := r.Questions
	complete, err := questions.AllComplete()
	if err != nil {
		return err
	}
	if complete {
		r.Complete = true
	}
	return nil
}

// ---------------------------------------------------------------------------

func AskTo(w io.Writer, r io.Reader, question string) string {
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
		CreatedAt string `json:"createdAt"`
	}{
		CreatedAt: formattedDate,
	}

	// Marshal the auxiliary struct into JSON.
	return json.Marshal(aux)
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
		return MyReview{}, false
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
