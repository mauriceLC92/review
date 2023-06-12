package reviews

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

type Schedule string

type Review struct {
	ID        string
	CreatedAt int64 // unix timestamp
	UpdatedAt int64 // unix timestamp
	Complete  bool
	Questions []Question
	Schedule
}

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
	// we will a way to check if all questions have been answered
	return nil
}
