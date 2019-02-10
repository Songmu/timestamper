package timestamper

import (
	"testing"
	"time"
)

func TestFormatTimestamp(t *testing.T) {
	s := &stamper{layout: defaultLayout}
	testCases := []struct {
		name   string
		input  time.Time
		expect string
	}{
		{
			name:   "simple",
			input:  time.Date(2019, time.November, 4, 11, 12, 13, 123456000, time.UTC),
			expect: "2019-11-04T11:12:13.123456Z ",
		},
		{
			name:   "pad zero",
			input:  time.Date(2019, time.November, 9, 11, 12, 13, 123400000, time.UTC),
			expect: "2019-11-09T11:12:13.123400Z ",
		},
		{
			name:   "no millisec",
			input:  time.Date(2019, time.November, 3, 11, 12, 13, 0, time.UTC),
			expect: "2019-11-03T11:12:13.000000Z ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := string(s.formatTimestamp(tc.input))
			if got != tc.expect {
				t.Errorf("something went wrong. expect: %s, got: %s", tc.expect, got)
			}
		})
	}
}
