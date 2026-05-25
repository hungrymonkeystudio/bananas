# 🛠 Development Guide

This document describes the high-level architecture and the state machine behind `bananas`.

## 🏗 High-Level Architecture

`bananas` is built using the **Model-View-Update (MVU)** architecture, also known as The Elm Architecture, through the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

The main application is split into three core packages:

- **`pkg/typer/`:** The core typing logic, word bank loading, and character color management.
- **`pkg/settings/`:** Manages user configuration, timer durations, and persistence.
- **`pkg/analysis/`:** Calculates WPM, accuracy, and other post-test statistics.

### 🔄 State Machine

The application transitions through several key states during a typing test:

1. **Initial State:** Waiting for the user's first keypress.
2. **Typing State:** The test is active, characters are being tracked, and the timer is running.
3. **Finish State:** The test has ended (timer hit zero or word count reached), and analysis begins.
4. **Settings State:** The user is configuring the test parameters (ESC key).

## 🧩 Components

### `TyperModel` (pkg/typer)
- **Lines:** A 2D slice of strings representing the lines of words on the screen.
- **LinesColor:** A parallel 2D slice of colors (encoded as 'g', 'r', 'w') for each character.
- **Skips:** Tracks locations of skipped words for proper backspace behavior.

### `SettingsModel` (pkg/settings)
- **Show:** Boolean flag for rendering the settings screen overlay.
- **ActiveTime / ActiveWords:** Persistent values loaded from `settings.json`.

### `Analysis` (pkg/analysis)
- **WPM Calculation:** Uses the formula `(TotalTyped / 5) / (TimeInSeconds / 60)`.
- **Accuracy:** `(TotalCorrect / TotalTyped) * 100`.

## 🛠 Contributing

1. **Fork the repository.**
2. **Create a feature branch:** `git checkout -b feat/my-new-feature`
3. **Make your changes.**
4. **Test your changes:** Run `go run main.go` to verify the TUI behavior.
5. **Submit a Pull Request.**

### 🔍 Guidelines
- **Go Conventions:** Adhere to standard Go formatting (`gofmt`).
- **Minimal Dependencies:** Avoid adding external libraries unless necessary. `bananas` aims to be lightweight.
- **Clean UI:** Ensure any UI changes maintain the minimal aesthetic.

## Component Diagram

```
+-----------------------------------------------------------+
|                        main.go                            |
|              (entry point, build config)                   |
+----------------------------+------------------------------+
                             |
                             v
+-----------------------------------------------------------+
|              internal/session (Model)                      |
|     state machine: typing <-> settings -> results         |
|     owns: timer, progress tracking, transitions           |
+----------+-----------------+------------------------------+
           |                 |                 |
           v                 v                 v
+----------------+  +----------------+  +----------------+
| ui/typing      |  | ui/settings    |  | ui/results     |
|                |  |                |  |                |
| word gen       |  | UI + controls  |  | WPM + accuracy |
| input handling |  | JSON persist   |  | retry flow     |
| rendering      |  |                |  |                |
+----------------+  +----------------+  +----------------+
           |                 |
           v                 v
+-----------------------------------------------------------+
|                   foundation/                              |
|  theme (styles)  |  logger (file log)  |  paths (resolve) |
+-----------------------------------------------------------+
```

Dependency rule: each layer only imports from layers below it.

