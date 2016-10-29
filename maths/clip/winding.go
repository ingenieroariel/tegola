package clip

type WindingOrder bool

const (
	Clockwise        = WindingOrder(false)
	CounterClockwise = WindingOrder(true)
)

func (w WindingOrder) String() string {
	if w {
		return "Counter Clockwise"
	}
	return "Clockwise"
}

// IsClockwise returns weather the winding order is clockwise.
func (w WindingOrder) IsClockwise() bool {
	return w == Clockwise
}

// IsCounterClockwise returns weather the winding order is counter clockwise.
func (w WindingOrder) IsCounterClockwise() bool {
	return w == CounterClockwise
}
