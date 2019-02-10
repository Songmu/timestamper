package timestamper

import (
	"bytes"
	"time"

	"golang.org/x/text/transform"
)

// Option type of timestamper
type Option func(*stamper)

// UTC for using utc in timestamp
func UTC() Option {
	return func(s *stamper) {
		s.utc = true
	}
}

// Layout for specifying custom timestamp layout
func Layout(layout string) Option {
	return func(s *stamper) {
		s.layout = layout
	}
}

// New returns new timestamper
func New(opts ...Option) transform.Transformer {
	s := &stamper{layout: defaultLayout}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type stamper struct {
	layout    string
	midOfLine bool
	utc       bool
}

func (s *stamper) stampLen() int {
	return len(s.layout)
}

// Reset implements transform.Transformer.Reset.
func (s *stamper) Reset() {
	s.midOfLine = false
}

const defaultLayout = "2006-01-02T15:04:05.000000Z07:00 " // RFC3339Micro

// Transform implements transform.Transformer.Transform.
func (s *stamper) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	var buf bytes.Buffer
	var nDstTemp int
	for _, chr := range src {
		if !s.midOfLine {
			buf.Write(s.timestampBytes())
			nDstTemp += s.stampLen()
			s.midOfLine = true
		}
		if chr == '\n' {
			s.midOfLine = false
		}
		buf.WriteByte(chr)
		nDstTemp++
		nSrc++
	}
	nDst = copy(dst, buf.Bytes())
	if nDst < nDstTemp {
		err = transform.ErrShortDst
	}
	return
}

func (s *stamper) timestampBytes() []byte {
	t := time.Now()
	if s.utc {
		t = t.UTC()
	}
	return s.formatTimestamp(t)
}

func (s *stamper) formatTimestamp(t time.Time) []byte {
	b := make([]byte, 0, s.stampLen())
	return t.AppendFormat(b, s.layout)
}
