package reviews_test

import (
	"testing"
	"time"

	"github.com/mauriceLC92/reviews"

	"github.com/google/go-cmp/cmp"
)

func TestReview(t *testing.T) {
	t.Parallel()
	currentDate := time.Now()
	_ = reviews.Review{
		Id:            "123",
		Date:          currentDate.UnixMilli(),
		DateFormatted: currentDate.Format(time.UnixDate),
		Questions: []reviews.Question{
			{
				Title:       "What is hello world?",
				Description: "A questio to get you thinking",
			},
		},
		Schedule: "Monthly",
	}
}

func TestCreateReview(t *testing.T) {
	t.Parallel()
	currentDate := time.Now()
	want := reviews.Review{
		Id:            "123",
		Date:          currentDate.UnixMilli(),
		DateFormatted: currentDate.Format(time.UnixDate),
		Questions: []reviews.Question{
			{
				Title:       "What is hello world?",
				Description: "A questio to get you thinking",
			},
		},
		Schedule: "Monthly",
	}

	questions := []reviews.Question{
		{
			Title:       "What is hello world?",
			Description: "A questio to get you thinking",
		},
	}

	got := reviews.Create("monthly", questions)

	if !cmp.Equal(got, want) {
		t.Error("created credentials do not match the expected", cmp.Diff(want, got))
	}
}
