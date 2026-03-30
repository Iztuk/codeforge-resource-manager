package views

import (
	"log"

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

	resourceLevel                   ResourceViewLevel
	selectedResource                string
	selectedResourceTable           string
	selectedResourceTableCell       int
	selectedResourceTableCellLength int

	width        int
	height       int
	menuWidth    int
	contentWidth int
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

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		if m.activeScreen == ScreenMenubar {
			switch msg.String() {
			case "j":
				if m.menuIndex < len(m.CurrentMenuItems())-1 {
					m.menuIndex++
					m.ResetContentSettings()
				}
				return m, nil
			case "k":
				if m.menuIndex > 0 {
					m.menuIndex--
					m.ResetContentSettings()
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
				if m.resourceLevel == ResourceLevelTables {
					m.resourceLevel = ResourceLevelList
					m.currentPage = ResourcesPage

					m.menuIndex = 0
					return m, nil
				}
				m.currentPage = page
				m.menuIndex = 0
				return m, nil
			}
		}

		const cols = 7
		const minSelectableCell = cols
		if m.activeScreen == ScreenContent && m.resourceLevel == ResourceLevelTables {
			switch msg.String() {

			case "h":
				col := m.selectedResourceTableCell % cols
				if m.selectedResourceTableCell > minSelectableCell && col > 0 {
					m.selectedResourceTableCell--
				}
				return m, nil

			case "l":
				col := m.selectedResourceTableCell % cols
				if col < cols-1 && m.selectedResourceTableCell+1 < m.selectedResourceTableCellLength {
					m.selectedResourceTableCell++
				}
				return m, nil

			case "j":
				next := m.selectedResourceTableCell + cols
				if next < m.selectedResourceTableCellLength {
					m.selectedResourceTableCell = next
				}
				return m, nil

			case "k":
				prev := m.selectedResourceTableCell - cols
				if prev >= minSelectableCell {
					m.selectedResourceTableCell = prev
				}
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		v := tea.NewView("loading...")
		v.AltScreen = true
		return v
	}

	m.menuWidth = 24
	m.contentWidth = max(20, m.width-m.menuWidth)

	title := m.titleView(m.width, 0)
	menu := m.menuBarView(m.menuWidth, m.height-5)
	content := m.contentView(m.contentWidth, m.height-5)

	body := lipgloss.JoinHorizontal(lipgloss.Top, menu, content)

	ui := lipgloss.JoinVertical(lipgloss.Top, title, body)

	v := tea.NewView(ui)
	v.AltScreen = true
	return v
}

func (m *model) menuBarView(width, height int) string {
	style := boxStyle
	if m.activeScreen == ScreenMenubar {
		style = focusedBoxStyle
	}

	return style.Width(width).Height(height).Render(m.GenerateMenuItems(m.CurrentMenuItems(), width))
}

func (m *model) contentView(width, height int) string {
	style := boxStyle
	if m.activeScreen == ScreenContent {
		style = focusedBoxStyle
	}

	return style.Width(width).Height(height).Render(m.GenerateContent())
}

func (m *model) titleView(width, height int) string {
	style := boxStyle.
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	return style.Render("CodeForge Resource Manager")
}

func StartView() {
	newStack := Stack[Page]{}
	newStack.Push(HomePage)
	m := &model{
		selectedResourceTableCell: 7,
		pageHistory:               newStack,
	}
	p := tea.NewProgram(m)

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
