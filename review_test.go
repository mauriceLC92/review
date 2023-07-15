package review_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mauriceLC92/review"
)

func resetFile(filePath string) {
	err := os.WriteFile(filePath, []byte("[]"), 0644)
	if err != nil {
		fmt.Printf("error resetting file: %v\n", err)
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
	if bufString != fmt.Sprintln(question) {
		t.Errorf("Wanted %q but got %q", question, bufString)
	}

	if userInputStr != got {
		t.Errorf("Wanted %q but got %q", userInputStr, got)
	}
}

func TestMyReview(t *testing.T) {
	t.Parallel()
	currentDate := time.Now()
	_ = review.MyReview{
		CreatedAt: currentDate,
	}
}

func TestDueChecksDueDateAndReturnsTrueIfDue(t *testing.T) {
	t.Parallel()

	testDate := "20-05-2023"
	myTime, _ := time.Parse(review.DAY_MONTH_YEAR_FORMAT, testDate)
	myReview := review.MyReview{
		CreatedAt: myTime,
	}

	wantDue := true
	gotDue := myReview.Due()

	if wantDue != gotDue {
		t.Errorf("wanted %v but got %v", wantDue, gotDue)
	}
}

func TestDueChecksDueDateAndReturnsFalseIfNotDue(t *testing.T) {
	t.Parallel()

	testDate := "20-06-2025"
	myTime, _ := time.Parse(review.DAY_MONTH_YEAR_FORMAT, testDate)
	myReview := review.MyReview{
		CreatedAt: myTime,
	}

	wantDue := false
	gotDue := myReview.Due()

	if wantDue != gotDue {
		t.Errorf("wanted %v but got %v", wantDue, gotDue)
	}
}

func TestCreatedTodayChecksIfAReviewWasCreatedTodayAndReturnsTrue(t *testing.T) {
	t.Parallel()

	r := review.MyReview{
		CreatedAt: time.Now(),
	}

	want := true
	got := r.CreatedToday()

	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

func TestCreatedTodayChecksIfAReviewWasCreatedTodayAndReturnsFalse(t *testing.T) {
	t.Parallel()

	r := review.MyReview{
		CreatedAt: time.Now().AddDate(5, 0, 1),
	}

	want := false
	got := r.CreatedToday()

	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

func TestNextDueDateReturnsTheNextDueDate(t *testing.T) {
	t.Parallel()

	testDate := "20-05-2023"
	myTime, _ := time.Parse(review.DAY_MONTH_YEAR_FORMAT, testDate)
	myReview := review.MyReview{
		CreatedAt: myTime,
	}

	wantDate := myTime.AddDate(0, 1, 0)
	gotDate := myReview.NextDueDate()

	if wantDate != gotDate {
		t.Errorf("wanted %v but got %v", wantDate, gotDate)
	}
}

func TestCheckDeterminesIfNoReviewsHaveTakenPlace(t *testing.T) {
	t.Parallel()
	testTime := time.Now()

	review.Now = func() time.Time {
		return testTime
	}
	reviews := []review.MyReview{}

	want := review.MyReview{
		CreatedAt: testTime,
		Questions: review.DefaultQuestions,
	}
	got, ok := review.Check(reviews)

	if ok != false {
		t.Error("expected false when no review has been done but got 'true'")
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCheckDeterminesIfAReviewHasTakenPlaceAndReturnsLatestReview(t *testing.T) {
	t.Parallel()

	reviews := []review.MyReview{
		{CreatedAt: time.Date(2025, time.June, 20, 0, 0, 0, 0, time.UTC)},
		{CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC)},
		{CreatedAt: time.Date(2025, time.July, 20, 0, 0, 0, 0, time.UTC)},
	}
	want := review.MyReview{
		CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC),
	}
	got, ok := review.Check(reviews)

	if ok != true {
		t.Error("expected 'true' when a review has taken place but got 'false'")
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseReadsJSONFileOfReviewsAndReturnsASliceOfReviews(t *testing.T) {
	t.Parallel()

	wantedReviews := []review.MyReview{
		{CreatedAt: time.Date(2025, time.June, 20, 0, 0, 0, 0, time.UTC)},
		{CreatedAt: time.Date(2025, time.July, 20, 0, 0, 0, 0, time.UTC)},
		{CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC)},
	}

	got, err := review.Parse("testdata/reviews-non-empty.json")
	if err != nil {
		t.Fatal("error parsing file", err)
	}

	if !cmp.Equal(got, wantedReviews) {
		t.Error(cmp.Diff(wantedReviews, got))
	}
}

func TestParseReadsJSONFileOfReviewsAndErrorsWhenFileNotFound(t *testing.T) {
	t.Parallel()

	_, err := review.Parse("nowheretobefound/reviews-non-empty.json")
	if err == nil {
		t.Fatal("expected error reading a file that does not exist but got nil instead")
	}
}

func TestParseReadsJSONFileOfReviewsAndErrorsParsingInvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := review.Parse("testdata/reviews-invalid.json")
	if err == nil {
		t.Fatal("expected error reading a file with invalid JSON but got nil instead", err)
	}
}

func TestReviewAsksQuestionsAndGetsAnswers(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	r := &review.MyReview{
		CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC),
		Questions: []review.MyQuestion{
			{Title: "What could you have done better last month?", Answer: ""},
		},
	}

	userAnswer := "Read 100 more pages of AWS Fundamentals"
	wantedReview := review.MyReview{
		CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC),
		Questions: []review.MyQuestion{
			{Title: "What could you have done better last month?", Answer: userAnswer},
		},
	}

	r.Review(buf, strings.NewReader(userAnswer))
	rev := *r
	if !cmp.Equal(rev, wantedReview) {
		t.Error(cmp.Diff(wantedReview, r))
	}
}

func TestSaveWillSaveAReview(t *testing.T) {
	t.Parallel()

	r := review.MyReview{
		CreatedAt: time.Date(2025, time.July, 9, 0, 0, 0, 0, time.UTC),
		Questions: []review.MyQuestion{
			{Title: "How are you today?", Answer: "Fantastic"},
			{Title: "What was your biggest win this month?", Answer: "Writing this test"},
		},
	}

	err := review.SaveTo(r, "testdata/reviews-save.json")
	if err != nil {
		t.Fatal(err)
	}

	reviews, err := review.Parse("testdata/reviews-save.json")
	if err != nil {
		t.Fatal(err)
	}

	got := reviews[0]
	if !cmp.Equal(got, r) {
		t.Error(cmp.Diff(got, r))
	}

	resetFile("testdata/reviews-save.json")
}

func TestAnsweredChecksIfAnyQuestionWasAnswered(t *testing.T) {
	t.Parallel()

	want := true
	r := review.MyReview{
		CreatedAt: time.Date(2025, time.July, 9, 0, 0, 0, 0, time.UTC),
		Questions: []review.MyQuestion{
			{Title: "How are you today?", Answer: "Fantastic"},
			{Title: "What was your biggest win this month?", Answer: "Writing this test"},
		},
	}

	got := r.Answered()
	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}

	want = false
	r = review.MyReview{
		CreatedAt: time.Date(2025, time.July, 9, 0, 0, 0, 0, time.UTC),
		Questions: []review.MyQuestion{
			{Title: "How are you today?", Answer: ""},
			{Title: "What was your biggest win this month?", Answer: ""},
		},
	}
	got = r.Answered()
	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

func TestOpenJSONStoreOpensFileAndReturnsJSONStore(t *testing.T) {
	t.Parallel()

	want := []review.MyReview{
		{
			CreatedAt: time.Date(2025, time.July, 9, 0, 0, 0, 0, time.UTC),
		},
	}

	store, err := review.OpenJSONStore("testdata/reviews-json-store.json")
	if err != nil {
		t.Fatal(err)
	}

	reviews := store.GetAll()
	if !cmp.Equal(want, reviews) {
		t.Error(cmp.Diff(want, reviews))
	}
}

func TestGetLatestReviewFetchesTheLatestReviewFromJSONStore(t *testing.T) {
	t.Parallel()

	want := review.MyReview{
		CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC),
	}

	store, err := review.OpenJSONStore("testdata/reviews-json-store-latest-review.json")
	if err != nil {
		t.Fatal(err)
	}

	got, ok := store.GetLatestReview()
	if ok != true {
		t.Error("expected 'true' when a review has taken place but got 'false'")
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestSaveWillSaveAReviewToTheJSONStore(t *testing.T) {
	t.Parallel()

	store, err := review.OpenJSONStore("testdata/reviews-json-store-save.json")
	if err != nil {
		t.Fatal(err)
	}
	want := review.MyReview{
		CreatedAt: time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC),
	}

	err = store.Save(want)
	if err != nil {
		t.Fatal(err)
	}

	store, err = review.OpenJSONStore("testdata/reviews-json-store-save.json")
	if err != nil {
		t.Fatal(err)
	}

	reviews := store.GetAll()
	got := reviews[0]
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
