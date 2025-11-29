package qui

import (
	"github.com/qbradq/q2d"
	"golang.org/x/image/font"
)

type Tab struct {
	Title   string
	Icon    Icon
	Content Widget
}

type TabContainer struct {
	BaseWidget
	Tabs      []Tab
	ActiveTab int
}

func NewTabContainer(tabs ...Tab) *TabContainer {
	t := &TabContainer{
		Tabs:      tabs,
		ActiveTab: 0,
	}
	t.Fill = true
	return t
}

func (t *TabContainer) MinSize() Size {
	theme := t.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	metrics := theme.Font.Metrics()
	headerHeight := (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding*2
	if IconSize+theme.Padding*2 > headerHeight {
		headerHeight = IconSize + theme.Padding*2
	}

	maxW, maxH := 0, 0
	for _, tab := range t.Tabs {
		if tab.Content != nil {
			sz := tab.Content.MinSize()
			if sz.Width > maxW {
				maxW = sz.Width
			}
			if sz.Height > maxH {
				maxH = sz.Height
			}
		}
	}

	// Width must also accommodate tabs
	tabsW := 0
	for _, tab := range t.Tabs {
		w := font.MeasureString(theme.Font, tab.Title).Ceil() + theme.Padding*2
		if tab.Icon != IconNone {
			w += IconSize + theme.Spacing
		}
		tabsW += w
	}
	if tabsW > maxW {
		maxW = tabsW
	}

	return Size{maxW, maxH + headerHeight}
}

func (t *TabContainer) Layout(available Size) Size {
	theme := t.GetTheme()
	if theme == nil || theme.Font == nil {
		return Size{0, 0}
	}
	metrics := theme.Font.Metrics()
	headerHeight := (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding*2
	if IconSize+theme.Padding*2 > headerHeight {
		headerHeight = IconSize + theme.Padding*2
	}

	contentAvailable := Size{available.Width, available.Height - headerHeight}

	if t.ActiveTab >= 0 && t.ActiveTab < len(t.Tabs) {
		child := t.Tabs[t.ActiveTab].Content
		if child != nil {
			// Position child below header
			child.SetRect(q2d.Rectangle{t.Rect.X(), t.Rect.Y() + headerHeight, contentAvailable.Width, contentAvailable.Height})
			child.Layout(contentAvailable)
		}
	}

	return available
}

func (t *TabContainer) Event(evt Event) bool {
	theme := t.GetTheme()
	if theme == nil || theme.Font == nil {
		return false
	}
	metrics := theme.Font.Metrics()
	headerHeight := (metrics.Ascent + metrics.Descent).Ceil() + theme.Padding*2
	if IconSize+theme.Padding*2 > headerHeight {
		headerHeight = IconSize + theme.Padding*2
	}

	switch event := evt.(type) {
	case MouseEvent:
		if t.Rect.Contains(event.Pos) {
			relY := event.Pos.Y() - t.Rect.Y()
			if relY < headerHeight {
				// Header click
				if event.TypeVal == EventMouseDown {
					// Find which tab
					x := 0
					relX := event.Pos.X() - t.Rect.X()
					for i, tab := range t.Tabs {
						w := font.MeasureString(theme.Font, tab.Title).Ceil() + theme.Padding*2
						if tab.Icon != IconNone {
							w += IconSize + theme.Spacing
						}

						if relX >= x && relX < x+w {
							t.ActiveTab = i
							return true
						}
						x += w
					}
				}
				return true // Consume header events
			} else {
				// Pass to content
				if t.ActiveTab >= 0 && t.ActiveTab < len(t.Tabs) {
					if t.Tabs[t.ActiveTab].Content != nil {
						return t.Tabs[t.ActiveTab].Content.Event(evt)
					}
				}
			}
		}
	default:
		// Pass non-mouse events to content
		if t.ActiveTab >= 0 && t.ActiveTab < len(t.Tabs) {
			if t.Tabs[t.ActiveTab].Content != nil {
				return t.Tabs[t.ActiveTab].Content.Event(evt)
			}
		}
	}
	return false
}

func (t *TabContainer) FindWidgetAt(pos q2d.Point) Widget {
	if !t.Rect.Contains(pos) {
		return nil
	}

	// Check content
	if t.ActiveTab >= 0 && t.ActiveTab < len(t.Tabs) {
		if t.Tabs[t.ActiveTab].Content != nil {
			if found := t.Tabs[t.ActiveTab].Content.FindWidgetAt(pos); found != nil {
				return found
			}
		}
	}

	// If on header, return t
	return t
}

func (t *TabContainer) Draw(img *q2d.Image) {
	theme := t.GetTheme()
	if theme == nil {
		return
	}

	metrics := theme.Font.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()
	headerHeight := textHeight + theme.Padding*2
	if IconSize+theme.Padding*2 > headerHeight {
		headerHeight = IconSize + theme.Padding*2
	}

	// Draw Header
	img.PushSubImage(t.Rect)

	// Draw Header Background
	headerRect := q2d.Rectangle{0, 0, t.Rect.Width(), headerHeight}
	img.PushSubImage(headerRect)
	img.Fill(theme.BackgroundColor.Darken(0.2))

	x := 0
	for i, tab := range t.Tabs {
		w := font.MeasureString(theme.Font, tab.Title).Ceil() + theme.Padding*2
		if tab.Icon != IconNone {
			w += IconSize + theme.Spacing
		}

		bg := theme.ButtonColor
		if i == t.ActiveTab {
			bg = theme.BackgroundColor
		}

		// Draw tab rect
		tabRect := q2d.Rectangle{x, 0, w, headerHeight}
		img.PushSubImage(tabRect)
		img.Fill(bg)
		img.Border(theme.BorderColor)

		contentX := theme.Padding
		if tab.Icon != IconNone {
			iconY := (headerHeight - IconSize) / 2
			DrawIcon(img, tab.Icon, q2d.Point{contentX, iconY}, theme.TextColor)
			contentX += IconSize + theme.Spacing
		}

		textY := (headerHeight - textHeight) / 2
		img.Text(q2d.Point{contentX, textY}, theme.TextColor, theme.Font, false, "%s", tab.Title)
		img.PopSubImage()

		x += w
	}
	img.PopSubImage() // Header Background
	img.PopSubImage() // t.Rect

	// Draw Content Background
	contentRect := q2d.Rectangle{0, headerHeight, t.Rect.Width(), t.Rect.Height() - headerHeight}
	img.PushSubImage(t.Rect)
	img.PushSubImage(contentRect)
	img.Fill(theme.BackgroundColor.Darken(0.2))
	img.PopSubImage()
	img.PopSubImage()

	// Draw Content
	if t.ActiveTab >= 0 && t.ActiveTab < len(t.Tabs) {
		if t.Tabs[t.ActiveTab].Content != nil {
			t.Tabs[t.ActiveTab].Content.Draw(img)
		}
	}
}
