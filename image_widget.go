package qui

import (
	"image"

	"github.com/qbradq/q2d"
)

type ImageWidget struct {
	BaseWidget
	Img image.Image
}

func NewImageWidget(img image.Image) *ImageWidget {
	return &ImageWidget{Img: img}
}

func (w *ImageWidget) MinSize() Size {
	if w.Img == nil {
		return Size{0, 0}
	}
	b := w.Img.Bounds()
	return Size{b.Dx(), b.Dy()}
}

func (w *ImageWidget) Draw(img *q2d.Image) {
	if w.Img == nil {
		return
	}

	img.PushSubImage(w.Rect)
	defer img.PopSubImage()

	// We need to draw the image.
	// Since q2d doesn't have DrawImage, we do pixel copy.
	// This is slow but functional.

	b := w.Img.Bounds()
	wW, wH := w.Rect.Width(), w.Rect.Height()

	// Draw 1:1 for now, clipped to widget rect

	for y := 0; y < wH; y++ {
		if y >= b.Dy() {
			break
		}
		for x := 0; x < wW; x++ {
			if x >= b.Dx() {
				break
			}

			r, g, b, a := w.Img.At(x+b.Min.X, y+b.Min.Y).RGBA()
			// RGBA is 0-65535 premultiplied.
			// q2d.Color expects uint8 non-premultiplied? Or whatever.
			// Let's assume standard conversion.

			// If alpha is 0, skip?
			if a == 0 {
				continue
			}

			c := q2d.Color{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			img.Set(q2d.Point{x, y}, c)
		}
	}
}
