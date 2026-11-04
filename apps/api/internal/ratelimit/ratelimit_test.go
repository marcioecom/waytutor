package ratelimit

import "testing"

func TestAllowSpendsTheBurstThenDenies(t *testing.T) {
	limiter := New(3)

	for i := range 3 {
		if !limiter.Allow("account") {
			t.Fatalf("Allow() = false on request %d, want true", i+1)
		}
	}
	if limiter.Allow("account") {
		t.Fatal("Allow() = true after the burst, want false")
	}
}

func TestAllowIsPerKey(t *testing.T) {
	limiter := New(1)

	if !limiter.Allow("first") {
		t.Fatal("Allow(first) = false, want true")
	}
	if !limiter.Allow("second") {
		t.Fatal("Allow(second) = false, want true")
	}
	if limiter.Allow("first") {
		t.Fatal("Allow(first) = true after its burst, want false")
	}
}
