package reviews_test

import (
	"testing"
	"time"

	"github.com/mauriceLC92/reviews"
)

func TestReview(t *testing.T) {
	t.Parallel()
	currentDate := time.Now()
	_ = reviews.Review{
		Id:            "123",
		CreatedAt:     currentDate.UnixMilli(),
		UpdatedAt:     currentDate.UnixMilli(),
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
