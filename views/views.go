package views

import (
	"log"
	"resource-manager/internal/resources"
	"resource-manager/internal/state"

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

	bindLevel                BindViewLevel
	selectedPath             string
	selectedPathItem         string
	selectedBindResourceCell int
	bindResourceCellLength   int
	bindResourceRows         []BindResourceRow

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

	debug string
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

		contentHeight := m.height - 3
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
			if m.contentMode == ContentBindResource {
				m.contentMode = ContentPreview
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
				if m.menuIndex != 0 {
					m.menuIndex--
				}
				return m, nil
			}
		}

		if m.currentPage == BindResourcePage && m.bindLevel == PathItem {
			switch msg.String() {
			case "ctrl+a":
				m.activeScreen = ScreenContent
				m.contentMode = ContentBindResource
				return m, nil
			case "ctrl+d":
				m.RemoveResourceBinding()
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
				if m.CurrentMenuSelection() == "Help" {
					return m, nil
				}

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

		if m.activeScreen == ScreenContent && m.contentMode == ContentBindResource {
			switch msg.String() {
			case "h":
				if m.selectedBindResourceCell > 0 {
					m.selectedBindResourceCell--
				}
				return m, nil

			case "l":
				if m.selectedBindResourceCell+1 < m.bindResourceCellLength {
					m.selectedBindResourceCell++
				}
				return m, nil

			case "j":
				row, col, ok := m.findBindResourcePosition()
				if !ok {
					return m, nil
				}

				nextRow := row + 1
				if nextRow >= len(m.bindResourceRows) {
					return m, nil
				}

				if col >= len(m.bindResourceRows[nextRow].Indices) {
					col = len(m.bindResourceRows[nextRow].Indices) - 1
				}

				m.selectedBindResourceCell = m.bindResourceRows[nextRow].Indices[col]
				return m, nil

			case "k":
				row, col, ok := m.findBindResourcePosition()
				if !ok {
					return m, nil
				}

				prevRow := row - 1
				if prevRow < 0 {
					return m, nil
				}

				if col >= len(m.bindResourceRows[prevRow].Indices) {
					col = len(m.bindResourceRows[prevRow].Indices) - 1
				}

				m.selectedBindResourceCell = m.bindResourceRows[prevRow].Indices[col]
				return m, nil

			case "enter":
				selected := m.currentSelectedBindResource()
				if selected != "" {
					m.BindResourceToEndpoint(selected)
					m.contentMode = ContentPreview
					m.activeScreen = ScreenMenubar
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

	title := m.titleView(m.width, 0)
	menu := m.menuBarView(m.menuWidth, m.height-3)
	content := m.contentView(m.contentWidth, m.height-3)

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
		switch m.bindLevel {
		case PathList:
			return style.Render("Bind Resource")
		case PathItem:
			return style.Render(m.debug)
			// return style.Render(m.selectedPath)
		}
	}

	return style.Render("CodeForge Resource Manager")
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
