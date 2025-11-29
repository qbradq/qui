package qui

import (
	"github.com/qbradq/tremor/lib/q2d"
	"golang.org/x/image/font"
)

type List struct {
	BaseWidget
	Items         []ListItem
	SelectedIndex int
	OnSelect      func(index int)

	hoveredIndex int
	focused      bool
	ScrollOffset int

	dragging     bool
	dragStart    q2d.Point
	startScrollY int
}

func NewList(items []ListItem, onSelect func(index int)) *List {
	return &List{
		Items:         items,
		SelectedIndex: -1,
		OnSelect:      onSelect,
		hoveredIndex:  -1,
	}
}

func (l *List) MinSize() Size {
	theme := l.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	maxWidth := 0
	metrics := theme.Font.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	if IconSize > lineHeight {
		lineHeight = IconSize
	}
	lineHeight += 2

	for _, item := range l.Items {
		w := font.MeasureString(theme.Font, item.Text).Ceil()
		if item.Icon != IconNone {
			w += IconSize + theme.Spacing
		}
		if w > maxWidth {
			maxWidth = w
		}
	}

	// Default to 5 items tall
	h := 5*lineHeight + 2
	if len(l.Items) < 5 {
		h = len(l.Items)*lineHeight + 2
	}

	return Size{maxWidth + theme.Padding*2 + 10, h} // Add space for scrollbar
}

func (l *List) Event(evt Event) bool {
	theme := l.GetTheme()
	if theme == nil || theme.Font == nil {
		return false
	}
	metrics := theme.Font.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	if IconSize > lineHeight {
		lineHeight = IconSize
	}
	lineHeight += 2

	barSize := 10
	contentHeight := len(l.Items) * lineHeight
	viewportHeight := l.Rect.Height() - 2
	maxScroll := contentHeight - viewportHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch event := evt.(type) {
	case ScrollEvent:
		// Scroll event doesn't have pos usually, but we assume it's global or focused?
		// Actually ScrollEvent in qui.go has no position. Master handles routing?
		// Master.Event: "If a widget is focused, it receives keyboard events. Mouse events are sent to the widget under the cursor."
		// Scroll is usually sent to widget under cursor.
		// Let's assume Master sends it if we are hovered.
		// But wait, Master.Event implementation for ScrollEvent?
		// We need to check Master.Event.
		// Assuming we get it:

		l.ScrollOffset -= int(event.DeltaY * float64(lineHeight))
		if l.ScrollOffset < 0 {
			l.ScrollOffset = 0
		}
		if l.ScrollOffset > maxScroll {
			l.ScrollOffset = maxScroll
		}
		return true
	case MouseEvent:
		if event.TypeVal == EventMouseDown {
			if l.Rect.Contains(event.Pos) {
				// Check scrollbar
				if maxScroll > 0 && event.Pos.X() >= l.Rect.X()+l.Rect.Width()-barSize {
					l.dragging = true
					l.dragStart = event.Pos
					l.startScrollY = l.ScrollOffset
					return true
				}

				relY := event.Pos.Y() - l.Rect.Y() - 1 + l.ScrollOffset
				index := relY / lineHeight

				if index >= 0 && index < len(l.Items) {
					l.SelectedIndex = index
					if l.OnSelect != nil {
						l.OnSelect(index)
					}
					return true
				}
			}
		} else if event.TypeVal == EventMouseUp {
			l.dragging = false
			if l.Rect.Contains(event.Pos) {
				return true
			}
		} else if event.TypeVal == EventMouseMove {
			if l.dragging && maxScroll > 0 {
				deltaY := event.Pos.Y() - l.dragStart.Y()
				trackH := viewportHeight
				thumbH := int(float64(trackH) * float64(trackH) / float64(contentHeight))
				if thumbH < 20 {
					thumbH = 20
				}

				scrollDelta := int(float64(deltaY) * float64(maxScroll) / float64(trackH-thumbH))
				l.ScrollOffset = l.startScrollY + scrollDelta
				if l.ScrollOffset < 0 {
					l.ScrollOffset = 0
				}
				if l.ScrollOffset > maxScroll {
					l.ScrollOffset = maxScroll
				}
				return true
			}

			if l.Rect.Contains(event.Pos) {
				// Check if over scrollbar
				if maxScroll > 0 && event.Pos.X() >= l.Rect.X()+l.Rect.Width()-barSize {
					l.hoveredIndex = -1
					return true
				}

				relY := event.Pos.Y() - l.Rect.Y() - 1 + l.ScrollOffset
				index := relY / lineHeight

				if index >= 0 && index < len(l.Items) {
					l.hoveredIndex = index
					return true
				} else {
					l.hoveredIndex = -1
				}
			} else {
				l.hoveredIndex = -1
			}
		}
	case KeyEvent:
		if l.focused && event.TypeVal == EventKeyDown {
			if event.Key == KeyUp { // Up
				if l.SelectedIndex > 0 {
					l.SelectedIndex--
					// Scroll to show
					itemTop := l.SelectedIndex * lineHeight
					if itemTop < l.ScrollOffset {
						l.ScrollOffset = itemTop
					}

					if l.OnSelect != nil {
						l.OnSelect(l.SelectedIndex)
					}
					return true
				}
			} else if event.Key == KeyDown { // Down
				if l.SelectedIndex < len(l.Items)-1 {
					l.SelectedIndex++
					// Scroll to show
					itemBottom := (l.SelectedIndex + 1) * lineHeight
					if itemBottom > l.ScrollOffset+l.Rect.Height()-2 {
						l.ScrollOffset = itemBottom - (l.Rect.Height() - 2)
					}

					if l.OnSelect != nil {
						l.OnSelect(l.SelectedIndex)
					}
					return true
				}
			}
		}
	}
	return false
}

