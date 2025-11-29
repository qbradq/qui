package qui

import (
	"github.com/qbradq/tremor/lib/q2d"
)

type ScrolledContainer struct {
	BaseWidget
	Content Widget
	ScrollX int
	ScrollY int
	Width   int // Optional fixed width
	Height  int // Optional fixed height

	draggingX    bool
	draggingY    bool
	dragStart    q2d.Point
	startScrollX int
	startScrollY int
}

func NewScrolledContainer(content Widget) *ScrolledContainer {
	return &ScrolledContainer{
		Content: content,
	}
}

func (s *ScrolledContainer) MinSize() Size {
	w := 100
	h := 100
	if s.Width > 0 {
		w = s.Width
	} else {
		// If no fixed width, try to fit content width + scrollbar
		if s.Content != nil {
			sz := s.Content.MinSize()
			w = sz.Width + 10
		} else {
			w = 10
		}
	}

	if s.Height > 0 {
		h = s.Height
	}

	return Size{w, h}
}

func (s *ScrolledContainer) Layout(available Size) Size {
	viewportW := s.Rect.Width()
	viewportH := s.Rect.Height()

	var contentMin Size
	if s.Content != nil {
		contentMin = s.Content.MinSize()
	}

	// Check if we need scrollbars
	// Initial check
	needV := contentMin.Height > viewportH
	needH := contentMin.Width > viewportW

	barSize := 10

	// If we need V, width reduces
	if needV {
		if contentMin.Width > viewportW-barSize {
			needH = true
		}
	}
	// If we need H, height reduces
	if needH {
		if contentMin.Height > viewportH-barSize {
			needV = true
		}
	}

	// Adjust viewport for content layout
	if needV {
		viewportW -= barSize
	}
	if needH {
		viewportH -= barSize
	}

	// Clamp Scroll
	maxScrollX := contentMin.Width - viewportW
	if maxScrollX < 0 {
		maxScrollX = 0
	}
	if s.ScrollX > maxScrollX {
		s.ScrollX = maxScrollX
	}

	maxScrollY := contentMin.Height - viewportH
	if maxScrollY < 0 {
		maxScrollY = 0
	}
	if s.ScrollY > maxScrollY {
		s.ScrollY = maxScrollY
	}

	// Calculate content size
	contentW := contentMin.Width
	if contentW < viewportW {
		contentW = viewportW
	}
	contentH := contentMin.Height
	if contentH < viewportH {
		contentH = viewportH
	}

	if s.Content != nil {
		s.Content.SetRect(q2d.Rectangle{
			s.Rect.X() - s.ScrollX,
			s.Rect.Y() - s.ScrollY,
			contentW,
			contentH,
		})
		s.Content.Layout(Size{contentW, contentH})
	}

	return available
}

