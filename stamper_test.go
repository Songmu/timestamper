package timestamper_test

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
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
		name:    "timestamp at buffer boundary",
		input:   strings.Repeat("x", 4085) + "\nY\n",
		options: []timestamper.Option{timestamper.Layout("[ts] ")},
		reg:     regexp.MustCompile(`\A\[ts\] x{1000}x{1000}x{1000}x{1000}x{85}\n\[ts\] Y\n\z`),
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
				b, err := io.ReadAll(r)
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

func TestStamperShortDst(t *testing.T) {
	for _, tc := range []struct {
		name    string
		size    int
		want    string
		wantSrc int
		wantErr error
	}{
		{name: "empty", size: 0, wantErr: transform.ErrShortDst},
		{name: "timestamp only", size: 5, wantErr: transform.ErrShortDst},
		{name: "partial line at EOF", size: 7, want: "[ts] AB", wantSrc: 2, wantErr: transform.ErrShortDst},
		{name: "exact fit", size: 8, want: "[ts] ABC", wantSrc: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := timestamper.New(timestamper.Layout("[ts] "))
			src := []byte("ABC")
			dst := make([]byte, tc.size)
			nDst, nSrc, err := s.Transform(dst, src, true)
			if string(dst[:nDst]) != tc.want || nSrc != tc.wantSrc || err != tc.wantErr {
				t.Fatalf("Transform = (%q, %d, %v), want (%q, %d, %v)",
					dst[:nDst], nSrc, err, tc.want, tc.wantSrc, tc.wantErr)
			}
			got := string(dst[:nDst])
			dst = make([]byte, 64)
			nDst, _, err = s.Transform(dst, src[nSrc:], true)
			if err != nil {
				t.Fatal(err)
			}
			if got += string(dst[:nDst]); got != "[ts] ABC" {
				t.Fatalf("output = %q, want %q", got, "[ts] ABC")
			}
		})
	}
}

type eofReader struct {
	*strings.Reader
}

func (r eofReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if r.Len() == 0 {
		err = io.EOF
	}
	return n, err
}

func TestStamperReaderEOF(t *testing.T) {
	input := strings.Repeat("x", 4095)
	r := transform.NewReader(eofReader{strings.NewReader(input)},
		timestamper.New(timestamper.Layout("[ts] ")))
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if want := "[ts] " + input; string(got) != want {
		t.Fatalf("output length = %d, want %d; timestamps = %d, want 1",
			len(got), len(want), bytes.Count(got, []byte("[ts] ")))
	}
}

func TestWriterConcurrent(t *testing.T) {
	var out bytes.Buffer
	w := timestamper.NewWriter(&out, timestamper.Layout("[ts] "))
	const workers = 16
	want := make(map[string]bool, workers)
	errs := make(chan error, workers+1)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		line := fmt.Sprintf("worker-%02d:", i) + strings.Repeat("x", 5000)
		want["[ts] "+line] = true
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := io.WriteString(w, line+"\n")
			if err != nil {
				errs <- err
			} else if n != len(line)+1 {
				errs <- fmt.Errorf("wrote %d bytes, want %d", n, len(line)+1)
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := w.Close(); err != nil {
			errs <- err
		}
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n") {
		if !want[line] {
			t.Fatalf("unexpected or duplicate line: %.80q", line)
		}
		delete(want, line)
	}
	if len(want) != 0 {
		t.Fatalf("missing %d lines", len(want))
	}
}
