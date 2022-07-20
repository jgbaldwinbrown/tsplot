package tsplot

import (
	"regexp"
	"os/exec"
	"fmt"
	"io"
	"strings"
	"strconv"
	"errors"
	"github.com/jgbaldwinbrown/lscan/lscan"
	"os"
	"bufio"
	"compress/gzip"
)

type Plottable struct {
	Outprefix string
	Plottable [][]string
}

func ScanPath(path string) (*bufio.Scanner, *os.File, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	return s, r, nil
}

type InfoE struct {
	Line []string
	Gen int64
}

func ParseInfoE(line []string) (InfoE, error) {
	var ie InfoE
	if len(line) < 2 {
		return ie, errors.New("Info line too short")
	}

	ie.Line = line
	var err error
	ie.Gen, err = strconv.ParseInt(line[0], 0, 64)
	return ie, err
}

func ReadInfo(path string) ([]InfoE, error) {
	s, r, err := ScanPath(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	scanf := lscan.ByByte('\t')
	out := []InfoE{}
	s.Scan()
	for s.Scan() {
		line := []string{}
		line = lscan.SplitByFunc(line, s.Text(), scanf)
		ie, err := ParseInfoE(line)
		if err != nil {
			return out, err
		}
		out = append(out, ie)
	}
	return out, nil
}

type BedE struct {
	Chr string
	Start int64
	End int64
}

func ParseBedE(line []string) (BedE, error) {
	var b BedE
	if len(line) < 3 {
		return b, errors.New("Bed line too short")
	}

	b.Chr = line[0]

	var err error
	b.Start, err = strconv.ParseInt(line[1], 0, 64)
	if err != nil {
		return b, err
	}

	b.End, err = strconv.ParseInt(line[2], 0, 64)
	if err != nil {
		return b, err
	}

	return b, nil
}

func ReadBed(path string) ([]BedE, error) {
	s, r, err := ScanPath(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	out, err := ReadBedScanner(s)
	return out, err
}

func ReadBedScanner(s *bufio.Scanner) ([]BedE, error) {
	scanf := lscan.ByByte('\t')
	out := []BedE{}
	line := []string{}
	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), scanf)
		b, err := ParseBedE(line)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func ResizeBed(bed []BedE, size int64) []BedE {
	var out []BedE
	for _, b := range bed {
		mid := (b.End + b.Start / 2)
		b.End = mid + size/2
		b.Start = mid - size/2
		out = append(out, b)
	}
	return out
}

func WriteBed(w io.Writer, bed []BedE) {
	for _, b := range bed {
		fmt.Fprintf(w, "%v\t%v\t%v\n", b.Chr, b.Start, b.End)
	}
}

type SyncE struct {
	Chr string
	Pos int64
	Afs [][]int64
	Line []string
}

func ParseSyncChrPos(line []string) (chr string, pos int64, err error) {
	if len(line) < 2 {
		return chr, pos, errors.New("Sync line too short")
	}

	chr = line[0]
	pos, err = strconv.ParseInt(line[1], 0, 64)
	pos -= 1
	return chr, pos, err
}

func ParseSyncCol(col string) ([]int64, error) {
	af_strs := strings.Split(col, ":")
	if len(af_strs) != 5 {
		return nil, errors.New("Sync column the wrong length")
	}

	afs := make([]int64, 5)
	for i, af_str := range af_strs {
		af, err := strconv.ParseInt(af_str, 0, 64)
		if err != nil {
			return afs, err
		}
		afs[i] = af
	}
	return afs, nil
}

func ParseSyncE(line []string) (SyncE, error) {
	var err error
	s := SyncE{}
	s.Chr, s.Pos, err = ParseSyncChrPos(line)
	if err != nil {
		return s, err
	}
	for _, col := range line[3:] {
		af, err := ParseSyncCol(col)
		if err != nil {
			return s, err
		}
		s.Afs = append(s.Afs, af)
	}
	s.Line = make([]string, len(line))
	copy(s.Line, line)
	return s, nil
}

func InBed(chr string, pos int64, bed []BedE) bool {
	for _, b := range bed {
		if chr == b.Chr && pos >= b.Start && pos < b.End {
			return true
		}
	}
	return false
}

func ReadSync(path string, bed []BedE) ([]SyncE, error) {
	conn, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var g *gzip.Reader
	var r io.Reader

	re := regexp.MustCompile(`\.gz$`)
	if re.MatchString(path) {
		g, err = gzip.NewReader(conn)
		if err != nil {
			return nil, err
		}
		r = g
		defer g.Close()
	} else {
		r = conn
	}

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	scanf := lscan.ByByte('\t')
	out := []SyncE{}
	line := []string{}
	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), scanf)
		chr, pos, err := ParseSyncChrPos(line)
		if err != nil {
			return nil, err
		}
		if InBed(chr, pos, bed) {
			sy, err := ParseSyncE(line)
			if err != nil {
				return nil, err
			}
			out = append(out, sy)
		}
	}
	return out, nil
}

