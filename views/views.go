package views

import (
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Screen int

const (
	ScreenMenubar Screen = iota
	ScreenContent
)

type model struct {
	active Screen
	width  int
	height int
}

var (
	boxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#0087ff")).
			Foreground(lipgloss.Color("#808080")).
			Padding(1, 1)

	focusedBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#0087ff")).
			Padding(1, 1).
			Bold(true)
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+h":
			if m.active > ScreenMenubar {
				m.active--
			}
			return m, nil
		case "ctrl+l":
			if m.active < ScreenContent {
				m.active++
			}
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		v := tea.NewView("loading...")
		v.AltScreen = true
		return v
	}

	menuWidth := 24
	contentWidth := max(20, m.width-menuWidth)

	title := m.titleView(m.width, 0)
	menu := m.menuBarView(menuWidth, m.height-5)
	content := m.contentView(contentWidth, m.height-5)

	body := lipgloss.JoinHorizontal(lipgloss.Top, menu, content)

	ui := lipgloss.JoinVertical(lipgloss.Top, title, body)

	v := tea.NewView(ui)
	v.AltScreen = true
	return v
}

func (m model) menuBarView(width, height int) string {
	style := boxStyle
	if m.active == ScreenMenubar {
		style = focusedBoxStyle
	}

	var b strings.Builder
	b.WriteString("Menu\n\n")
	b.WriteString("Resources\n")
	b.WriteString("API Paths\n")
	b.WriteString("Quit\n\n")
	b.WriteString("ctrl+h / ctrl+l to switch focus")

	return style.Width(width).Height(height).Render(b.String())
}

func (m model) contentView(width, height int) string {
	style := boxStyle
	if m.active == ScreenContent {
		style = focusedBoxStyle
	}

	var b strings.Builder
	b.WriteString("Content\n\n")
	b.WriteString("This is the content pane.\n")
	b.WriteString("Later you can show details here.")

	return style.Width(width).Height(height).Render(b.String())
}

func (m model) titleView(width, height int) string {
	style := boxStyle.
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	return style.Render("CodeForge Resource Manager")
}

func StartView() {
	p := tea.NewProgram(model{})

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
