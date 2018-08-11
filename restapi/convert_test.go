package restapi

import (
	"testing"
	"time"
)

func TestTimeToEpochMs(t *testing.T) {
	for n, tc := range []struct {
		in  time.Time
		out int64
	}{
		{time.Time{}, 0},
		{time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2000, 1, 2, 3, 4, 5, 234000000, time.UTC), 946782245234},
	} {
		if got, want := timeToEpochMs(tc.in), tc.out; got != want {
			t.Errorf("[%d] got %d, want %d", n, got, want)
		}
	}
}

func TestEpochMsToTime(t *testing.T) {
	for n, tc := range []struct {
		in  int64
		out time.Time
	}{
		{0, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)},
		{946782245234, time.Date(2000, 1, 2, 3, 4, 5, 234000000, time.UTC)},
	} {
		out := epochMsToTime(tc.in)

		if got, want := out, tc.out; !got.Equal(want) {
			t.Errorf("[%d] got %s, want %s", n, got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
		}

		if got, want := out.Location(), time.UTC; got != want {
			t.Errorf("[%d] location is %q, want %q", n, got, want)
		}
	}
}
