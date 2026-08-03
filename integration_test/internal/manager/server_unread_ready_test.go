package manager

import "testing"

func TestServerUnreadCountsReady_whenAllCountsMatch(t *testing.T) {
	// Given: every offline user has reached the expected server unread count.
	actual := map[string]int{"u16": 416, "u17": 394}
	expected := func(userID string) int {
		return map[string]int{"u16": 416, "u17": 394}[userID]
	}

	// When: readiness is evaluated.
	ready, err := serverUnreadCountsReady(actual, expected)

	// Then: the login barrier is satisfied.
	if err != nil {
		t.Fatalf("serverUnreadCountsReady() error = %v", err)
	}
	if !ready {
		t.Fatal("serverUnreadCountsReady() = false, want true")
	}
}

func TestServerUnreadCountsReady_whenCountIsPending(t *testing.T) {
	// Given: one offline user is still missing a server-visible message.
	actual := map[string]int{"u16": 415, "u17": 394}
	expected := func(userID string) int {
		return map[string]int{"u16": 416, "u17": 394}[userID]
	}

	// When: readiness is evaluated.
	ready, err := serverUnreadCountsReady(actual, expected)

	// Then: the barrier remains pending without failing.
	if err != nil {
		t.Fatalf("serverUnreadCountsReady() error = %v", err)
	}
	if ready {
		t.Fatal("serverUnreadCountsReady() = true, want false")
	}
}

func TestServerUnreadCountsReady_whenCountExceedsExpected(t *testing.T) {
	// Given: a server count exceeds the deterministic test expectation.
	actual := map[string]int{"u16": 417}
	expected := func(string) int { return 416 }

	// When: readiness is evaluated.
	ready, err := serverUnreadCountsReady(actual, expected)

	// Then: the barrier fails instead of polling forever.
	if err == nil {
		t.Fatal("serverUnreadCountsReady() error = nil, want non-nil")
	}
	if ready {
		t.Fatal("serverUnreadCountsReady() = true, want false")
	}
}
