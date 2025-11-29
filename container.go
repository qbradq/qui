package qui

import (
	"github.com/qbradq/q2d"
)

type LayoutDirection int

const (
	LayoutVertical LayoutDirection = iota
	LayoutHorizontal
)

type Container struct {
	BaseWidget
	Children  []Widget
	Direction LayoutDirection
	Stretch   bool
}

func NewContainer(dir LayoutDirection, children ...Widget) *Container {
	return &Container{
		Children:  children,
		Direction: dir,
	}
}

func (c *Container) Add(w Widget) {
	c.Children = append(c.Children, w)
}

func (c *Container) MinSize() Size {
	theme := c.GetTheme()
	spacing := 5
	if theme != nil {
		spacing = theme.Spacing
	}

	w, h := 0, 0
	for i, child := range c.Children {
		sz := child.MinSize()
		if c.Direction == LayoutVertical {
			if sz.Width > w {
				w = sz.Width
			}
			h += sz.Height
			if i < len(c.Children)-1 {
				h += spacing
			}
		} else {
			if sz.Height > h {
				h = sz.Height
			}
			w += sz.Width
			if i < len(c.Children)-1 {
				w += spacing
			}
		}
	}
	return Size{w, h}
}

func (c *Container) Layout(available Size) Size {
	theme := c.GetTheme()
	spacing := 5
	if theme != nil {
		spacing = theme.Spacing
	}

	// Calculate total spacing
	totalSpacing := 0
	if len(c.Children) > 1 {
		totalSpacing = spacing * (len(c.Children) - 1)
	}

	// First pass: Calculate space taken by non-fill widgets
	usedSpace := 0
	fillCount := 0
	for _, child := range c.Children {
		if child.IsFill() {
			fillCount++
		} else {
			sz := child.MinSize()
			if c.Direction == LayoutVertical {
				usedSpace += sz.Height
			} else {
				usedSpace += sz.Width
			}
		}
	}

	// Calculate space per fill widget
	fillSpace := 0
	if fillCount > 0 {
		availableSpace := 0
		if c.Direction == LayoutVertical {
			availableSpace = available.Height - totalSpacing - usedSpace
		} else {
			availableSpace = available.Width - totalSpacing - usedSpace
		}
		if availableSpace > 0 {
			fillSpace = availableSpace / fillCount
		}
	}

	// Second pass: Layout
	x, y := c.Rect.X(), c.Rect.Y()
	totalW, totalH := 0, 0

	for i, child := range c.Children {
		sz := child.MinSize()

		// Determine size for this child
		childW, childH := sz.Width, sz.Height

		if c.Direction == LayoutVertical {
			if c.Stretch {
				childW = available.Width
			}
			if child.IsFill() {
				childH = fillSpace
				// If it's the last fill widget, give it any remainder?
				// For simplicity, just use integer division result.
				// Or better: ensure we fill exactly?
				// Let's stick to simple for now.
				if childH < sz.Height {
					childH = sz.Height // Respect min size?
				}
			}

			child.SetRect(q2d.Rectangle{x, y, childW, childH})
			child.Layout(Size{childW, childH})

			y += childH + spacing
			if childW > totalW {
				totalW = childW
			}
			totalH += childH
		} else {
			if c.Stretch {
				childH = available.Height
			}
			if child.IsFill() {
				childW = fillSpace
				if childW < sz.Width {
					childW = sz.Width
				}
			}

			child.SetRect(q2d.Rectangle{x, y, childW, childH})
			child.Layout(Size{childW, childH})

			x += childW + spacing
			if childH > totalH {
				totalH = childH
			}
			totalW += childW
		}

		if i < len(c.Children)-1 {
			if c.Direction == LayoutVertical {
				totalH += spacing
			} else {
				totalW += spacing
			}
		}
	}

	return Size{totalW, totalH}
}

func (c *Container) Draw(img *q2d.Image) {
	for _, child := range c.Children {
		child.Draw(img)
	}
}

func (c *Container) Event(e Event) bool {
	for _, child := range c.Children {
		if child.Event(e) {
			return true
		}
	}
	return false
}

func (c *Container) FindWidgetAt(pos q2d.Point) Widget {
	if !c.Rect.Contains(pos) {
		return nil
	}
	// Check children in reverse order (top-most first)
	for i := len(c.Children) - 1; i >= 0; i-- {
		child := c.Children[i]
		if w := child.FindWidgetAt(pos); w != nil {
			return w
		}
	}
	return c
}
