package views

import (
	"log"
	"resource-manager/internal/resources"
	"resource-manager/internal/state"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
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

	bindLevel        BindViewLevel
	selectedPath     string
	selectedPathItem string

	contentMode  ContentMode
	focusedInput int

	// Add Resource Form
	nameInput             textinput.Model
	addrInput             textinput.Model
	addResourceFormErrors []error

	// Delete Resource
	deleteResourceErrors error

	width           int
	height          int
	menuWidth       int
	contentWidth    int
	menuViewport    viewport.Model
	contentViewport viewport.Model
}

var (
	boxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 1)

	focusedBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#0087ff")).
			Padding(1, 1)
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.menuWidth = m.width / 5
		m.contentWidth = max(20, m.width-m.menuWidth)

		contentHeight := m.height - 6
		if contentHeight < 1 {
			contentHeight = 1
		}

		m.menuViewport = viewport.New(
			viewport.WithWidth(max(1, m.menuWidth-4)),
			viewport.WithHeight(max(1, contentHeight-4)),
		)
		m.menuViewport.SetContent(m.GenerateMenuItems(m.CurrentMenuItems(), m.menuWidth))

		m.contentViewport = viewport.New(
			viewport.WithWidth(max(1, m.contentWidth-4)),
			viewport.WithHeight(max(1, contentHeight-4)),
		)
		m.contentViewport.SetContent(m.GenerateContent())

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
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

		if m.currentPage == ResourcesPage && m.resourceLevel == ResourceLevelList {
			switch msg.String() {
			case "ctrl+a":
				m.activeScreen = ScreenContent
				m.contentMode = ContentAddResource
				m.initAddResourceForm()
				return m, textinput.Blink

			case "ctrl+d":
				items := m.CurrentMenuItems()
				var menuItemName string
				for i, item := range items {
					if i == m.menuIndex {
						menuItemName = item
						break
					}
				}
				m.deleteResourceErrors = state.DeleteResource(menuItemName)
				return m, nil
			}
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
				if m.resourceLevel != ResourceLevelTables && m.bindLevel != PathItem {
					m.SelectMenuItem()
					m.menuIndex = 0
				}
				return m, nil
			case "backspace":
				if m.resourceLevel == ResourceLevelTables {
					m.resourceLevel = ResourceLevelList
					m.currentPage = ResourcesPage

					m.menuIndex = 0
					return m, nil
				}
				if m.bindLevel == PathItem {
					m.bindLevel = PathList
					m.currentPage = BindResourcePage

					m.menuIndex = 0
					return m, nil
				}

				page, t := m.pageHistory.Pop()
				if !t {
					m.currentPage = HomePage
					return m, nil
				}
				m.currentPage = page
				m.menuIndex = 0
				return m, nil
			case "ctrl+n":
				m.menuViewport.ScrollDown(1)
				return m, nil
			case "ctrl+p":
				m.menuViewport.ScrollUp(1)
				return m, nil
			}
		}

		if m.activeScreen == ScreenContent {
			switch msg.String() {
			case "ctrl+n":
				m.contentViewport.ScrollDown(1)
				return m, nil
			case "ctrl+p":
				m.contentViewport.ScrollUp(1)
				return m, nil
			}
		}

		if m.activeScreen == ScreenContent && m.contentMode == ContentAddResource {
			switch msg.String() {
			case "esc":
				m.contentMode = ContentPreview
				m.activeScreen = ScreenMenubar
				return m, nil
			case "tab", "down":
				m.moveFormFocus(1)
				return m, nil
			case "shift+tab", "up":
				m.moveFormFocus(-1)
				return m, nil
			case "enter":
				m.addResourceFormErrors = make([]error, 0)
				m.addResourceFormErrors = append(m.addResourceFormErrors, resources.AddDb(m.nameInput.Value(), m.addrInput.Value())...)

				if len(m.addResourceFormErrors) == 0 {
					m.nameInput.Reset()
					m.addrInput.Reset()
					m.activeScreen = ScreenMenubar
					m.contentMode = ContentPreview
				}

				return m, nil
			}

			var cmd1, cmd2 tea.Cmd
			m.nameInput, cmd1 = m.nameInput.Update(msg)
			m.addrInput, cmd2 = m.addrInput.Update(msg)

			return m, tea.Batch(cmd1, cmd2)
		}

		const cols = 7
		const minSelectableCell = cols
		const selectableStart = cols - 3
		const selectableEnd = cols - 1
		if m.activeScreen == ScreenContent && m.resourceLevel == ResourceLevelTables {
			switch msg.String() {

			case "h":
				col := m.selectedResourceTableCell % cols
				if m.selectedResourceTableCell > minSelectableCell && col > selectableStart {
					m.selectedResourceTableCell--
				}
				return m, nil

			case "l":
				col := m.selectedResourceTableCell % cols
				if col < selectableEnd && m.selectedResourceTableCell+1 < m.selectedResourceTableCellLength {
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
			case "enter", "space":
				m.ToggleResourceTableCell(m.selectedResource, m.CurrentMenuSelection())
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

	title := m.titleView(m.width, 0)
	menu := m.menuBarView(m.menuWidth, m.height-6)
	content := m.contentView(m.contentWidth, m.height-6)
	cmdHelp := m.cmdView(m.width, 0)

	body := lipgloss.JoinHorizontal(lipgloss.Top, menu, content)

	ui := lipgloss.JoinVertical(lipgloss.Top, title, body, cmdHelp)

	v := tea.NewView(ui)
	v.AltScreen = true
	return v
}

func (m *model) menuBarView(width, height int) string {
	style := boxStyle
	if m.activeScreen == ScreenMenubar {
		style = focusedBoxStyle
	}

	m.menuViewport.SetContent(m.GenerateMenuItems(m.CurrentMenuItems(), width))
	return style.Width(width).Height(height).Render(m.menuViewport.View())
}

func (m *model) contentView(width, height int) string {
	style := boxStyle
	if m.activeScreen == ScreenContent {
		style = focusedBoxStyle
	}

	m.contentViewport.SetContent(m.GenerateContent())
	return style.Width(width).Height(height).Render(m.contentViewport.View())
}

func (m *model) titleView(width, height int) string {
	style := boxStyle.
		Padding(0).
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true)

	switch m.currentPage {
	case HomePage:
		return style.Render("CodeForge Resource Manager")
	case ResourcesPage:
		switch m.resourceLevel {
		case ResourceLevelList:
			return style.Render("Resources")
		case ResourceLevelTables:
			return style.Render(m.selectedResource)
		}
	case BindResourcePage:
		return style.Render("Bind Resource")
	}

	return style.Render("CodeForge Resource Manager")
}

func (m *model) cmdView(width, height int) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0).
		Width(width).
		Height(height).
		Foreground(lipgloss.Color("#ffffff"))

	var b strings.Builder

	switch m.currentPage {
	case HomePage:
	case ResourcesPage:
		switch m.resourceLevel {
		case ResourceLevelList:
			if m.contentMode == ContentAddResource {
				b.WriteString(" ")
				b.WriteString("Navigation: tab (down)/shift+tab (up)")
				b.WriteString(" | ")
				b.WriteString("Cancel: esc")
				b.WriteString(" | ")
				b.WriteString("Save: enter")

			} else {
				b.WriteString(" ")
				b.WriteString("Add Resource: ctrl + a")
				b.WriteString(" | ")
				b.WriteString("Delete Resource: ctrl + d")

			}
		case ResourceLevelTables:
		}
	case BindResourcePage:
	}

	return style.Render(b.String())
}

func StartView() {
	newStack := Stack[Page]{}
	newStack.Push(HomePage)
	m := &model{
		selectedResourceTableCell: 11,
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
