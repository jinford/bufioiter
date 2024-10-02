package bufioiter

import (
	"bufio"
	"io"
	"iter"
)

type scanner struct {
	*bufio.Scanner
}

func NewScanner(r io.Reader, opts ...ScannerOption) iter.Seq2[string, error] {
	s := &scanner{
		Scanner: bufio.NewScanner(r),
	}
	for _, opt := range opts {
		opt.apply(s)
	}

	return func(yield func(string, error) bool) {
		for s.Scan() {
			if !yield(s.Text(), nil) {
				break
			}
		}
		if err := s.Err(); err != nil {
			yield("", err)
		}
	}
}

type ScannerOption interface {
	apply(*scanner)
}

type bufioScannerOptionFunc func(*scanner)

func (f bufioScannerOptionFunc) apply(s *scanner) {
	f(s)
}

func Bufferr(buf []byte, max int) ScannerOption {
	return bufioScannerOptionFunc(func(s *scanner) {
		s.Buffer(buf, max)
	})
}

func Split(split bufio.SplitFunc) ScannerOption {
	return bufioScannerOptionFunc(func(s *scanner) {
		s.Split(split)
	})
}
