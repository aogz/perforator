package commands

import (
	"testing"
	"time"

	"github.com/google/go-github/v40/github"
)

func TestCalculateReviewTimeByAuthor(t *testing.T) {
	now := time.Now()
	in6Hours := now.Add(time.Hour * 6)
	in12Hours := now.Add(time.Hour * 12)
	tomorrow := now.Add(time.Hour * 24)
	username1 := "Foo"
	username2 := "Bar"
	prId := 1
	prs := []*github.PullRequest{
		{
			User: &github.User{
				Login: &username1,
			},
			Number:    &prId,
			CreatedAt: &now,
			MergedAt:  &in6Hours,
		},
		{
			User: &github.User{
				Login: &username1,
			},
			Number:    &prId,
			CreatedAt: &now,
			MergedAt:  &tomorrow,
		},
		{
			User: &github.User{
				Login: &username2,
			},
			Number:    &prId,
			CreatedAt: &now,
			MergedAt:  &in12Hours,
		},
		{
			User: &github.User{
				Login: &username2,
			},
			Number:    &prId,
			CreatedAt: &now,
			MergedAt:  nil,
		},
	}

	initial := map[string][]time.Duration{}
	result := calculateReviewTimeByAuthor(initial, prs, 100)
	if len(result) != 2 {
		t.Fatalf("Unexpected result length, expected 2 got %d", len(result))
	}
	if len(result[username1]) != 2 {
		t.Fatalf("Unexpected len of results for user1, expected 2 got %d", len(result[username1]))
	}
	if result[username1][0] != time.Hour*6 && result[username1][1] != time.Hour*24 {
		t.Fatalf("Unexpected results for user1, expected 6 and 24 got %d %d", result[username1][0], result[username1][1])
	}
	if len(result[username2]) != 1 {
		t.Fatalf("Unexpected len of results for user1, expected 2 got %d", len(result[username1]))
	}
	if result[username2][0] != time.Hour*12 {
		t.Fatalf("Unexpected results for user1, expected 12 got %d", result[username2][0])
	}
}

func TestGetResults(t *testing.T) {
	stats := map[string][]time.Duration{
		"foo": {time.Minute * 1, time.Minute * 2, time.Minute * 3},
		"bar": {time.Minute * 4, time.Minute * 5, time.Minute * 6},
		"baz": {time.Minute * 7, time.Minute * 8, time.Minute * 9},
	}

	results := calculateAggregatedResultsPerUser(stats)
	expected := 5.0 / 60.0
	if results != float64(expected) {
		t.Fatalf("Invalid aggregated results, expected %f got %f", expected, results)
	}
}

func TestCalculateUserReviewTime(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(time.Hour * 24)

	pr := github.PullRequest{
		CreatedAt: &now,
	}

	review := github.PullRequestReview{
		SubmittedAt: &tomorrow,
	}

	reviewTime := calculateUserReviewTime(&pr, &review)
	if reviewTime != time.Hour*24 {
		t.Fatalf("Invalid time in review, expected 24h got %s", reviewTime)
	}
}

func TestCalculatePRInReviewTime(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(time.Hour * 24)

	pr := github.PullRequest{
		CreatedAt: &now,
		MergedAt:  &tomorrow,
	}

	timeInReview := calculatePRInReviewTime(&pr)
	if timeInReview != time.Hour*24 {
		t.Fatalf("Invalid time in review, expected 24h got %s", timeInReview)
	}
}

func TestCalculatePRInReviewTimeNotMerged(t *testing.T) {
	now := time.Now()

	pr := github.PullRequest{
		CreatedAt: &now,
		MergedAt:  nil,
	}

	timeInReview := calculatePRInReviewTime(&pr)
	if timeInReview != 0 {
		t.Fatalf("Invalid time in review, expected 0 got %s", timeInReview)
	}
}

func TestCalculateAverageReviewTimePerUser(t *testing.T) {
	durations := []time.Duration{
		time.Minute * 1,
		time.Minute * 2,
		time.Minute * 3,
	}
	averageReviewTime := calculateAverageReviewTimePerUser(durations)
	expected := 120.0 / 60 / 60
	if averageReviewTime != expected {
		t.Fatalf("Invalid averageReviewTime, expected %f got %f", expected, averageReviewTime)
	}

}

func TestCalculateAggregatedAverageReviewTime(t *testing.T) {
	if calculateAggregatedAverageReviewTime(100.0, 10) != 10.0 {
		t.Fatal("Invalid result")
	}

	if calculateAggregatedAverageReviewTime(24.0, 10) != 2.4 {
		t.Fatal("Invalid result")
	}

	if calculateAggregatedAverageReviewTime(15.0, 10) != 1.5 {
		t.Fatal("Invalid result")
	}
}
