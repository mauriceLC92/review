package reviews

/**
	1. As a user I would like to create a monthly review
	2. As a user I should be able to mark a review as complete or not
		- Fill in the questions
	3. As a user I should be able to view my past reviews
**/

type Question struct {
	Title       string
	Description string
}

type Schedule = string

type Review struct {
	Id            string
	Date          int64
	DateFormatted string
	Complete      bool
	Questions     []Question
	Schedule
}

func Create(schedule string, questions []Question) Review {
	return Review{}
}
