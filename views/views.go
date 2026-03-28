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

type Page int

const (
	HomePage Page = iota
	ResourcesPage
	BindResourcePage
)

type model struct {
	activeScreen Screen
	menuIndex    int

	currentPage Page
	pageHistory Stack[Page]

	resourceLevel    ResourceViewLevel
	selectedResource string
	selectedTable    string

	width  int
	height int
}

var (
	boxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("#808080")).
			Padding(1, 1)

	focusedBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 1).
			Bold(true)
)

func (m model) Init() tea.Cmd {
	m.pageHistory.Push(HomePage)
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
			if m.activeScreen > ScreenMenubar {
				m.activeScreen--
			}
			return m, nil
		case "ctrl+l":
			if m.activeScreen < ScreenContent {
				m.activeScreen++
			}
			return m, nil
		}

		if m.activeScreen == 0 {
			switch msg.String() {
			case "j":
				if m.menuIndex < len(m.CurrentMenuItems())-1 {
					m.menuIndex++
				}
				return m, nil
			case "k":
				if m.menuIndex > 0 {
					m.menuIndex--
				}
				return m, nil
			case "enter":
				m.SelectMenuItem()
				m.menuIndex = 0
				return m, nil
			case "backspace":
				page, t := m.pageHistory.Pop()
				if !t {
					m.currentPage = HomePage
					return m, nil
				}
				m.currentPage = page
				m.menuIndex = 0
				return m, nil
			}
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
	if m.activeScreen == ScreenMenubar {
		style = focusedBoxStyle
	}

	return style.Width(width).Height(height).Render(m.GenerateMenuItems(m.CurrentMenuItems(), width))
}

func (m model) contentView(width, height int) string {
	style := boxStyle
	if m.activeScreen == ScreenContent {
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
