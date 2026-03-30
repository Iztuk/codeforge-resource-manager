package views

import (
	"resource-manager/internal/state"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
)

type ResourceViewLevel int

const (
	ResourceLevelList ResourceViewLevel = iota
	ResourceLevelTables
)

func (m model) GenerateMenuItems(menuItems []string, width int) string {
	width = width - 4
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Width(width).
		Align(lipgloss.Center, lipgloss.Center)

	selectedStyle := style.
		BorderForeground(lipgloss.Color("#0087ff"))

	var b strings.Builder
	for i, item := range menuItems {
		if i == m.menuIndex {
			b.WriteString(selectedStyle.Render(item))
		} else {
			b.WriteString(style.Render(item))
		}

		if i < len(menuItems)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m model) CurrentMenuItems() []string {
	switch m.currentPage {
	case HomePage:
		return []string{
			"Resources",
			"Bind Resource",
		}
	case ResourcesPage:
		switch m.resourceLevel {
		case ResourceLevelList:
			var resources []string
			for _, item := range state.AppState.ResourceContract.Resources {
				resources = append(resources, item.Name)
			}

			sort.Strings(resources)
			return resources
		case ResourceLevelTables:
			var tables []string

			if state.AppState.ResourceContract.Resources[m.selectedResource].DB != nil {
				for key := range state.AppState.ResourceContract.Resources[m.selectedResource].DB.Tables {
					tables = append(tables, key)
				}
			}

			sort.Strings(tables)
			return tables
		}
	case BindResourcePage:
	default:
		return []string{
			"Resources",
			"Bind Resource",
		}
	}
	return []string{}
}

func (m *model) SelectMenuItem() {
	switch m.currentPage {
	case HomePage:
		switch m.menuIndex {
		case 0:
			m.pageHistory.Push(m.currentPage)
			m.currentPage = ResourcesPage
			return
		case 1:
			m.pageHistory.Push(m.currentPage)
			m.currentPage = BindResourcePage
			return
		}
	case ResourcesPage:
		switch m.resourceLevel {
		case ResourceLevelList:
			resources := m.CurrentMenuItems()

			m.selectedResource = resources[m.menuIndex]
			m.resourceLevel = ResourceLevelTables
			return
		case ResourceLevelTables:
			tables := m.CurrentMenuItems()

			m.selectedResourceTable = tables[m.menuIndex]
			return
		}
	case BindResourcePage:
	}
}

func (m model) CurrentMenuSelection() string {
	items := m.CurrentMenuItems()
	if len(items) == 0 {
		return ""
	}
	if m.menuIndex < 0 || m.menuIndex >= len(items) {
		return ""
	}
	return items[m.menuIndex]
}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}
