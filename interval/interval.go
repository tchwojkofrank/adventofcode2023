package interval

import "errors"

type Interval struct {
	Start int
	End   int
}

func (a Interval) Intersection(b Interval) (Interval, error) {
	var result Interval
	var err error
	if a.End < b.Start || b.End < a.Start {
		result = a
		err = errors.New("No intersection")
		return result, err
	}
	return Interval{max(a.Start, b.Start), min(a.End, b.End)}, nil
}

func (a Interval) Union(b Interval) (Interval, error) {
	if a.End < b.Start || b.End < a.Start {
		return Interval{}, errors.New("No union")
	}

	return Interval{min(a.Start, b.Start), max(a.End, b.End)}, nil
}

func (a Interval) Contains(b Interval) bool {
	return a.Start <= b.Start && b.End <= a.End
}

func (a Interval) Equals(b Interval) bool {
	return a.Start == b.Start && a.End == b.End
}

func (a Interval) Intersects(b Interval) bool {
	return a.End >= b.Start && b.End >= a.Start
}

func (a Interval) String() string {
	return "[" + string(a.Start) + "," + string(a.End) + "]"
}

func (a Interval) Less(b Interval) bool {
	return a.Start < b.Start
}
