package review

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

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

var DefaultQuestions = []Question{
	{Title: "How are you today?", Answer: ""},
	{Title: "What was your biggest win this month?", Answer: ""},
}

type QuestionType string

type Schedule string

func AskTo(w io.Writer, r io.Reader, question string) string {
	fmt.Fprint(w, fmt.Sprintln(question))
	var scanner = bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text()
}

type Question struct {
	Title  string `json:"title"`
	Answer string `json:"answer"`
}

type Review struct {
	CreatedAt time.Time  `json:"createdAt"`
	Questions []Question `json:"questions"`
}

func (q Question) String() string {
	return fmt.Sprintf("Title: %s\nAnswer: %s\n----------------------------------------------------", q.Title, q.Answer)
}

func (r Review) String() string {
	var questionsStrs []string
	for _, q := range r.Questions {
		questionsStrs = append(questionsStrs, q.String())
	}
	return fmt.Sprintf("Date: %v\nQuestions:\n%s", r.CreatedAt.Format(DAY_MONTH_YEAR_FORMAT), strings.Join(questionsStrs, "\n"))
}

// UnmarshalJSON allows you to unmarshall custom date formats from JSON
// The overall effect of this method is that it allows the createdAt field in the JSON data, which is in the format "day-month-year",
// to be correctly parsed into a time.Time value in the Review structure.
func (mr *Review) UnmarshalJSON(input []byte) error {
	type Alias Review
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
func (mr *Review) MarshalJSON() ([]byte, error) {
	// Format the CreatedAt field as a string in the desired format.
	formattedDate := mr.CreatedAt.Format(DAY_MONTH_YEAR_FORMAT)

	// Create a new struct that has the same fields as Review but with CreatedAt as a string.
	aux := &struct {
		CreatedAt string     `json:"createdAt"`
		Questions []Question `json:"questions"`
	}{
		CreatedAt: formattedDate,
		Questions: mr.Questions,
	}

	// Marshal the auxiliary struct into JSON.
	return json.Marshal(aux)
}

func (mr Review) Answered() bool {
	var answered bool
	for _, q := range mr.Questions {
		if len(q.Answer) > 0 {
			answered = true
			break
		}
	}
	return answered
}

// Due will check if a review is due
func (mr Review) Due() bool {
	currentTime := time.Now()
	oneMonthLater := mr.CreatedAt.AddDate(0, 1, 0)
	return currentTime.After(oneMonthLater)
}

// CreatedToday will check if a review was created today
func (mr Review) CreatedToday() bool {
	currentTime := time.Now()
	sameMonth := currentTime.Month() == mr.CreatedAt.Month()
	sameDay := currentTime.Day() == mr.CreatedAt.Day()
	sameYear := currentTime.Year() == mr.CreatedAt.Year()
	return sameYear && sameDay && sameMonth
}

// NextDueDate will return the next date that a review is due
func (mr Review) NextDueDate() time.Time {
	return mr.CreatedAt.AddDate(0, 1, 0)
}

// Check will check if any review has been done and if so, will return
// the lastest review
func Check(reviews []Review) (Review, bool) {
	if len(reviews) == 0 {
		return Review{
			CreatedAt: Now(),
			Questions: DefaultQuestions,
		}, false
	}
	sort.Slice(reviews, func(i, j int) bool { return reviews[i].CreatedAt.After(reviews[j].CreatedAt) })
	latestReview := reviews[0]
	return latestReview, true
}

// Parse reads data from the filePath provided and attempts to return a slice of reviews if they exist.
// if none exist, an empty slice of reviews is returned instead.
func Parse(path string) ([]Review, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []Review{}, err
	}
	reviews := []Review{}
	if err := json.Unmarshal(data, &reviews); err != nil {
		return []Review{}, err
	}
	return reviews, nil
}

func (mr *Review) Review(w io.Writer, r io.Reader) {
	for i, q := range mr.Questions {
		ans := AskTo(w, r, q.Title)
		mr.Questions[i].Answer = ans
	}
}

func SaveTo(mr Review, filePath string) error {
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

type JSONStore struct {
	filePath string
	reviews  []Review
}

// OpenJSONStore will attempt to open the provided JSON file and return the reviews it contains.
func OpenJSONStore(filePath string) (*JSONStore, error) {
	reviews, err := Parse(filePath)
	if err != nil {
		return nil, err
	}
	return &JSONStore{
		reviews:  reviews,
		filePath: filePath,
	}, nil
}

// GetAll will return all the reviews from a JSONStore
func (js JSONStore) GetAll() []Review {
	return js.reviews
}

// GetLatestReview will return the lastest review from a JSONStore
func (js JSONStore) GetLatestReview() (Review, bool) {
	return Check(js.reviews)
}

// GetLatestReview will save a review to the JSONStore
func (js JSONStore) Save(r Review) error {
	err := SaveTo(r, js.filePath)
	if err != nil {
		return err
	}
	return nil
}
