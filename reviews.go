package reviews

/**
	1. As a user I would like to create a review to review on a pre-determined schedule. Monthly, weekly is fine for now
	2. As a user I should be able to mark a review as complete or not
		- Fill in the questions
	3. As a user I should be able to view my past reviews
**/

const (
	MONTHLY = "monthly"
	WEEKLY  = "weekly"
)

type Question struct {
	Title       string
	Description string
	Answer      string
}

type Review struct {
	Id            string
	CreatedAt     int64  // unix timestamp
	UpdatedAt     int64  // unix timestamp
	DateFormatted string // formatted into a utc date
	Complete      bool
	Questions     []Question
	Schedule      string
}
