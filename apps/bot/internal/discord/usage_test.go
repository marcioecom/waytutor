package discord

import "testing"

func TestSparkline(t *testing.T) {
	if got := sparkline(nil); got != "—" {
		t.Fatalf("unexpected empty sparkline: %q", got)
	}
	if got := sparkline([]int64{0, 1, 2}); got != "▁▅█" {
		t.Fatalf("unexpected sparkline: %q", got)
	}
	if got := sparkline([]int64{0, 0, 0}); got != "▁▁▁" {
		t.Fatalf("unexpected flat sparkline: %q", got)
	}
}

func TestAllZero(t *testing.T) {
	if !allZero(nil) || !allZero([]int64{0, 0}) {
		t.Fatal("expected an empty series to be all zero")
	}
	if allZero([]int64{0, 1}) {
		t.Fatal("expected a series with usage not to be all zero")
	}
}

func TestNumber(t *testing.T) {
	cases := map[int64]string{
		0:       "0",
		999:     "999",
		1000:    "1.000",
		1234567: "1.234.567",
		-1200:   "-1.200",
	}
	for value, want := range cases {
		if got := number(value); got != want {
			t.Fatalf("number(%d) = %q, want %q", value, got, want)
		}
	}
}

func TestPercent(t *testing.T) {
	if got := percent(50, 200); got != "25%" {
		t.Fatalf("unexpected percent: %q", got)
	}
	if got := percent(1, 0); got != "—" {
		t.Fatalf("unexpected empty percent: %q", got)
	}
}

func TestFormatDay(t *testing.T) {
	if got := formatDay("2026-11-07"); got != "07/11" {
		t.Fatalf("unexpected day: %q", got)
	}
	if got := formatDay("not a day"); got != "not a day" {
		t.Fatalf("unexpected fallback: %q", got)
	}
}
