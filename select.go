package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Select struct {
	BaseWidget
	Items         []ListItem
	SelectedIndex int
	OnSelect      func(index int)

	expanded     bool
	hoveredIndex int

	OverlayManager OverlayManager
}

func NewSelect(items []ListItem, onSelect func(index int)) *Select {
	return &Select{
		Items:         items,
		SelectedIndex: -1,
		OnSelect:      onSelect,
		hoveredIndex:  -1,
	}
}

func (s *Select) MinSize() Size {
	theme := s.GetTheme()
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

	for _, item := range s.Items {
		w := font.MeasureString(theme.Font, item.Text).Ceil()
		if item.Icon != IconNone {
			w += IconSize + theme.Spacing
		}
		if w > maxWidth {
			maxWidth = w
		}
	}

	h := lineHeight + theme.Padding*2
	// No longer include list height in MinSize as it's an overlay

	return Size{maxWidth + theme.Padding*4, h}
}

func (s *Select) Event(evt Event) bool {
	theme := s.GetTheme()
	if theme == nil || theme.Font == nil {
		return false
	}

	switch event := evt.(type) {
	case MouseEvent:
		if s.Rect.Contains(event.Pos) {
			if event.TypeVal == EventMouseDown {
				// Toggle
				if s.expanded {
					// Close
					if s.OverlayManager != nil {
						s.OverlayManager.PopOverlay()
					}
					s.expanded = false
				} else {
					// Open
					s.expanded = true
					if s.OverlayManager != nil {
						// Create and push overlay
						list := &SelectList{
							Select: s,
							OnDismissFunc: func() {
								s.expanded = false
							},
						}
						// Calculate size and pos
						metrics := theme.Font.Metrics()
						lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
						if IconSize > lineHeight {
							lineHeight = IconSize
						}
						lineHeight += 2

						h := len(s.Items)*lineHeight + 2
						if len(s.Items) > 5 {
							h = 5*lineHeight + 2
						}

						list.SetRect(q2d.Rectangle{
							s.Rect.X(),
							s.Rect.Y() + s.Rect.Height(),
							s.Rect.Width(),
							h,
						})

						s.OverlayManager.PushOverlay(list)
					}
				}
				return true
			}
		}
	}
	return false
}

func (s *Select) Draw(img *q2d.Image) {
	theme := s.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(s.Rect)
	defer img.PopSubImage()

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineHeight := textHeight
	if IconSize > lineHeight {
		lineHeight = IconSize
	}
	lineHeight += 2
	headerHeight := lineHeight + theme.Padding*2

	// Draw Header
	img.Fill(theme.ButtonColor)
	img.Border(theme.BorderColor)

	text := "Select..."
	icon := IconNone
	if s.SelectedIndex >= 0 && s.SelectedIndex < len(s.Items) {
		text = s.Items[s.SelectedIndex].Text
		icon = s.Items[s.SelectedIndex].Icon
	}

	x := theme.Padding
	if icon != IconNone {
		iconY := (headerHeight - IconSize) / 2
		DrawIcon(img, icon, q2d.Point{x, iconY}, theme.TextColor)
		x += IconSize + theme.Spacing
	}

	textY := (headerHeight - textHeight) / 2
	img.Text(q2d.Point{x, textY}, theme.TextColor, theme.Font, false, "%s", text)
}

type SelectList struct {
	BaseWidget
	Select        *Select
	OnDismissFunc func()
	ScrollOffset  int

	dragging     bool
	dragStart    q2d.Point
	startScrollY int
}

func (l *SelectList) OnDismiss() {
	if l.OnDismissFunc != nil {
		l.OnDismissFunc()
	}
}

func (l *SelectList) Event(evt Event) bool {
	theme := l.Select.GetTheme()
	if theme == nil {
		return false
	}
	metrics := theme.Font.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()
	if IconSize > lineHeight {
		lineHeight = IconSize
	}
	lineHeight += 2

	barSize := 10
	contentHeight := len(l.Select.Items) * lineHeight
	viewportHeight := l.Rect.Height() - 2
	maxScroll := contentHeight - viewportHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch event := evt.(type) {
	case ScrollEvent:
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

				if index >= 0 && index < len(l.Select.Items) {
					l.Select.SelectedIndex = index
					if l.Select.OnSelect != nil {
						l.Select.OnSelect(index)
					}
					// Close overlay
					if l.Select.OverlayManager != nil {
						l.Select.OverlayManager.PopOverlay()
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
					l.Select.hoveredIndex = -1
					return true
				}

				relY := event.Pos.Y() - l.Rect.Y() - 1 + l.ScrollOffset
				index := relY / lineHeight

				if index >= 0 && index < len(l.Select.Items) {
					l.Select.hoveredIndex = index
					return true
				} else {
					l.Select.hoveredIndex = -1
				}
			} else {
				l.Select.hoveredIndex = -1
			}
		}
	}
	return false
}

func (l *SelectList) Draw(img *q2d.Image) {
	theme := l.Select.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(l.Rect)
	defer img.PopSubImage()

	img.Fill(theme.BackgroundColor)
	img.Border(theme.BorderColor)

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	lineHeight := textHeight
	if IconSize > lineHeight {
		lineHeight = IconSize
	}
	lineHeight += 2

	barSize := 10
	viewportHeight := l.Rect.Height() - 2
	contentHeight := len(l.Select.Items) * lineHeight
	maxScroll := contentHeight - viewportHeight

	// Clip content to inside border (excluding scrollbar if needed)
	contentWidth := l.Rect.Width() - 2
	if maxScroll > 0 {
		contentWidth -= barSize
	}

	// Clip content
	contentRect := q2d.Rectangle{1, 1, contentWidth, viewportHeight}
	img.PushSubImage(contentRect)

	startIdx := l.ScrollOffset / lineHeight
	endIdx := (l.ScrollOffset + contentRect.Height() + lineHeight - 1) / lineHeight
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(l.Select.Items) {
		endIdx = len(l.Select.Items)
	}

	for i := startIdx; i < endIdx; i++ {
		item := l.Select.Items[i]
		y := i*lineHeight - l.ScrollOffset

		bg := theme.BackgroundColor
		if i == l.Select.SelectedIndex {
			bg = theme.PrimaryColor
		} else if i == l.Select.hoveredIndex {
			bg = theme.SecondaryColor
		}

		if i == l.Select.SelectedIndex || i == l.Select.hoveredIndex {
			img.PushSubImage(q2d.Rectangle{0, y, contentRect.Width(), lineHeight})
			img.Fill(bg)
			img.PopSubImage()
		}

		itemX := theme.Padding
		if item.Icon != IconNone {
			iconY := y + (lineHeight-IconSize)/2
			DrawIcon(img, item.Icon, q2d.Point{itemX, iconY}, theme.TextColor)
			itemX += IconSize + theme.Spacing
		}

		itemTextY := y + (lineHeight-textHeight)/2
		img.Text(q2d.Point{itemX, itemTextY}, theme.TextColor, theme.Font, false, "%s", item.Text)
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
