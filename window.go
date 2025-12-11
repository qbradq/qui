package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Window struct {
	BaseWidget
	Title      string
	Content    Widget
	ShowHeader bool
	ShowFrame  bool
	Closable   bool
	OnClose    func()

	dragging  bool
	dragStart q2d.Point

	closeBtn *Button

	overlayManager OverlayManager
}

func NewWindow(title string, content Widget) *Window {
	w := &Window{
		Title:      title,
		Content:    content,
		ShowHeader: true,
		ShowFrame:  true,
		Closable:   true,
	}

	w.closeBtn = NewButton("", func() {
		if w.overlayManager != nil {
			w.overlayManager.PopOverlay()
		} else if w.OnClose != nil {
			// Fallback if not managed
			w.OnClose()
		}
	})
	w.closeBtn.Icon = IconClose

	return w
}

func (w *Window) MinSize() Size {
	var contentSz Size
	if w.Content != nil {
		contentSz = w.Content.MinSize()
	}
	width := contentSz.Width
	height := contentSz.Height

	if w.ShowFrame {
		width += 4 // 2px each side
		height += 4
	}

	if w.ShowHeader {
		theme := w.GetTheme()
		if theme != nil && theme.Font != nil {
			metrics := theme.Font.Metrics()
			headerH := (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding.Top + theme.Padding.Bottom
			if IconSize+theme.Padding.Top+theme.Padding.Bottom > headerH {
				headerH = IconSize + theme.Padding.Top + theme.Padding.Bottom
			}
			height += headerH

			// Min width for title + close btn
			titleW := font.MeasureString(theme.Font, w.Title).Ceil() + theme.Padding.Left + theme.Padding.Right
			if w.Closable {
				titleW += IconSize + theme.Padding.Left + theme.Padding.Right // Close btn size approx
			}
			if titleW > width {
				width = titleW
			}
		}
	}

	return Size{width, height}
}

func (w *Window) Layout(available Size) Size {
	// Window determines its own size usually, or takes available?
	// If it's a window, it usually has a set size.
	// But Layout is called to tell it how much space it HAS.
	// It should layout its content.

	// We assume w.Rect is already set by parent or self.

	contentRect := w.Rect

	if w.ShowFrame {
		contentRect = q2d.Rectangle{
			contentRect.X() + 2,
			contentRect.Y() + 2,
			contentRect.Width() - 4,
			contentRect.Height() - 4,
		}
	}

	if w.ShowHeader {
		headerH := 0
		theme := w.GetTheme()
		if theme != nil && theme.Font != nil {
			metrics := theme.Font.Metrics()
			headerH = (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding.Top + theme.Padding.Bottom
			if IconSize+theme.Padding.Top+theme.Padding.Bottom > headerH {
				headerH = IconSize + theme.Padding.Top + theme.Padding.Bottom
			}
		}

		// Layout Close Button
		if w.Closable {
			// Header takes space from top.

			// Let's adjust contentRect for header
			// Header is at top of frame (inside frame)

			// Actual header rect
			headerRect := q2d.Rectangle{
				contentRect.X(),
				contentRect.Y(), // Start at top of content area (inside frame)
				contentRect.Width(),
				headerH,
			}

			// Close btn inside header
			w.closeBtn.SetRect(q2d.Rectangle{
				headerRect.X() + headerRect.Width() - headerH, // Square button
				headerRect.Y(),
				headerH,
				headerH,
			})

			contentRect = q2d.Rectangle{
				contentRect.X(),
				contentRect.Y() + headerH,
				contentRect.Width(),
				contentRect.Height() - headerH,
			}
		} else {
			// Just adjust for header height
			contentRect = q2d.Rectangle{
				contentRect.X(),
				contentRect.Y() + headerH,
				contentRect.Width(),
				contentRect.Height() - headerH,
			}
		}
	}

	if w.Content != nil {
		w.Content.SetRect(contentRect)
		w.Content.Layout(Size{contentRect.Width(), contentRect.Height()})
	}

	return Size{w.Rect.Width(), w.Rect.Height()}
}

func (w *Window) Event(evt Event) bool {
	// Handle dragging
	switch event := evt.(type) {
	case MouseEvent:
		if w.dragging {
			if event.TypeVal == EventMouseUp {
				w.dragging = false
				return true
			}
			if event.TypeVal == EventMouseMove {
				delta := event.Pos.Sub(w.dragStart)
				w.Rect = q2d.Rectangle{
					w.Rect.X() + delta.X(),
					w.Rect.Y() + delta.Y(),
					w.Rect.Width(),
					w.Rect.Height(),
				}
				w.dragStart = event.Pos
				// Need to re-layout children because absolute positions changed
				w.Layout(Size{w.Rect.Width(), w.Rect.Height()})
				return true
			}
		}

		if w.Rect.Contains(event.Pos) {
			// Check header for drag
			if w.ShowHeader {
				headerH := 0
				theme := w.GetTheme()
				if theme != nil && theme.Font != nil {
					metrics := theme.Font.Metrics()
					headerH = (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding.Top + theme.Padding.Bottom
					if IconSize+theme.Padding.Top+theme.Padding.Bottom > headerH {
						headerH = IconSize + theme.Padding.Top + theme.Padding.Bottom
					}
				}

				frameOffset := 0
				if w.ShowFrame {
					frameOffset = 2
				}

				headerRect := q2d.Rectangle{
					w.Rect.X() + frameOffset,
					w.Rect.Y() + frameOffset,
					w.Rect.Width() - frameOffset*2,
					headerH,
				}

				if headerRect.Contains(event.Pos) {
					// Check close button
					if w.Closable && w.closeBtn.Event(evt) {
						return true
					}

					if event.TypeVal == EventMouseDown {
						w.dragging = true
						w.dragStart = event.Pos
						return true
					}
				}
			}

			// Pass to content
			if w.Content != nil {
				if w.Content.Event(evt) {
					return true
				}
			}
			return true
		}
	}

	return false
}

func (w *Window) FindWidgetAt(pos q2d.Point) Widget {
	if !w.Rect.Contains(pos) {
		return nil
	}

	// Check Close Button
	if w.ShowHeader && w.closeBtn != nil {
		if found := w.closeBtn.FindWidgetAt(pos); found != nil {
			return found
		}
	}

	// Check Content
	if w.Content != nil {
		if found := w.Content.FindWidgetAt(pos); found != nil {
			return found
		}
	}

	return w
}

func (w *Window) Draw(img *q2d.Image) {
	theme := w.GetTheme()
	if theme == nil {
		return
	}

	// Draw Frame
	if w.ShowFrame {
		img.PushSubImage(w.Rect)
		img.Fill(theme.BorderColor) // Frame color
		// Fill inside with background
		inner := q2d.Rectangle{2, 2, w.Rect.Width() - 4, w.Rect.Height() - 4}
		img.PushSubImage(inner)
		img.Fill(theme.BackgroundColor)
		img.PopSubImage()
		img.PopSubImage()
	} else {
		img.PushSubImage(w.Rect)
		img.Fill(theme.BackgroundColor)
		img.PopSubImage()
	}

	// Draw Header
	if w.ShowHeader {
		headerH := 0
		if theme.Font != nil {
			metrics := theme.Font.Metrics()
			headerH = (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding.Top + theme.Padding.Bottom
			if IconSize+theme.Padding.Top+theme.Padding.Bottom > headerH {
				headerH = IconSize + theme.Padding.Top + theme.Padding.Bottom
			}
		}

		frameOffset := 0
		if w.ShowFrame {
			frameOffset = 2
		}

		headerRect := q2d.Rectangle{
			w.Rect.X() + frameOffset,
			w.Rect.Y() + frameOffset,
			w.Rect.Width() - frameOffset*2,
			headerH,
		}

		img.PushSubImage(headerRect)
		img.Fill(theme.PrimaryColor) // Header color

		// Title
		metrics := theme.Font.Metrics()
		textHeight := (metrics.Ascent + metrics.Descent).Ceil()
		textY := (headerH - textHeight) / 2
		img.Text(q2d.Point{theme.Padding.Left, textY}, theme.TextColor, theme.Font, false, "%s", w.Title)

		img.PopSubImage()

		if w.Closable {
			w.closeBtn.Draw(img)
		}
	}

	if w.Content != nil {
		w.Content.Draw(img)
	}
}

func (w *Window) SetOverlayManager(m OverlayManager) {
	w.overlayManager = m
}

func (w *Window) OnDismiss() {
	if w.OnClose != nil {
		w.OnClose()
	}
}
