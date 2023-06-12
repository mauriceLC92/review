package reviews_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mauriceLC92/reviews"
)

func TestReview(t *testing.T) {
	t.Parallel()
	currentDate := time.Now()
	_ = reviews.Review{
		ID:        "123",
		CreatedAt: currentDate.UnixMilli(),
		UpdatedAt: currentDate.UnixMilli(),
		Questions: []reviews.Question{
			{
				Title:       "What were my biggest wins?",
				Description: "A questio to get you thinking",
			},
		},
		Schedule: reviews.MONTHLY,
	}
}

func TestQuestion(t *testing.T) {
	t.Parallel()
	_ = reviews.Question{
		ID:          "321",
		Title:       "What were my biggest wins?",
		Description: "Some description goes here.",
		Answer:      "",
		Type:        reviews.SINGLE,
	}
}

func TestCreateANewReview(t *testing.T) {
	t.Parallel()

	want := reviews.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  false,
		Questions: reviews.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "",
				Type:        reviews.SINGLE,
			},
		},
		Schedule: reviews.MONTHLY,
	}

	review := reviews.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Questions: reviews.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "",
				Type:        reviews.SINGLE,
			},
		},
		Schedule: reviews.MONTHLY,
	}

	got := reviews.Create(review)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFinishReview(t *testing.T) {
	t.Parallel()

	want := reviews.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  true,
		Questions: reviews.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "Reading 10 pages each day of Atomic Habits",
				Type:        reviews.SINGLE,
			},
		},
		Schedule: reviews.MONTHLY,
	}

	myReview := reviews.Review{
		ID:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  false,
		Questions: reviews.Questions{
			{
				Title:       "What were my biggest wins?",
				Description: "Some description goes here.",
				Answer:      "Reading 10 pages each day of Atomic Habits",
				Type:        reviews.SINGLE,
			},
		},
		Schedule: reviews.MONTHLY,
	}

	err := myReview.Finish()
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, myReview) {
		t.Error(cmp.Diff(want, myReview))
	}
}

func TestAllQuestionsCompleteValid(t *testing.T) {
	t.Parallel()

	questions := reviews.Questions{
		{
			Title:       "What were my biggest wins?",
			Description: "Some description goes here.",
			Answer:      "Reading 10 pages each day of Atomic Habits",
			Type:        reviews.SINGLE,
		},
		{
			Title:       "What had the most impact on my month?",
			Description: "Some description goes here.",
			Answer:      "Reading 10 everday",
			Type:        reviews.SINGLE,
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

	questions = reviews.Questions{
		{
			Title:       "What were my biggest wins?",
			Description: "Some description goes here.",
			Answer:      "",
			Type:        reviews.SINGLE,
		},
		{
			Title:       "What had the most impact on my month?",
			Description: "Some description goes here.",
			Answer:      "Reading 10 everday",
			Type:        reviews.SINGLE,
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
	questions := reviews.Questions{}

	_, err := questions.AllComplete()
	if err == nil {
		t.Fatalf("want error checking all questions complete but got nil")
	}
}