func (s *ScrolledContainer) Event(evt Event) bool {
	barSize := 10

	var contentMin Size
	if s.Content != nil {
		contentMin = s.Content.MinSize()
	}
	viewportW := s.Rect.Width()
	viewportH := s.Rect.Height()

	// Recalculate needed scrollbars to know active areas
	needV := contentMin.Height > viewportH
	needH := contentMin.Width > viewportW
	if needV && contentMin.Width > viewportW-barSize {
		needH = true
	}
	if needH && contentMin.Height > viewportH-barSize {
		needV = true
	}

	effViewportW := viewportW
	effViewportH := viewportH
	if needV {
		effViewportW -= barSize
	}
	if needH {
		effViewportH -= barSize
	}

	maxScrollX := contentMin.Width - effViewportW
	if maxScrollX < 0 {
		maxScrollX = 0
	}
	maxScrollY := contentMin.Height - effViewportH
	if maxScrollY < 0 {
		maxScrollY = 0
	}

	switch event := evt.(type) {
	case ScrollEvent:
		s.ScrollY -= int(event.DeltaY * 20)
		if s.ScrollY < 0 {
			s.ScrollY = 0
		}
		if s.ScrollY > maxScrollY {
			s.ScrollY = maxScrollY
		}
		return true

	case MouseEvent:
		if event.TypeVal == EventMouseDown {
			if s.Rect.Contains(event.Pos) {
				// Check scrollbars
				relX := event.Pos.X() - s.Rect.X()
				relY := event.Pos.Y() - s.Rect.Y()

				// Vertical Scrollbar
				if needV && relX >= viewportW-barSize {
					// In vertical bar area
					trackH := effViewportH
					thumbH := int(float64(trackH) * float64(trackH) / float64(contentMin.Height))
					if thumbH < 20 {
						thumbH = 20
					}
					thumbY := int(float64(s.ScrollY) / float64(maxScrollY) * float64(trackH-thumbH))

					if relY >= thumbY && relY < thumbY+thumbH {
						s.draggingY = true
						s.dragStart = event.Pos
						s.startScrollY = s.ScrollY
						return true
					}
				}

				// Horizontal Scrollbar
				if needH && relY >= viewportH-barSize {
					// In horizontal bar area
					trackW := effViewportW
					thumbW := int(float64(trackW) * float64(trackW) / float64(contentMin.Width))
					if thumbW < 20 {
						thumbW = 20
					}
					thumbX := int(float64(s.ScrollX) / float64(maxScrollX) * float64(trackW-thumbW))

					if relX >= thumbX && relX < thumbX+thumbW {
						s.draggingX = true
						s.dragStart = event.Pos
						s.startScrollX = s.ScrollX
						return true
					}
				}

				// Pass to content
				if s.Content != nil {
					return s.Content.Event(evt)
				}
				return false
			}
		} else if event.TypeVal == EventMouseUp {
			s.draggingX = false
			s.draggingY = false
			if s.Rect.Contains(event.Pos) {
				if s.Content != nil {
					return s.Content.Event(evt)
				}
				return false
			}
		} else if event.TypeVal == EventMouseMove {
			if s.draggingY && maxScrollY > 0 {
				deltaY := event.Pos.Y() - s.dragStart.Y()
				trackH := effViewportH
				thumbH := int(float64(trackH) * float64(trackH) / float64(contentMin.Height))
				if thumbH < 20 {
					thumbH = 20
				}
				// deltaY corresponds to how much scroll?
				// thumb moves (trackH - thumbH) for maxScrollY
				scrollDelta := int(float64(deltaY) * float64(maxScrollY) / float64(trackH-thumbH))
				s.ScrollY = s.startScrollY + scrollDelta
				if s.ScrollY < 0 {
					s.ScrollY = 0
				}
				if s.ScrollY > maxScrollY {
					s.ScrollY = maxScrollY
				}
				return true
			}
			if s.draggingX && maxScrollX > 0 {
				deltaX := event.Pos.X() - s.dragStart.X()
				trackW := effViewportW
				thumbW := int(float64(trackW) * float64(trackW) / float64(contentMin.Width))
				if thumbW < 20 {
					thumbW = 20
				}
				scrollDelta := int(float64(deltaX) * float64(maxScrollX) / float64(trackW-thumbW))
				s.ScrollX = s.startScrollX + scrollDelta
				if s.ScrollX < 0 {
					s.ScrollX = 0
				}
				if s.ScrollX > maxScrollX {
					s.ScrollX = maxScrollX
				}
				return true
			}

			if s.Rect.Contains(event.Pos) {
				if s.Content != nil {
					s.Content.Event(evt)
				}
				// Always return true if hovering?
				// Maybe not, but we want to track hover.
				return true
			}
		}
	}
	if s.Content != nil {
		return s.Content.Event(evt)
	}
	return false
}

