// Stats package includes simple statistics utilities
package stats

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

// Series calculates cumulative statistics on a series of numbers
// For streamlined use, all operations are designed to be error-free, results
// may not always be useful though.
type Series[T Number] struct {
	count uint
	min   T
	max   T
	total T
}

func (s *Series[T]) Add(v T) {
	if s.count <= 0 {
		s.min = v
		s.max = v
		s.total = v
	} else {
		if s.min > v {
			s.min = v
		}
		if s.max < v {
			s.max = v
		}
		s.total += v
	}
	s.count += 1
}

func (s Series[T]) Len() uint { return s.count }
func (s Series[T]) Min() T    { return s.min }
func (s Series[T]) Max() T    { return s.max }
func (s Series[T]) Total() T  { return s.total }

func (s Series[T]) Avg() T {
	if s.count <= 0 {
		return 0
	}
	return s.total / T(s.count)
}

// Implement the fmt.Formatter interface
func (s Series[T]) Format(f fmt.State, verb rune) {
	fldFmtStr := rebuildFmtStr(f, verb)
	fmtStr := fmt.Sprintf("min: %s max: %s avg: %s", fldFmtStr, fldFmtStr, fldFmtStr)
	fmt.Fprintf(f, fmtStr, s.Min(), s.Max(), s.Avg())
}

func rebuildFmtStr(f fmt.State, verb rune) string {
	fldFmtStr := "%"
	for flag := range [...]int{'+', '=', '#', ' ', '0'} {
		if f.Flag(flag) {
			fldFmtStr += string(rune(flag))
		}
	}
	if w, ok := f.Width(); ok {
		fldFmtStr += fmt.Sprintf("%d", w)
	}
	if p, ok := f.Precision(); ok {
		if p == 0 {
			fldFmtStr += "."
		} else {
			fldFmtStr += fmt.Sprintf(".%d", p)
		}
	}
	fldFmtStr += string(verb)
	return fldFmtStr
}
