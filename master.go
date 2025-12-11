package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Master struct {
	Root Widget

	// Overlay state
	Overlays []Widget

	HoveredWidget Widget
	FocusedWidget Focusable
	MousePos      q2d.Point

	// Theme
	Theme *Theme
}

func NewMaster(root Widget, theme *Theme) *Master {
	return &Master{
		Root:     root,
		Theme:    theme,
		Overlays: make([]Widget, 0),
	}
}

func (m *Master) PushOverlay(w Widget) {
	if mo, ok := w.(ManagedOverlay); ok {
		mo.SetOverlayManager(m)
	}
	m.Overlays = append(m.Overlays, w)
}

func (m *Master) PopOverlay() {
	if len(m.Overlays) > 0 {
		overlay := m.Overlays[len(m.Overlays)-1]
		m.Overlays = m.Overlays[:len(m.Overlays)-1]
		if d, ok := overlay.(Dismissable); ok {
			d.OnDismiss()
		}
	}
}

func (m *Master) Layout(size Size) {
	if m.Theme != nil {
		DefaultTheme = m.Theme
	}
	if m.Root != nil {
		m.Root.SetRect(q2d.Rectangle{0, 0, size.Width, size.Height})
		m.Root.Layout(size)
	}
	// Layout overlays
	for _, overlay := range m.Overlays {
		rect := overlay.GetRect()
		overlay.Layout(Size{rect.Width(), rect.Height()})
	}
}

func (m *Master) Event(e Event) bool {
	// Update MousePos
	if mouse, ok := e.(MouseEvent); ok {
		m.MousePos = mouse.Pos

		if mouse.TypeVal == EventMouseDown {
			// Handle Focus
			var target Widget
			// Check overlays
			for i := len(m.Overlays) - 1; i >= 0; i-- {
				if w := m.Overlays[i].FindWidgetAt(mouse.Pos); w != nil {
					target = w
					break
				}
			}
			// Check root
			if target == nil && m.Root != nil {
				target = m.Root.FindWidgetAt(mouse.Pos)
			}

			if target != nil {
				if f, ok := target.(Focusable); ok {
					if m.FocusedWidget != f {
						if m.FocusedWidget != nil {
							m.FocusedWidget.Unfocus()
						}
						f.Focus()
						m.FocusedWidget = f
					}
				} else {
					// Clicked non-focusable
					if m.FocusedWidget != nil {
						m.FocusedWidget.Unfocus()
						m.FocusedWidget = nil
					}
				}
			} else {
				// Clicked nothing
				if m.FocusedWidget != nil {
					m.FocusedWidget.Unfocus()
					m.FocusedWidget = nil
				}
			}
		}
	}

	// Handle Keyboard/Text events via FocusedWidget
	switch e.(type) {
	case KeyEvent, TextInputEvent:
		if m.FocusedWidget != nil {
			// We need to cast Focusable back to Widget to call Event?
			// Focusable interface doesn't have Event().
			// But all widgets have Event().
			// We can assume FocusedWidget is a Widget.
			if w, ok := m.FocusedWidget.(Widget); ok {
				if w.Event(e) {
					return true
				}
			}
		}
	case ScrollEvent:
		// Route to widget under mouse
		var target Widget
		// Check overlays
		for i := len(m.Overlays) - 1; i >= 0; i-- {
			if w := m.Overlays[i].FindWidgetAt(m.MousePos); w != nil {
				target = w
				break
			}
		}
		// Check root
		if target == nil && m.Root != nil {
			target = m.Root.FindWidgetAt(m.MousePos)
		}

		if target != nil {
			if target.Event(e) {
				return true
			}
		}
	}

	// 1. Handle Overlays (Top to Bottom)
	for i := len(m.Overlays) - 1; i >= 0; i-- {
		overlay := m.Overlays[i]
		if overlay.Event(e) {
			return true
		}
		// If click outside top overlay, should we close it?
		if mouse, ok := e.(MouseEvent); ok && mouse.TypeVal == EventMouseDown {
			if !overlay.GetRect().Contains(mouse.Pos) {
				m.PopOverlay()
				return true
			}
		}
	}

	// 2. Handle Root
	if m.Root != nil {
		if m.Root.Event(e) {
			return true
		}
	}

	return false
}

func (m *Master) Draw(img *q2d.Image) {
	if m.Theme != nil {
		DefaultTheme = m.Theme // Set global default theme for convenience
	}

	// Draw Root
	if m.Root != nil {
		m.Root.Draw(img)
	}

	// Draw Overlays (Bottom to Top)
	for _, overlay := range m.Overlays {
		overlay.Draw(img)
	}

	// Draw Tooltip
	m.UpdateHover(m.MousePos)
	if m.HoveredWidget != nil {
		text := m.HoveredWidget.GetTooltip()
		if text != "" {
			m.drawTooltip(img, text)
		}
	}
}

func (m *Master) UpdateHover(p q2d.Point) {
	m.HoveredWidget = nil

	// Check overlays first (top to bottom)
	for i := len(m.Overlays) - 1; i >= 0; i-- {
		overlay := m.Overlays[i]
		if found := overlay.FindWidgetAt(p); found != nil {
			m.HoveredWidget = found
			return
		}
	}

	// Check Root
	if m.Root != nil {
		m.HoveredWidget = m.Root.FindWidgetAt(p)
	}
}

func (m *Master) drawTooltip(img *q2d.Image, text string) {
	if m.Theme == nil || m.Theme.Font == nil {
		return
	}

	// Calculate size
	width := font.MeasureString(m.Theme.Font, text).Ceil() + m.Theme.Padding.Left + m.Theme.Padding.Right
	metrics := m.Theme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil() + m.Theme.Padding.Top + m.Theme.Padding.Bottom

	// Position: near mouse
	x := m.MousePos.X() + 10
	y := m.MousePos.Y() + 10

	// Clamp to screen
	if x+width > img.Rect.Width() {
		x = img.Rect.Width() - width
	}
	if y+height > img.Rect.Height() {
		y = img.Rect.Height() - height
	}

	rect := q2d.Rectangle{x, y, width, height}

	img.PushSubImage(rect)
	img.Fill(m.Theme.BackgroundColor)
	img.Border(m.Theme.BorderColor)

	textY := (height - (metrics.Ascent + metrics.Descent).Ceil()) / 2
	img.Text(q2d.Point{m.Theme.Padding.Left, textY}, m.Theme.TextColor, m.Theme.Font, false, "%s", text)
	img.PopSubImage()
}
