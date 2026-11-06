package timestamper

import (
	"bytes"
	"io"
	"sync"
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

// NewWriter returns a timestamping writer safe for concurrent Write and Close calls.
// Close flushes the transformer without closing w.
func NewWriter(w io.Writer, opts ...Option) io.WriteCloser {
	return &writer{w: transform.NewWriter(w, New(opts...))}
}

type writer struct {
	mu sync.Mutex
	w  *transform.Writer
}

func (w *writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Write(p)
}

func (w *writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Close()
}

type stamper struct {
	layout    string
	midOfLine bool
	utc       bool

	mu sync.Mutex
}

// Reset implements transform.Transformer.Reset.
func (s *stamper) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reset()
}

func (s *stamper) reset() {
	s.midOfLine = false
}

const defaultLayout = "2006-01-02T15:04:05.000000Z07:00 " // RFC3339Micro

// Transform implements transform.Transformer.Transform.
func (s *stamper) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var buf bytes.Buffer
	for _, chr := range src {
		if buf.Len() == len(dst) {
			err = transform.ErrShortDst
			break
		}
		if !s.midOfLine {
			ts := s.timestampBytes()
			if len(dst)-buf.Len() < len(ts)+1 {
				err = transform.ErrShortDst
				break
			}
			buf.Write(ts)
		}
		buf.WriteByte(chr)
		nSrc++
		s.midOfLine = chr != '\n'
	}
	nDst = copy(dst, buf.Bytes())
	if atEOF && err == nil {
		s.reset()
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
	const defaultMax = 64
	max := len(s.layout) + 10
	if max < defaultMax {
		max = defaultMax
	}
	b := make([]byte, 0, max)
	return t.AppendFormat(b, s.layout)
}
