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
		Id:        "123",
		CreatedAt: currentDate.UnixMilli(),
		UpdatedAt: currentDate.UnixMilli(),
		Questions: []reviews.Question{
			{
				Title:       "What were my biggest wins?",
				Description: "A questio to get you thinking",
			},
		},
		Schedule: "Monthly",
	}
}

func TestQuestion(t *testing.T) {
	t.Parallel()
	_ = reviews.Question{
		Title:       "What were my biggest wins?",
		Description: "Some description goes here.",
		Answer:      "",
		Type:        reviews.SINGLE,
	}
}

func TestCreateANewReview(t *testing.T) {
	t.Parallel()

	want := reviews.Review{
		Id:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Complete:  false,
		Questions: []reviews.Question{
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
		Id:        "12",
		CreatedAt: 1685799475,
		UpdatedAt: 1685799487,
		Questions: []reviews.Question{
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