func (s *ScrolledContainer) FindWidgetAt(pos q2d.Point) Widget {
	if !s.Rect.Contains(pos) {
		return nil
	}
	// Check content (which is clipped, so we should check clip?)
	// Content rect is offset.
	// If pos is inside container rect, we check content.
	// But content might be larger.
	// If we are over scrollbars, we should return ScrolledContainer (to handle drag).
	// But FindWidgetAt is for tooltips.
	// If we hover scrollbar, maybe no tooltip?
	// Or ScrolledContainer tooltip?

	// Check if over scrollbars
	// We need to recalculate scrollbar areas or store them.
	// For simplicity, let's check content first.
	// But content is offset.
	// FindWidgetAt expects absolute pos.
	// Content has absolute pos set during layout.

	if s.Content != nil {
		if found := s.Content.FindWidgetAt(pos); found != nil {
			// But wait, content is clipped by ScrolledContainer rect.
			// If found widget is outside ScrolledContainer rect (but inside Content rect),
			// it shouldn't be hovered.
			// So we must check intersection.
			// We already checked s.Rect.Contains(pos).
			return found
		}
	}

	return s
}

func (s *ScrolledContainer) Draw(img *q2d.Image) {
	img.PushClip(s.Rect)
	defer img.PopClip()

	if s.Content != nil {
		s.Content.Draw(img)
	}

	var contentMin Size
	if s.Content != nil {
		contentMin = s.Content.MinSize()
	}
	viewportW := s.Rect.Width()
	viewportH := s.Rect.Height()
	barSize := 10

	needV := contentMin.Height > viewportH
	needH := contentMin.Width > viewportW
	if needV && contentMin.Width > viewportW-barSize {
		needH = true
	}
	if needH && contentMin.Height > viewportH-barSize {
		needV = true
	}

	effViewportW := viewportW
	effViewportH := viewportH
	if needV {
		effViewportW -= barSize
	}
	if needH {
		effViewportH -= barSize
	}

	img.PushSubImage(s.Rect)
	defer img.PopSubImage()

	// Draw Vertical Scrollbar
	if needV {
		trackH := effViewportH
		thumbH := int(float64(trackH) * float64(trackH) / float64(contentMin.Height))
		if thumbH < 20 {
			thumbH = 20
		}
		maxScrollY := contentMin.Height - effViewportH
		thumbY := 0
		if maxScrollY > 0 {
			thumbY = int(float64(s.ScrollY) / float64(maxScrollY) * float64(trackH-thumbH))
		}

		theme := s.GetTheme()
		if theme == nil {
			theme = DefaultTheme
		}

		// Track
		img.PushSubImage(q2d.Rectangle{viewportW - barSize, 0, barSize, trackH})
		img.Fill(theme.BackgroundColor.Lighten(0.1))
		img.PopSubImage()

		// Thumb
		img.PushSubImage(q2d.Rectangle{viewportW - barSize, thumbY, barSize, thumbH})
		img.Fill(theme.BorderColor)
		img.PopSubImage()
	}

	// Draw Horizontal Scrollbar
	if needH {
		trackW := effViewportW
		thumbW := int(float64(trackW) * float64(trackW) / float64(contentMin.Width))
		if thumbW < 20 {
			thumbW = 20
		}
		maxScrollX := contentMin.Width - effViewportW
		thumbX := 0
		if maxScrollX > 0 {
			thumbX = int(float64(s.ScrollX) / float64(maxScrollX) * float64(trackW-thumbW))
		}

		theme := s.GetTheme()
		if theme == nil {
			theme = DefaultTheme
		}

		// Track
		img.PushSubImage(q2d.Rectangle{0, viewportH - barSize, trackW, barSize})
		img.Fill(theme.BackgroundColor.Lighten(0.1))
		img.PopSubImage()

		// Thumb
		img.PushSubImage(q2d.Rectangle{thumbX, viewportH - barSize, thumbW, barSize})
		img.Fill(theme.BorderColor)
		img.PopSubImage()
	}

	// Corner
	if needV && needH {
		theme := s.GetTheme()
		if theme == nil {
			theme = DefaultTheme
		}
		img.PushSubImage(q2d.Rectangle{viewportW - barSize, viewportH - barSize, barSize, barSize})
		img.Fill(theme.BackgroundColor.Darken(0.3))
		img.PopSubImage()
	}
}
