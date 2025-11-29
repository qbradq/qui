package qui

import (
	"github.com/qbradq/tremor/lib/q2d"
	"golang.org/x/image/font"
)

type MenuItem struct {
	BaseWidget
	Text   string
	Icon   Icon
	Action func()

	hovered bool
	parent  *PopupMenu
}

func NewMenuItem(text string, icon Icon, action func()) *MenuItem {
	return &MenuItem{
		Text:   text,
		Icon:   icon,
		Action: action,
	}
}

func (m *MenuItem) MinSize() Size {
	theme := m.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	width := font.MeasureString(theme.Font, m.Text).Ceil()
	width += IconSize + theme.Spacing + theme.Padding*2

	metrics := theme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()
	if IconSize > height {
		height = IconSize
	}
	height += theme.Padding * 2

	return Size{width, height}
}

func (m *MenuItem) Event(e Event) bool {
	switch evt := e.(type) {
	case MouseEvent:
		if m.Rect.Contains(evt.Pos) {
			if evt.TypeVal == EventMouseMove {
				m.hovered = true
				return true
			}
			if evt.TypeVal == EventMouseUp {
				if m.Action != nil {
					m.Action()
				}
				if m.parent != nil {
					m.parent.Close()
				}
				return true
			}
		} else {
			m.hovered = false
		}
	}
	return false
}

func (m *MenuItem) Draw(img *q2d.Image) {
	theme := m.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(m.Rect)
	defer img.PopSubImage()

	if m.hovered {
		img.Fill(theme.PrimaryColor)
	}

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	contentHeight := textHeight
	if IconSize > contentHeight {
		contentHeight = IconSize
	}

	y := (m.Rect.Height() - contentHeight) / 2
	if y < 0 {
		y = 0
	}

	x := theme.Padding
	if m.Icon != IconNone {
		iconY := y + (contentHeight-IconSize)/2
		DrawIcon(img, m.Icon, q2d.Point{x, iconY}, theme.TextColor)
		x += IconSize + theme.Spacing
	}

	textY := y + (contentHeight-textHeight)/2
	img.Text(q2d.Point{x, textY}, theme.TextColor, theme.Font, false, "%s", m.Text)
}

type PopupMenu struct {
	BaseWidget
	Items         []*MenuItem
	OnDismissFunc func()

	overlayManager OverlayManager
}

func NewPopupMenu(items ...*MenuItem) *PopupMenu {
	p := &PopupMenu{Items: items}
	for _, item := range items {
		item.parent = p
	}
	return p
}

func (p *PopupMenu) MinSize() Size {
	w, h := 0, 0
	for _, item := range p.Items {
		sz := item.MinSize()
		if sz.Width > w {
			w = sz.Width
		}
		h += sz.Height
	}
	return Size{w, h}
}

func (p *PopupMenu) Layout(available Size) Size {
	y := p.Rect.Y()
	w := p.Rect.Width() // Use set width (usually min width)

	for _, item := range p.Items {
		sz := item.MinSize()
		item.SetRect(q2d.Rectangle{p.Rect.X(), y, w, sz.Height})
		y += sz.Height
	}
	return Size{w, y - p.Rect.Y()}
}

func (p *PopupMenu) Draw(img *q2d.Image) {
	// Draw background
	theme := p.GetTheme()
	img.PushSubImage(p.Rect)
	img.Fill(theme.BackgroundColor.Darken(0.1))
	img.PopSubImage()

	for _, item := range p.Items {
		item.Draw(img)
	}

	// Draw border on top
	img.PushSubImage(p.Rect)
	img.Border(theme.BorderColor)
	img.PopSubImage()
}

func (p *PopupMenu) Event(e Event) bool {
	handled := false
	for _, item := range p.Items {
		if item.Event(e) {
			handled = true
			if e.Type() != EventMouseMove {
				return true
			}
		}
	}
	if handled {
		return true
	}
	if mouse, ok := e.(MouseEvent); ok {
		if p.Rect.Contains(mouse.Pos) {
			return true
		}
	}
	return false
}

func (p *PopupMenu) OnDismiss() {
	if p.OnDismissFunc != nil {
		p.OnDismissFunc()
	}
}

func (p *PopupMenu) SetOverlayManager(m OverlayManager) {
	p.overlayManager = m
}

