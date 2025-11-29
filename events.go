package qui

import "github.com/qbradq/q2d"

type EventType int

const (
	EventMouseMove EventType = iota
	EventMouseDown
	EventMouseUp
	EventKeyDown
	EventKeyUp
	EventTextInput
	EventScroll
)

type Event interface {
	Type() EventType
}

type MouseEvent struct {
	TypeVal EventType
	Pos     q2d.Point
	Button  int // 0: Left, 1: Right, 2: Middle
}

func (e MouseEvent) Type() EventType { return e.TypeVal }

type ScrollEvent struct {
	TypeVal EventType
	DeltaX  float64
	DeltaY  float64
}

func (e ScrollEvent) Type() EventType { return e.TypeVal }

type KeyEvent struct {
	TypeVal EventType
	Key     int // Ebiten key code or similar
}

func (e KeyEvent) Type() EventType { return e.TypeVal }

type TextInputEvent struct {
	Text string
}

func (e TextInputEvent) Type() EventType { return EventTextInput }

const (
	KeyBackspace = 8
	KeyEnter     = 13
	KeyHome      = 36
	KeyLeft      = 37
	KeyUp        = 38
	KeyRight     = 39
	KeyDown      = 40
	KeyEnd       = 35
	KeyDelete    = 46
)
