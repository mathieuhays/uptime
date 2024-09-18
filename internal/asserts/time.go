package asserts

import (
	"testing"
	"time"
)

func SimilarTime(t testing.TB, got, want time.Time) {
	t.Helper()
	diff := got.Sub(want).Abs()
	if diff > time.Minute {
		t.Errorf("time doesn't match. expected: %s. got: %s", want, got)
	}
}
