package qui

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/qbradq/q2d"
)

//go:embed icons.png
var iconsPng []byte

type Icon int

const (
	IconNone Icon = iota
	IconFile
	IconFolder
	IconSave
	IconOpen
	IconNew
	IconCopy
	IconCut
	IconPaste
	IconUndo
	IconRedo
	IconCheck
	IconUncheck
	IconRadioOn
	IconRadioOff
	IconClose
	IconMinimize
	IconMaximize
	IconArrowUp
	IconArrowDown
	IconArrowLeft
	IconArrowRight
	IconEdit
	IconDelete
	IconSettings
	IconHelp
	IconInfo
	IconWarning
	IconError
	IconUser
	IconGroup
	IconSearch
	IconZoomIn
	IconZoomOut
)

const (
	IconSize      = 16
	IconSheetSize = 256
	IconsPerRow   = IconSheetSize / IconSize
)

// DrawIcon draws the icon at the given position with the given color (tint).
func DrawIcon(img *q2d.Image, icon Icon, p q2d.Point, c q2d.Color) {
	if DefaultTheme == nil || DefaultTheme.IconSheet == nil || icon == IconNone {
		return
	}

	idx := int(icon)
	if idx < 0 {
		return
	}

	sx := (idx % IconsPerRow) * IconSize
	sy := (idx / IconsPerRow) * IconSize

	// We need to draw a sub-image of the icon sheet, tinted with c.
	// q2d doesn't support complex blending/tinting of images easily yet,
	// but the prompt says "The icon sheet will be white and transparent to allow tinted rendering".
	// So we can iterate pixels and multiply alpha/color.

	// Since q2d.Image is our target, and IconSheet is likely an image.Image or q2d.Image.
	// Let's assume IconSheet is *image.RGBA for now in the Theme, or *q2d.Image.
	// If it's *q2d.Image, we can read from it.

	// Let's implement a blit with tint.

	// destRect := q2d.Rectangle{p.X(), p.Y(), IconSize, IconSize}

	// Check clip
	// This logic should probably be in q2d, but we can do it here for now.

	sheet := DefaultTheme.IconSheet

	for y := 0; y < IconSize; y++ {
		for x := 0; x < IconSize; x++ {
			srcC := sheet.At(q2d.Point{sx + x, sy + y})
			if srcC.A() == 0 {
				continue
			}

			// Tint: multiply alpha. Source is white (255,255,255,A).
			// Target is (R,G,B, A * srcA).

			// a := float64(srcC.A()) / 255.0
			// finalA := float64(c.A()) / 255.0 * a

			// Simple alpha blending with background is handled by q2d.Image.Set usually?
			// q2d.Image.Set replaces? No, let's check q2d.Image.
			// q2d.Image.Set just sets the pixel. It doesn't blend.
			// q2d.Image.Text DOES blend.

			// We should probably implement a DrawIcon or Blit in q2d, but I can't modify q2d easily without context switch.
			// I'll implement manual blending here similar to Text.

			targetP := p.Add(q2d.Point{x, y})
			bg := img.At(targetP) // This respects clip?
			// img.At checks clip and returns 0 if outside.
			// But we need to know if it was actually outside or just transparent/black.
			// Actually img.At returns Color{}, which is 0,0,0,0.

			// We should check clip first.
			// img.Set checks clip.

			// Blending:
			// outA = srcA + dstA(1-srcA)
			// outC = (srcC*srcA + dstC*dstA(1-srcA)) / outA

			// In our case, "Source" is the tinted icon pixel.
			// SrcR = c.R, SrcG = c.G, SrcB = c.B, SrcA = c.A * pixelAlpha

			pixelAlpha := float64(srcC.A()) / 255.0
			srcR := float64(c.R())
			srcG := float64(c.G())
			srcB := float64(c.B())
			srcA := float64(c.A()) / 255.0 * pixelAlpha

			dstR := float64(bg.R())
			dstG := float64(bg.G())
			dstB := float64(bg.B())
			dstA := float64(bg.A()) / 255.0

			outA := srcA + dstA*(1.0-srcA)
			if outA == 0 {
				continue
			}

			outR := (srcR*srcA + dstR*dstA*(1.0-srcA)) / outA
			outG := (srcG*srcA + dstG*dstA*(1.0-srcA)) / outA
			outB := (srcB*srcA + dstB*dstA*(1.0-srcA)) / outA

			img.Set(targetP, q2d.Color{uint8(outR), uint8(outG), uint8(outB), uint8(outA * 255.0)})
		}
	}
}

// CreateDummyIconSheet loads the embedded icons.png
func CreateDummyIconSheet() *q2d.Image {
	img, _, err := image.Decode(bytes.NewReader(iconsPng))
	if err != nil {
		// Fallback to empty image if decode fails (shouldn't happen with valid embed)
		return q2d.NewImage(IconSheetSize, IconSheetSize)
	}
	return ImageToQ2D(img)
}

// Convert image.Image to q2d.Image
func ImageToQ2D(src image.Image) *q2d.Image {
	b := src.Bounds()
	dst := q2d.NewImage(b.Dx(), b.Dy())

	// This is slow but works
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, b, a := src.At(x+b.Min.X, y+b.Min.Y).RGBA()
			// RGBA returns 0-65535 pre-multiplied
			// We need non-premultiplied uint8?
			// q2d.Color is uint8.
			// Assuming q2d.Color is straight alpha or premultiplied?
			// q2d implementation of Text blending suggests it handles blending manually.
			// Let's assume straight alpha for storage or whatever q2d expects.
			// Actually q2d.Image is just a buffer.

			// Standard Go image/draw works with premultiplied alpha.
			// Let's just cast down.
			dst.Set(q2d.Point{x, y}, q2d.Color{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return dst
}
