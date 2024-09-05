package tsplot

import (
	"compress/gzip"
	"regexp"
	"io"
	"bufio"
	"os"
	"github.com/jgbaldwinbrown/lscan/lscan"
	"fmt"
)

type SyncIter struct {
	io.Closer
	Scanner *bufio.Scanner
	Line []string
	Scanf lscan.Splitter
}

// Iterate over entries of a sync file
func NewSyncIter(r io.ReadCloser) *SyncIter {
	s := new(SyncIter)
	s.Closer = r
	s.Scanner = bufio.NewScanner(r)
	s.Scanner.Buffer([]byte{}, 1e12)
	s.Scanf = lscan.ByByte('\t')
	return s
}

type GzReader struct {
	FileCloser io.Closer
	*gzip.Reader
}

func (g GzReader) Close() error {
	g.Reader.Close()
	return g.FileCloser.Close()
}

// Open sync file, which may be gzipped, as a *SyncIter
func OpenSyncIter(path string) (*SyncIter, error) {

	conn, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("OpenSyncIter: %w", err)
	}

	var g *gzip.Reader
	var r io.ReadCloser

	re := regexp.MustCompile(`\.gz$`)
	if re.MatchString(path) {
		g, err = gzip.NewReader(conn)
		if err != nil {
			return nil, fmt.Errorf("OpenSyncIter: %w", err)
		}
		r = GzReader{conn, g}
	} else {
		r = conn
	}

	return NewSyncIter(r), nil
}

func (s *SyncIter) Next() (SyncE, bool) {
	ok := s.Scanner.Scan()
	if !ok {
		return SyncE{}, false
	}

	s.Line = lscan.SplitByFunc(s.Line, s.Scanner.Text(), s.Scanf)
	sy, err := ParseSyncE(s.Line)
	if err != nil {
		return SyncE{}, false
	}
	return sy, true
}
