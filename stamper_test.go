package timestamper_test

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/Songmu/timestamper"
	"golang.org/x/text/transform"
)

func TestStamper(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		options []timestamper.Option
		reg     *regexp.Regexp
	}{{
		name: "basic",
		input: `Hello
World!`,
		reg: func() *regexp.Regexp {
			timeRegStr := `[0-9]{4}-[0-9]{2}-[0-9]{2}T` +
				`[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{6}` +
				`(?:Z|[-+][0-9]{2}:[0-9]{2}) `
			return regexp.MustCompile(`\A` +
				timeRegStr + "Hello\n" +
				timeRegStr + `World!\z`)
		}(),
	}, {
		name: "utc",
		input: `Hello
World!`,
		options: []timestamper.Option{timestamper.UTC()},
		reg: func() *regexp.Regexp {
			timeRegStr := `[0-9]{4}-[0-9]{2}-[0-9]{2}T` +
				`[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{6}Z `
			return regexp.MustCompile(`\A` +
				timeRegStr + "Hello\n" +
				timeRegStr + `World!\z`)
		}(),
	}, {
		name: "layout",
		input: `Hello
World!`,
		reg: func() *regexp.Regexp {
			timeRegStr := `[0-9]{2}-[0-9]{2}-[0-9]{2} ` +
				`[0-9]{2}:[0-9]{2}:[0-9]{2} `
			return regexp.MustCompile(`\A` +
				timeRegStr + "Hello\n" +
				timeRegStr + `World!\z`)
		}(),
		options: []timestamper.Option{timestamper.Layout("06-01-02 15:04:05 ")},
	}, {
		name:  "long",
		input: strings.Repeat("s", 4100) + "\n" + strings.Repeat("S", 4100) + "\nHello!",
		reg: func() *regexp.Regexp {
			timeRegStr := `[0-9]{4}-[0-9]{2}-[0-9]{2}T` +
				`[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{6}` +
				`(?:Z|[-+][0-9]{2}:[0-9]{2}) `
			return regexp.MustCompile(`\A` +
				timeRegStr + "s{1000}s{1000}s{1000}s{1000}s{100}\n" +
				timeRegStr + "S{1000}S{1000}S{1000}S{1000}S{100}\n" +
				timeRegStr + `Hello!\z`)
		}(),
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.input
			reg := tc.reg
			s := timestamper.New(tc.options...)

			t.Run("Reader", func(t *testing.T) {
				r := transform.NewReader(strings.NewReader(input), s)
				b, err := ioutil.ReadAll(r)
				if err != nil {
					t.Errorf("something went wrong: %s", err)
				}
				out := string(b)
				if !reg.MatchString(out) {
					t.Errorf("something went wrong. output: %s", out)
				}
				t.Logf("Output:\n%s", out)
			})

			s.Reset()

			t.Run("Writer", func(t *testing.T) {
				buf := &bytes.Buffer{}
				w := transform.NewWriter(buf, s)
				w.Write([]byte(input))
				w.Close()
				out := buf.String()
				if !reg.MatchString(out) {
					t.Errorf("something went wrong. output: %s", out)
				}
				t.Logf("Output:\n%s", out)
			})
		})
	}
}