func PlottableLineOut(major, minor int64, chr string, pos int64, info InfoE) []string {
	out := make([]string, 10)
	total := major + minor
	majorf := float64(major) / float64(total)
	minorf := float64(minor) / float64(total)
	out[0] = chr
	out[1] = strconv.FormatInt(pos, 10)
	out[2] = out[0] + "_" + out[1]
	out[3] = strconv.FormatInt(major, 10)
	out[4] = strconv.FormatInt(minor, 10)
	out[5] = strconv.FormatFloat(majorf, 'e', -1, 64)
	out[6] = strconv.FormatFloat(minorf, 'e', -1, 64)
	out[7] = info.Line[0]
	out[8] = info.Line[1]
	out[9] = out[2] + "_" + out[8]
	return out
}

func CalcAfFreq(af []int64) []float64 {
	out := make([]float64, len(af))
	var total int64 = 0
	for _, val := range af {
		total += val
	}
	for i, val := range af {
		out[i] = float64(val) / float64(total)
	}
	return out
}

func CalcAfFreqs(afs [][]int64) [][]float64 {
	out := make([][]float64, len(afs))
	for i, af := range afs {
		out[i] = CalcAfFreq(af)
	}
	return out
}

func ColMeans(afs_f [][]float64) []float64 {
	if len(afs_f) < 1 {
		return nil
	}
	out := make([]float64, len(afs_f[0]))
	for i, _ := range out {
		for _, af_f := range afs_f {
			out[i] += af_f[i]
		}
		out[i] = out[i] / float64(len(afs_f))
	}
	return out
}

func GetMajorTwoAFs(afs [][]int64) (major, minor int) {
	freqs := CalcAfFreqs(afs)
	mean_freqs := ColMeans(freqs)

	major_val := 0.0
	minor_val := 0.0
	for i, val := range mean_freqs {
		if val > major_val {
			minor_val = major_val
			minor = major
			major_val = val
			major = i
		} else if val > minor_val {
			minor_val = val
			minor = i
		}
	}
	return major, minor
}

func ToPlottableLine(s SyncE, info []InfoE) [][]string {
	var out [][]string
	major_ind, minor_ind := GetMajorTwoAFs(s.Afs)
	for i, afset := range s.Afs {
		out = append(out, PlottableLineOut(afset[major_ind], afset[minor_ind], s.Chr, s.Pos, info[i]))
	}
	return out
}

func ToPlottable(sync []SyncE, info []InfoE) [][]string {
	var p [][]string
	for _, s := range sync {
		p = append(p, ToPlottableLine(s, info)...)
	}
	return p
}

func ToSeparatePlottables(sync []SyncE, bed []BedE, info []InfoE, outprefix string) []Plottable {
	var out []Plottable
	for _, b := range bed {
		var ss []SyncE
		bb := []BedE{b}
		for _, s := range sync {
			if InBed(s.Chr, s.Pos, bb) {
				ss = append(ss, s)
			}
		}
		p := ToPlottable(ss, info)
		pref := outprefix
		if len(p) > 0 {
			pref = fmt.Sprintf("%v_%v_%v_%v", outprefix, b.Chr, b.Start, b.End)
		}
		out = append(out, Plottable{pref, p})
	}
	return out
}

func PlottableHeader() []string {
	return []string{
		"chr",
		"pos",
		"chrpos",
		"major_c",
		"minor_c",
		"major_f",
		"minor_f",
		"gen",
		"repl",
		"chrposrepl",
	}
}

func WritePlottable(w io.Writer, plottable [][]string) {
	fmt.Fprintln(w, "chr\tpos\tchrpos\tmajor_c\tminor_c\tmajor_f\tminor_f\tgen\trepl\tchrposrepl")
	for _, p := range plottable {
		io.WriteString(w, strings.Join(p, "\t") + "\n")
	}
}

const (
	plottablesuffix string = "_plottable.txt"
	plottedsuffix string = "_plotted.pdf"
)

func WritePlottableToFile(p [][]string, outprefix string) error {
	txtout, err := os.Create(outprefix + plottablesuffix)
	if err != nil {
		return err
	}
	defer txtout.Close()
	b := bufio.NewWriter(txtout)
	defer b.Flush()

	WritePlottable(b, p)
	return nil
}

func WritePlottablesToFiles(plottables ...Plottable) error {
	for _, plottable := range plottables {
		err := WritePlottableToFile(plottable.Plottable, plottable.Outprefix)
		if err != nil {
			return err
		}
	}
	return nil
}

func PlotPlottableFile(inpath, outpath string) error {
	cmd := exec.Command("plotafs", inpath, outpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func PlotPlottable(p Plottable) error {
	err := PlotPlottableFile(p.Outprefix + plottablesuffix, p.Outprefix + plottedsuffix)
	return err
}

func PlotPlottables(ps ...Plottable) error {
	for _, p := range ps {
		err := PlotPlottable(p)
		if err != nil {
			return err
		}
	}
	return nil
}
