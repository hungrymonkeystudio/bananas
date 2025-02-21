package main

import (
    lg "github.com/charmbracelet/lipgloss"
)

// color styles
var (
    white = lg.NewStyle().Foreground(lg.Color("#ffffff"))
    red = lg.NewStyle().Foreground(lg.Color("#cf513e"))
    yellow = lg.NewStyle().Foreground(lg.Color("#d6c43c"))
    gray = lg.NewStyle().Foreground(lg.Color("#787878"))
    cursor = lg.NewStyle().Foreground(lg.Color("#000000")).Background(lg.Color("#ffffff"))
)
