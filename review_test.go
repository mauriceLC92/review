package review_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mauriceLC92/review"
)

func TestReview(t *testing.T) {
	t.Parallel()
	currentDate := time.Now()
	_ = review.Review{
		ID:        "123",
		CreatedAt: currentDate.UnixMilli(),
		UpdatedAt: currentDate.UnixMilli(),
		Questions: []review.Question{
			{
				Title:       "What were my biggest wins?",
				Description: "A questio to get you thinking",
			},
		},
		Schedule: review.MONTHLY,
	}
}

func TestQuestion(t *testing.T) {
	t.Parallel()
	_ = review.Question{
		ID:          "321",
		Title:       "What were my biggest wins?",
		Description: "Some description goes here.",
		Answer:      "",
		Type:        review.SINGLE,
	}
}

func TestCreateANewReview(t *testing.T) {
	t.Parallel()

	want := review.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  false,
		Questions: review.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "",
				Type:        review.SINGLE,
			},
		},
		Schedule: review.MONTHLY,
	}

	rev := review.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Questions: review.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "",
				Type:        review.SINGLE,
			},
		},
		Schedule: review.MONTHLY,
	}

	got := review.Create(rev)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestAllQuestionsCompleteValid(t *testing.T) {
	t.Parallel()

	questions := review.Questions{
		{
			Title:       "What were my biggest wins?",
			Description: "Some description goes here.",
			Answer:      "Reading 10 pages each day of Atomic Habits",
			Type:        review.SINGLE,
		},
		{
			Title:       "What had the most impact on my month?",
			Description: "Some description goes here.",
			Answer:      "Reading 10 everday",
			Type:        review.SINGLE,
		},
	}
	want := true

	got, err := questions.AllComplete()
	if err != nil {
		t.Fatal(err)
	}

	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}

	questions = review.Questions{
		{
			Title:       "What were my biggest wins?",
			Description: "Some description goes here.",
			Answer:      "",
			Type:        review.SINGLE,
		},
		{
			Title:       "What had the most impact on my month?",
			Description: "Some description goes here.",
			Answer:      "Reading 10 everday",
			Type:        review.SINGLE,
		},
	}
	want = false
	got, err = questions.AllComplete()
	if err != nil {
		t.Fatal(err)
	}

	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

func TestAllQuestionsCompleteInvalid(t *testing.T) {
	t.Parallel()
	questions := review.Questions{}

	_, err := questions.AllComplete()
	if err == nil {
		t.Fatalf("want error checking all questions complete but got nil")
	}
}

func TestFinishReview(t *testing.T) {
	t.Parallel()

	want := review.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  true,
		Questions: review.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "Reading 10 pages each day of Atomic Habits",
				Type:        review.SINGLE,
			},
		},
		Schedule: review.MONTHLY,
	}

	myReview := review.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  false,
		Questions: review.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "Reading 10 pages each day of Atomic Habits",
				Type:        review.SINGLE,
			},
		},
		Schedule: review.MONTHLY,
	}

	err := myReview.Finish()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, myReview) {
		t.Error(cmp.Diff(want, myReview))
	}
}

func TestFinishReviewInvalid(t *testing.T) {
	t.Parallel()

	myReview := review.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  false,
		Questions: review.Questions{},
		Schedule:  review.MONTHLY,
	}

	err := myReview.Finish()
	if err == nil {
		t.Fatalf("want error when finishing a review with no questions but got instead got nil")
	}
}

func TestAskPrintsGivenQuestionAndReturnsAnswer(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)

	question := "What were your biggest wins?"
	userInputStr := "Spent 10 minutes programming in Go"
	userInput := strings.NewReader(userInputStr)

	got := review.AskTo(buf, userInput, question)

	bufString := buf.String()
	if bufString != question {
		t.Errorf("Wanted %s but got %s", question, bufString)
	}

	if userInputStr != got {
		t.Errorf("Wanted %s but got %v", userInputStr, got)
	}
}

func TestDueChecksDueDateAndReturnsIfDueOrNot(t *testing.T) {
	t.Parallel()

	// date := time.Date(2023, time.April, 20)
}
