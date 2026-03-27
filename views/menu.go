package views

import (
	"resource-manager/internal/state"
	"strings"

	"charm.land/lipgloss/v2"
)

type MenuItem int

const (
	HomeResources MenuItem = iota
	HomeBindResource

	numMenuItem
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
		if i == int(m.menuItem) {
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
		var resources []string
		for _, res := range state.AppState.ResourceContract.Resources {
			resources = append(resources, res.Name)
		}
		return resources
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
		switch m.menuItem {
		case 0:
			m.currentPage = ResourcesPage
			m.menuItem = 0
			return
		case 1:
			m.currentPage = BindResourcePage
			m.menuItem = 0
			return
		}
	case ResourcesPage:
	case BindResourcePage:
	}
}
