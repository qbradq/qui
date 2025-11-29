package qui

import (
	"github.com/qbradq/q2d"
)

type Panel struct {
	BaseWidget
	Content Widget
}

func NewPanel(content Widget) *Panel {
	return &Panel{
		Content: content,
	}
}

func (p *Panel) MinSize() Size {
	theme := p.GetTheme()
	if theme == nil {
		return Size{0, 0}
	}
	if p.Content == nil {
		return Size{theme.Padding * 2, theme.Padding * 2}
	}
	sz := p.Content.MinSize()
	return Size{sz.Width + theme.Padding*2, sz.Height + theme.Padding*2}
}

func (p *Panel) Layout(available Size) Size {
	theme := p.GetTheme()
	if theme == nil {
		return available
	}

	if p.Content != nil {
		p.Content.SetRect(q2d.Rectangle{
			p.Rect.X() + theme.Padding,
			p.Rect.Y() + theme.Padding,
			p.Rect.Width() - theme.Padding*2,
			p.Rect.Height() - theme.Padding*2,
		})
		p.Content.Layout(Size{
			p.Rect.Width() - theme.Padding*2,
			p.Rect.Height() - theme.Padding*2,
		})
	}
	return available
}

func (p *Panel) Draw(img *q2d.Image) {
	theme := p.GetTheme()
	if theme == nil {
		return
	}

	img.PushSubImage(p.Rect)
	defer img.PopSubImage()

	if theme.BackgroundColor.A() > 0 {
		img.Fill(theme.BackgroundColor)
	}
	if theme.BorderColor.A() > 0 {
		img.Border(theme.BorderColor)
	}

	if p.Content != nil {
		p.Content.Draw(img)
	}
}

func (p *Panel) Event(evt Event) bool {
	if p.Content != nil {
		return p.Content.Event(evt)
	}
	return false
}

func (p *Panel) FindWidgetAt(pos q2d.Point) Widget {
	if !p.Rect.Contains(pos) {
		return nil
	}
	if p.Content != nil {
		if found := p.Content.FindWidgetAt(pos); found != nil {
			return found
		}
	}
	return p
}
