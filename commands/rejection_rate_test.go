package commands

import "testing"

func TestCalculateRejectionRatePerUser20(t *testing.T) {
	stats := UserRejectionRateStatistic{
		Rejected: 1,
		Total:    5,
	}

	result := calculateRejectionRatePerUser(stats)
	expected := 20.0
	if result != float64(expected) {
		t.Fatalf("Invalid result, expected %.2f got %.2f\n", expected, result)
	}
}

func TestCalculateRejectionRatePerUser0(t *testing.T) {
	stats := UserRejectionRateStatistic{
		Rejected: 0,
		Total:    5,
	}

	result := calculateRejectionRatePerUser(stats)
	expected := 0.0
	if result != float64(expected) {
		t.Fatalf("Invalid result, expected %.2f got %.2f\n", expected, result)
	}
}

func TestCalculateAggregatedRate10(t *testing.T) {
	stats := map[string]UserRejectionRateStatistic{
		"foo": {
			Rejected: 1,
			Total:    5,
		},
		"bar": {
			Rejected: 0,
			Total:    5,
		},
	}
	result := calculateAggregatedRate(stats)
	expected := 10.0
	if result != float64(expected) {
		t.Fatalf("Invalid result, expected %.2f got %.2f\n", expected, result)
	}
}

func TestTryCreateUserStats(t *testing.T) {
	stats := map[string]UserRejectionRateStatistic{}
	tryCreateUserStats(stats, "foo")
	tryCreateUserStats(stats, "bar")
	tryCreateUserStats(stats, "baz")
	tryCreateUserStats(stats, "baz")

	expected := 3
	if len(stats) != expected {
		t.Fatalf("Invalid stats len, expected %d got %d\n", expected, len(stats))
	}

	if userStats, ok := stats["foo"]; !ok {
		t.Fatalf("Expected foo, but not found")

		if userStats.Total != 0 {
			t.Fatalf("Unexpected userStats.Total, expected 0 got %d", userStats.Total)
		}

		if userStats.Rejected != 0 {
			t.Fatalf("Unexpected userStats.Rejected, expected 0 got %d", userStats.Rejected)
		}
	}
}
