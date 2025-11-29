# QUI - Retained Mode GUI Library for Go

This package was created with Gemini CLI and Antigravity using the Gemini 2.5
and 3 models. It is intended to be used by other LLMs.

QUI is a simple, modular, retained-mode GUI library for Go, built on top of
`q2d`. It provides a set of common user interface widgets and a theming system.

## Table of Contents

- [Installation](#installation)
- [Getting Started](#getting-started)
- [Widgets](#widgets)
  - [Label](#label)
  - [Button](#button)
  - [Container](#container)
  - [Entry](#entry)
  - [List](#list)
  - [Window](#window)
- [Theming](#theming)

## Installation

```bash
go get github.com/qbradq/qui
```

## Getting Started

To use QUI, you need to initialize a `Master` controller with a root widget and
a theme. You then integrate this into your `q2d` application loop.

```go
package main

import (
	"github.com/qbradq/qui"
	"github.com/qbradq/q2d"
)

func main() {
	// 1. Initialize Theme
	qui.InitTheme(q2d.FontThin)

	// 2. Create Root Widget
	root := qui.NewContainer(qui.LayoutVertical,
		qui.NewLabel("Hello, QUI!"),
		qui.NewButton("Click Me", func() {
			println("Button Clicked!")
		}),
	)

	// 3. Create Master Controller
	master := qui.NewMaster(root, qui.DefaultTheme)

	// 4. In your Game Loop:
	// - Pass events to master.Event(e)
	// - Call master.Layout(size) when window resizes
	// - Call master.Draw(img) to render
}
```

## Widgets

### Label

Displays text with an optional icon.

**Code Example:**

```go
label := qui.NewLabel("This is a label")
labelWithIcon := qui.NewLabel("Label with Icon")
labelWithIcon.Icon = qui.IconFile // Set an icon
```

**Expected Result:**

- `label`: Renders the text "This is a label".
- `labelWithIcon`: Renders a file icon followed by the text "Label with Icon".

### Button

A clickable button with text and an optional icon.

**Code Example:**

```go
btn := qui.NewButton("Submit", func() {
    println("Form submitted")
})
```

**Expected Result:**

- Renders a button with the text "Submit".
- When clicked, "Form submitted" is printed to the console.
- Visual feedback (color change) on hover and press.

### Container

Layouts children either vertically or horizontally.

**Code Example:**

```go
// Vertical Container
vBox := qui.NewContainer(qui.LayoutVertical,
    qui.NewLabel("Top"),
    qui.NewLabel("Bottom"),
)

// Horizontal Container
hBox := qui.NewContainer(qui.LayoutHorizontal,
    qui.NewButton("Left", nil),
    qui.NewButton("Right", nil),
)
```

**Expected Result:**

- `vBox`: Renders "Top" label above "Bottom" label.
- `hBox`: Renders "Left" button to the left of "Right" button.

### Entry

A single-line text input field. Supports different types like Text, Password,
Integer, and Float.

**Code Example:**

```go
// Text Entry
nameEntry := qui.NewEntry("John Doe", qui.EntryText)

// Password Entry
passEntry := qui.NewEntry("", qui.EntryPassword)

// Integer Entry
ageEntry := qui.NewEntry("25", qui.EntryInteger)
```

**Expected Result:**

- `nameEntry`: Input field showing "John Doe". Allows editing text.
- `passEntry`: Input field showing asterisks (\*). Hides actual text.
- `ageEntry`: Input field showing "25". Only accepts numeric input.

### List

A scrollable vertical list of items.

**Code Example:**

```go
items := []qui.ListItem{
    {Text: "Item 1"},
    {Text: "Item 2", Icon: qui.IconFolder},
    {Text: "Item 3"},
}

list := qui.NewList(items, func(index int) {
    println("Selected item:", index)
})
```

**Expected Result:**

- Renders a list with 3 items.
- "Item 2" displays a folder icon.
- Clicking an item highlights it and prints its index.
- Shows a scrollbar if items exceed the visible area.

### Window

A draggable container with a header, title, and close button.

**Code Example:**

```go
content := qui.NewLabel("Window Content")
win := qui.NewWindow("My Window", content)

// Optional: Handle close
win.OnClose = func() {
    println("Window closed")
}
```

**Expected Result:**

- Renders a window frame with title "My Window".
- Contains the "Window Content" label.
- Can be dragged by the header.
- Clicking the 'X' button triggers `OnClose`.

## Theming

QUI supports custom themes. You can generate a theme from a base color or
manually set all colors.

**Code Example:**

```go
// Generate a theme from a base color (e.g., Blue)
baseColor := q2d.Color{0, 0, 255, 255}
newTheme := qui.GenerateThemeFromColor(baseColor, basicfont.Face7x13)

// Apply to a specific widget
widget.SetTheme(newTheme)

// Or set as global default (before Master creation or in Master)
qui.DefaultTheme = newTheme
```

**Expected Result:**

- The UI elements will use a color palette derived from the base blue color.
- Backgrounds will be dark blue, text white, and primary accents bright blue.
