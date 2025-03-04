package styles

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	orange    = lipgloss.Color("#f66c00")
	blue      = lipgloss.Color("#20b9f7")
	red       = lipgloss.Color("#ec5353")
	green     = lipgloss.Color("#2ecc71")
	dark      = lipgloss.Color("#3c3f42")
	paragraph = lipgloss.Color("#8c969e")
	lightText = lipgloss.Color("#e0e2e9")
	light     = lipgloss.Color("#f2f3f6")
	purple    = lipgloss.Color("#6A5ACD")
)

func PlaceContent(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

type mainMenu struct {
	Page,
	Title,
	MenuTitle,
	Description,
	Selections,
	Selected,
	Updates,
	Summary,
	Error,
	Cli,
	Footer lipgloss.Style
}

type summary struct {
	Base,
	Container,
	Action,
	Status,
	KeyWord,
	Items,
	Footer lipgloss.Style
}

func SummaryStyles(status bool) summary {
	var s summary

	s.Base = lipgloss.NewStyle().
		Width(80).
		MarginTop(1).
		Align(lipgloss.Center)

	s.Container = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(orange).
		Height(30).
		Width(80)

	switch status {
	case true:
		s.Status = lipgloss.NewStyle().
			Inherit(s.Base).
			MarginTop(1).
			Foreground(green)
	default:
		s.Status = lipgloss.NewStyle().
			Inherit(s.Base).
			MarginTop(1).
			Foreground(red)
	}

	s.Items = lipgloss.NewStyle().
		Foreground(orange)

	s.Action = lipgloss.NewStyle().
		Inherit(s.Base).
		MarginTop(1).
		Foreground(red)

	s.Footer = lipgloss.NewStyle().
		Inherit(s.Base).
		MarginTop(1).
		Foreground(red)

	s.KeyWord = lipgloss.NewStyle().
		Foreground(blue)

	return s
}

func DefaultStyles() mainMenu {
	var s mainMenu

	s.Cli = lipgloss.NewStyle().
		Foreground(orange)

	s.Page = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(orange).
		Width(80).
		Height(30).
		Align(lipgloss.Center)

	s.Title = lipgloss.NewStyle().
		Foreground(red).
		MarginBottom(1)

	s.MenuTitle = lipgloss.NewStyle().
		Foreground(light).
		Bold(true).
		Align(lipgloss.Center).
		MarginTop(1).
		MarginBottom(1)

	s.Description = lipgloss.NewStyle().
		Foreground(orange).
		Width(80).
		MarginTop(1).
		MarginBottom(2).
		Align(lipgloss.Center)

	s.Selections = lipgloss.NewStyle().
		MarginTop(1).
		Width(25).
		AlignHorizontal(lipgloss.Center).
		Foreground(lightText)

	s.Selected = lipgloss.NewStyle().
		Inherit(s.Selections).Background(dark)

	s.Updates = lipgloss.NewStyle().
		Inherit(s.Description).
		Foreground(orange).
		Height(30).
		AlignVertical(lipgloss.Center)

	s.Footer = lipgloss.NewStyle().
		Foreground(red).
		MarginTop(4)

	s.Summary = lipgloss.NewStyle().
		Inherit(s.Page).
		Foreground(orange)

	s.Error = lipgloss.NewStyle().
		Foreground(red)

	return s
}

func FormBaseStyles() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Title = t.Focused.Title.Foreground(light).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(orange).Bold(true).MarginBottom(1).MarginTop(1).Width(80)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(red).Bold(true).AlignHorizontal(lipgloss.Center).Width(80)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(orange).Blink(true)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(lightText).Background(dark).Bold(true)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(dark)
	t.Focused.Base = t.Focused.Base.BorderLeft(false)
	t.Blurred = t.Focused

	t.Blurred.NextIndicator = t.Focused.NextIndicator.Foreground(red)
	t.Blurred.PrevIndicator = t.Focused.PrevIndicator.Foreground(red)

	t.Help = help.New().Styles

	t.Help.ShortKey = lipgloss.NewStyle().Foreground(red)
	t.Help.ShortDesc = lipgloss.NewStyle().Foreground(red)
	t.Help.ShortSeparator = lipgloss.NewStyle().Foreground(red)

	return t
}

func FormStyles() *huh.Theme {

	t := FormBaseStyles()
	t.Focused.Description = t.Focused.Description.Foreground(dark).Bold(true).MarginBottom(1).MarginTop(1).Width(80)

	// Input styles
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(orange)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(orange)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(blue)

	// Multi select styles
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(orange).Blink(true)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().SetString("[🔥] ").Foreground(orange)
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().SetString("[ ] ").Foreground(purple)
	t.Blurred = t.Focused

	return t
}

func AgentsPageStyle(agent string) *huh.Theme {
	t := FormStyles()

	switch agent {
	case "Telegraf":
		t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(purple).Bold(true).AlignHorizontal(lipgloss.Center).Width(80)
	default:
		t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(red).Bold(true).AlignHorizontal(lipgloss.Center).Width(80)
	}
	t.Blurred = t.Focused

	return t
}

func CustomKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.MultiSelect.Filter.Unbind()
	km.Select.Filter.Unbind()
	km.MultiSelect.Toggle.SetHelp("x", "select")

	return km
}

// Custom styling for Cobra CLI
func CustomUsageFunc(cmd *cobra.Command) error {
	var headerStyle = lipgloss.NewStyle().Bold(true).Foreground(orange)

	usage := headerStyle.Render("Usage:") + "\n  " + cmd.UseLine()

	if cmd.HasAvailableSubCommands() {
		usage += "\n  " + cmd.CommandPath() + " [command]"
	}
	if len(cmd.Aliases) > 0 {
		usage += "\n\n" + headerStyle.Render("Aliases:") + "\n  " + cmd.NameAndAliases()
	}
	if cmd.Example != "" {
		usage += "\n\n" + headerStyle.Render("Examples:") + "\n  " + cmd.Example
	}
	if cmd.HasAvailableSubCommands() {
		usage += "\n\n" + headerStyle.Render("Available Commands:") + "\n"
		for _, c := range cmd.Commands() {
			if c.IsAvailableCommand() || c.Name() == "help" {
				usage += "  " + c.Name() + "\t" + c.Short + "\n"
			}
		}
	}
	if cmd.HasAvailableLocalFlags() {
		usage += "\n\n" + headerStyle.Render("Flags:") + "\n" + cmd.LocalFlags().FlagUsages()
	}
	if cmd.HasAvailableInheritedFlags() {
		usage += "\n\n" + headerStyle.Render("Global Flags:") + "\n" + cmd.InheritedFlags().FlagUsages()
	}
	if cmd.HasHelpSubCommands() {
		usage += "\n\n" + headerStyle.Render("Additional help topics:") + "\n"
		for _, c := range cmd.Commands() {
			if c.IsAdditionalHelpTopicCommand() {
				usage += "  " + c.CommandPath() + "\t" + c.Short + "\n"
			}
		}
	}
	if cmd.HasAvailableSubCommands() {
		usage += "\n\nUse \"" + cmd.CommandPath() + " [command] --help\" for more information about a command."
	}
	fmt.Println(usage)
	return nil
}
