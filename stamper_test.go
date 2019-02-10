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
	input := `Hello
World!`
	timeRegStr := `[0-9]{4}-[0-9]{2}-[0-9]{2}T` +
		`[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{6}` +
		`(?:Z|[-+][0-9]{2}:[0-9]{2}) `
	reg := regexp.MustCompile(`\A` +
		timeRegStr + "Hello\n" +
		timeRegStr + `World!\z`)

	s := timestamper.New()

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
}

func TestStamper_OptionUTC(t *testing.T) {
	input := `Hello
World!`
	timeRegStr := `[0-9]{4}-[0-9]{2}-[0-9]{2}T` +
		`[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{6}Z `
	reg := regexp.MustCompile(`\A` +
		timeRegStr + "Hello\n" +
		timeRegStr + `World!\z`)

	s := timestamper.New(timestamper.UTC())

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
}

func TestStamper_OptionLayout(t *testing.T) {
	input := `Hello
World!`
	timeRegStr := `[0-9]{2}-[0-9]{2}-[0-9]{2} ` +
		`[0-9]{2}:[0-9]{2}:[0-9]{2} `
	reg := regexp.MustCompile(`\A` +
		timeRegStr + "Hello\n" +
		timeRegStr + `World!\z`)

	s := timestamper.New(timestamper.Layout("06-01-02 15:04:05 "))

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
}
