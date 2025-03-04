package types

import (
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (View, tea.Cmd)
	View() string
}