func (l *List) Focus() {
	l.focused = true
}

func (l *List) Unfocus() {
	l.focused = false
}

func (l *List) FindWidgetAt(pos q2d.Point) Widget {
	if l.Rect.Contains(pos) {
		return l
	}
	return nil
}

func (l *List) Draw(img *q2d.Image) {
	theme := l.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(l.Rect)
	defer img.PopSubImage()

	img.Fill(theme.BackgroundColor)

	borderColor := theme.BorderColor
	if l.focused {
		borderColor = theme.PrimaryColor
	}
	img.Border(borderColor)

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineHeight := textHeight
	if IconSize > lineHeight {
		lineHeight = IconSize
	}
	lineHeight += 2

	barSize := 10
	viewportHeight := l.Rect.Height() - 2
	contentHeight := len(l.Items) * lineHeight
	maxScroll := contentHeight - viewportHeight

	// Clip content to inside border (excluding scrollbar if needed)
	contentWidth := l.Rect.Width() - 2
	if maxScroll > 0 {
		contentWidth -= barSize
	}

	contentRect := q2d.Rectangle{1, 1, contentWidth, viewportHeight}
	img.PushSubImage(contentRect)

	startIdx := l.ScrollOffset / lineHeight
	endIdx := (l.ScrollOffset + contentRect.Height() + lineHeight - 1) / lineHeight
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(l.Items) {
		endIdx = len(l.Items)
	}

	for i := startIdx; i < endIdx; i++ {
		item := l.Items[i]
		y := i*lineHeight - l.ScrollOffset

		bg := theme.BackgroundColor
		if i == l.SelectedIndex {
			bg = theme.PrimaryColor
		} else if i == l.hoveredIndex {
			bg = theme.SecondaryColor
		}

		if i == l.SelectedIndex || i == l.hoveredIndex {
			img.PushSubImage(q2d.Rectangle{0, y, contentRect.Width(), lineHeight})
			img.Fill(bg)
			img.PopSubImage()
		}

		x := theme.Padding
		if item.Icon != IconNone {
			iconY := y + (lineHeight-IconSize)/2
			DrawIcon(img, item.Icon, q2d.Point{x, iconY}, theme.TextColor)
			x += IconSize + theme.Spacing
		}

		textY := y + (lineHeight-textHeight)/2
		img.Text(q2d.Point{x, textY}, theme.TextColor, theme.Font, false, "%s", item.Text)
	}
	img.PopSubImage() // Pop content clip

	// Draw Scrollbar
	if maxScroll > 0 {
		trackH := viewportHeight
		thumbH := int(float64(trackH) * float64(trackH) / float64(contentHeight))
		if thumbH < 20 {
			thumbH = 20
		}
		thumbY := int(float64(l.ScrollOffset) / float64(maxScroll) * float64(trackH-thumbH))

		// Track
		img.PushSubImage(q2d.Rectangle{l.Rect.Width() - barSize - 1, 1, barSize, trackH})
		img.Fill(theme.BackgroundColor.Lighten(0.1)) // Lighter track
		img.PopSubImage()

		// Thumb
		img.PushSubImage(q2d.Rectangle{l.Rect.Width() - barSize - 1, 1 + thumbY, barSize, thumbH})
		img.Fill(theme.BorderColor)
		img.PopSubImage()
	}
}
