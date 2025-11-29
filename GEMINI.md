# QUI Project Context

## Overview
**QUI** is a retained-mode GUI library for Go, built on top of `q2d` (likely a 2D graphics library, `github.com/qbradq/q2d`). It provides a set of common user interface widgets and a theming system, designed to be simple and modular.

The architecture follows a standard retained-mode pattern where the UI is constructed as a tree of `Widget` objects. A `Master` object manages this tree, handling input events, layout calculations, and rendering.

## Architecture

### Core Components
*   **Master (`master.go`):** The root controller. It holds the root widget, manages overlays (popups, menus), handles global input events, and coordinates the draw cycle.
*   **Widget Interface (`qui.go`):** The contract that all UI elements must fulfill. Key methods include:
    *   `Layout(available Size) Size`: Calculates size/position.
    *   `Draw(img *q2d.Image)`: Renders the widget.
    *   `Event(e Event) bool`: Handles input.
*   **Theme (`theme.go`):** Defines the visual style (colors, fonts, spacing) for widgets.
*   **Events (`events.go`):** Defines input events like `MouseEvent`, `KeyEvent`, `TextInputEvent`.

### Key Interfaces (from `qui.go`)
*   `Widget`: The primary interface for all UI components.
*   `OverlayManager`: Interface for widgets that need to manage overlays (like a `Select` dropdown).
*   `Focusable`: Interface for widgets that can accept keyboard focus.
*   `ManagedOverlay` & `Dismissable`: For controlling overlay lifecycles.

## Project Structure
*   `qui.go`: Core interfaces and `BaseWidget` implementation.
*   `master.go`: Main loop and coordination logic.
*   `*.go` (other): Individual widget implementations (e.g., `button.go`, `label.go`, `window.go`).
*   `DESIGN.md`: Project design goals and progress tracker.
*   `go.mod`: Go module definition, showing a dependency on `github.com/qbradq/q2d`.

## Development Conventions
*   **Widget Implementation:** Most widgets embed `BaseWidget` to handle common properties like `Rect`, `Tooltip`, and `Theme`.
*   **Layout:** Layout is top-down. The parent passes constraints to the child in `Layout`, and the child returns its actual used size.
*   **Drawing:** Drawing is done onto a `*q2d.Image` context.
*   **Events:** Events bubble or are routed by the `Master`. Returning `true` from `Event` consumes it.

## Building and Usage
This is a library, so it is typically imported by other projects.
To run tests (if available):
```bash
go test ./...
```
To build:
```bash
go build ./...
```
