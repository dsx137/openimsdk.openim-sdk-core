package checker

import (
	"testing"
)

func TestClassifyUnreadMismatch_distinguishes_failure_stage(t *testing.T) {
	tests := []struct {
		name    string
		local   int64
		server  int64
		correct int64
		want    string
	}{
		{name: "local unread state", local: 9, server: 10, correct: 10, want: "local_unread_state"},
		{name: "server state or expectation", local: 9, server: 9, correct: 10, want: "server_state_or_expectation"},
		{name: "mixed state", local: 8, server: 9, correct: 10, want: "mixed_local_and_server_state"},
		{name: "no mismatch", local: 10, server: 10, correct: 10, want: "no_mismatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyUnreadMismatch(tt.local, tt.server, tt.correct)
			if got != tt.want {
				t.Fatalf("classifyUnreadMismatch() = %q, want %q", got, tt.want)
			}
		})
	}
}
