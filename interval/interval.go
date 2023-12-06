package interval

import "errors"

type Interval struct {
	start int
	end   int
}

func (a Interval) intersect(b Interval) (Interval, error) {
	var result Interval
	var err error
	if a.end < b.start || b.end < a.start {
		result = a
		err = errors.New("No intersection")
		return result, err
	}
	return Interval{max(a.start, b.start), min(a.end, b.end)}, nil
}

func (a Interval) union(b Interval) (Interval, error) {
	if a.end < b.start || b.end < a.start {
		return Interval{}, errors.New("No union")
	}

	return Interval{min(a.start, b.start), max(a.end, b.end)}, nil
}

func (a Interval) contains(b Interval) bool {
	return a.start <= b.start && b.end <= a.end
}

func (a Interval) equals(b Interval) bool {
	return a.start == b.start && a.end == b.end
}

func (a Interval) intersects(b Interval) bool {
	return a.end >= b.start && b.end >= a.start
}

func (a Interval) String() string {
	return "[" + string(a.start) + "," + string(a.end) + "]"
}
