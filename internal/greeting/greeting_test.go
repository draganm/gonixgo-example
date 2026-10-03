package greeting

import "testing"

func TestFor(t *testing.T) {
	if got, want := For("gonixgo"), "Hello, gonixgo!"; got != want {
		t.Fatalf("For(%q) = %q, want %q", "gonixgo", got, want)
	}
}
