# QUI Design Document

## Overview
QUI is a retained-mode GUI library for Go, built on top of `q2d`. It provides a set of common UI widgets and a theming
system.

## Architecture

QUI follows a retained-mode architecture where the UI is built as a tree of widgets. The root of the tree is managed by
a `Master` object (or `Context`) which handles input events and rendering.

### Widget Interface

All UI elements implement the `Widget` interface:

```go
type Widget interface {
    // Layout calculates the size and position of the widget and its children.
    // It receives the available space constraints.
    Layout(constraints Size) Size

    // Draw renders the widget onto the q2d.Image.
    // It receives the drawing context (offset, clip).
    Draw(img *q2d.Image, ctx Context)

    // Event handles input events (mouse, keyboard).
    // Returns true if the event was consumed.
    Event(e Event) bool

    // MinSize returns the minimum size required by the widget.
    MinSize() Size
}
```

### Theme

The look and feel of the UI is controlled by a `Theme` struct:

```go
type Theme struct {
    BackgroundColor q2d.Color
    TextColor       q2d.Color
    PrimaryColor    q2d.Color
    SecondaryColor  q2d.Color
    Font            font.Face
    Spacing         int
    // ... other properties
}
```

## Widgets

The following widgets will be implemented:

- **Label**: Displays text with optional Icon.
- **Button**: Clickable button with text and optional Icon.
- **List**: Vertical list of items with optional Icons.
- **Select**: Dropdown selection list with optional Icons.
- **Entry**: Single-line text input (Text, Password, Integer, Float).
- **TextArea**: Multi-line text input.
- **TabContainer**: Container with tabs (Title + Icon) to switch between views.
- **ScrolledContainer**: Container with scrollbars for content larger than the view.
- **Image**: Displays an image.
- **Checkbox**: Binary toggle button.
- **RadioButton**: Single selection from a group.
- **Window**: Draggable container with header and frame.
- **MenuItem**: Menu item with text, icon, and action.
- **PopupMenu**: Vertical list of menu items.
- **MenuBar**: Horizontal menu bar.

## Implementation Plan

1.  **Core Infrastructure**:
    -   Define `Widget` interface.
    -   Define `Theme` struct.
    -   Create `Master` / `Context` for managing the UI tree.
    -   Implement event handling types.

2.  **Basic Widgets**:
    -   `Label`
    -   `Button`
    -   `Container` (Box layout)

3.  **Input Widgets**:
    -   `Entry`
    -   `TextArea`

4.  **Complex Widgets**:
    -   `List`
    -   `Select`
    -   `TabContainer`
    -   `ScrolledContainer`

## Progress

- [x] Core Infrastructure
- [x] Label
- [x] Button
- [x] Container
- [x] Entry
- [x] TextArea
- [x] List
- [x] Select
- [x] TabContainer
- [x] ScrolledContainer