func (p *PopupMenu) Close() {
	if p.overlayManager != nil {
		p.overlayManager.PopOverlay()
	}
}

func (p *PopupMenu) FindWidgetAt(pos q2d.Point) Widget {
	if !p.Rect.Contains(pos) {
		return nil
	}
	for _, item := range p.Items {
		if w := item.FindWidgetAt(pos); w != nil {
			return w
		}
	}
	return p
}

type MenuBar struct {
	BaseWidget
	Menus []*MenuItem // Using MenuItem as top level menu headers
	// We need to know which menu is open
	OpenMenuIndex int
	// And the popup menus associated
	Popups []*PopupMenu

	OverlayManager OverlayManager
}

func NewMenuBar() *MenuBar {
	return &MenuBar{
		Menus:         make([]*MenuItem, 0),
		OpenMenuIndex: -1,
		Popups:        make([]*PopupMenu, 0),
	}
}

func (m *MenuBar) AddMenu(title string, popup *PopupMenu) {
	item := NewMenuItem(title, IconNone, nil)
	m.Menus = append(m.Menus, item)
	m.Popups = append(m.Popups, popup)

	// Hook action to toggle menu
	index := len(m.Menus) - 1
	item.Action = func() {
		if m.OpenMenuIndex == index {
			// Close it
			if m.OverlayManager != nil {
				m.OverlayManager.PopOverlay()
			}
			m.OpenMenuIndex = -1
		} else {
			// If another menu is open, close it first?
			// If OverlayManager has a stack, we might need to pop the previous one if it's a sibling menu.
			// But Master handles "click outside" which might have already popped it?
			// If we clicked on this menu item, we are outside the previous popup.
			// So Master might have popped it.
			// But we need to be sure.
			// If m.OpenMenuIndex != -1, it means we think a menu is open.
			// If Master popped it, our OnDismiss should have cleared OpenMenuIndex.
			// So if OpenMenuIndex is still set, we should pop it.
			if m.OpenMenuIndex != -1 {
				if m.OverlayManager != nil {
					m.OverlayManager.PopOverlay()
				}
			}

			m.OpenMenuIndex = index
			if m.OverlayManager != nil {
				// Position popup
				menuItem := m.Menus[index]
				popup := m.Popups[index]
				sz := popup.MinSize()
				popup.SetRect(q2d.Rectangle{
					menuItem.Rect.X(),
					menuItem.Rect.Y() + menuItem.Rect.Height(),
					sz.Width,
					sz.Height,
				})

				// Set dismiss callback
				popup.OnDismissFunc = func() {
					m.OpenMenuIndex = -1
				}

				m.OverlayManager.PushOverlay(popup)
			}
		}
	}
}

func (m *MenuBar) MinSize() Size {
	theme := m.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	metrics := theme.Font.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding*2

	// Width is 100% usually, but min width is sum of items
	w := 0
	for _, item := range m.Menus {
		w += item.MinSize().Width
	}
	return Size{w, height}
}

func (m *MenuBar) Layout(available Size) Size {
	x := m.Rect.X()
	h := m.Rect.Height()

	for _, item := range m.Menus {
		sz := item.MinSize()
		item.SetRect(q2d.Rectangle{x, m.Rect.Y(), sz.Width, h})
		x += sz.Width
	}

	// Layout open popup if any
	// Popups are now handled by OverlayManager, so we don't layout them here.
	// But we need to ensure they are positioned correctly when opened.
	// (Handled in Action)

	return available
}

func (m *MenuBar) Draw(img *q2d.Image) {
	theme := m.GetTheme()
	img.PushSubImage(m.Rect)
	img.Fill(theme.BackgroundColor.Darken(0.1))
	img.PopSubImage()

	for _, item := range m.Menus {
		item.Draw(img)
	}
}

func (m *MenuBar) Event(e Event) bool {
	// Popups are handled by Master now.
	// We just handle menu bar clicks.

	handled := false
	for _, item := range m.Menus {
		if item.Event(e) {
			handled = true
			if e.Type() != EventMouseMove {
				return true
			}
		}
	}
	return handled
}

func (m *MenuBar) FindWidgetAt(pos q2d.Point) Widget {
	if !m.Rect.Contains(pos) {
		return nil
	}
	for _, item := range m.Menus {
		if w := item.FindWidgetAt(pos); w != nil {
			return w
		}
	}
	return m
}
