package colors 

import (
    lg "github.com/charmbracelet/lipgloss"
)

// color styles
var (
    White = lg.NewStyle().Foreground(lg.Color("#ffffff"))
    Red = lg.NewStyle().Foreground(lg.Color("#cf513e"))
    Yellow = lg.NewStyle().Foreground(lg.Color("#d6c43c"))
    Gray = lg.NewStyle().Foreground(lg.Color("#787878"))
    Cursor = lg.NewStyle().Foreground(lg.Color("#000000")).Background(lg.Color("#ffffff"))
    Instructions = lg.NewStyle().Foreground(lg.Color("#525252"))
    Underline = lg.NewStyle().Foreground(lg.Color("#ffffff")).Underline(true)
)
