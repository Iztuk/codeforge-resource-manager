package views

import (
	"resource-manager/internal/contracts"
	"resource-manager/internal/state"
	"sort"
	"strconv"

	"charm.land/lipgloss/v2"
)

func (m *model) ResetContentSettings() {
	m.selectedResourceTableCell = 7
}

func (m *model) GenerateContent() string {
	switch m.currentPage {
	case HomePage:
		return "Select an option from the menu."

	case ResourcesPage:
		switch m.resourceLevel {
		case ResourceLevelList:

		case ResourceLevelTables:
			tableName := m.CurrentMenuSelection()
			if tableName == "" {
				return "No table selected."
			}

			return m.GenerateTableContentGrid(m.selectedResource, tableName)
		}

	case BindResourcePage:
		return "Bind Resource page"
	}

	return ""
}

func (m *model) GeneratePreviewTableFromResourceList() string {
	resourceName := m.CurrentMenuSelection()
	if resourceName == "" {
		return "No resource selected."
	}

	resource, ok := state.AppState.ResourceContract.Resources[resourceName]
	if !ok || resource.DB == nil {
		return "Selected resource has no database tables."
	}

	if len(resource.DB.Tables) == 0 {
		return "Selected resource has no tables."
	}

	var tableNames []string
	for name := range resource.DB.Tables {
		tableNames = append(tableNames, name)
	}
	sort.Strings(tableNames)

	// Use the first table as the preview table
	tableName := tableNames[0]
	table := resource.DB.Tables[tableName]

	var fieldNames []string
	for name := range table.Fields {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	var fields []contracts.FieldSpec
	for _, name := range fieldNames {
		fields = append(fields, table.Fields[name])
	}

	return m.renderTableGrid(fields)
}

func (m *model) GenerateTableContentGrid(resourceName, tableName string) string {
	resource, ok := state.AppState.ResourceContract.Resources[resourceName]
	if !ok || resource.DB == nil {
		return ""
	}

	table, ok := resource.DB.Tables[tableName]
	if !ok {
		return ""
	}

	var fieldNames []string
	for name := range table.Fields {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	var fields []contracts.FieldSpec
	for _, name := range fieldNames {
		fields = append(fields, table.Fields[name])
	}

	return m.renderTableGrid(fields)
}

func (m *model) renderTableGrid(items []contracts.FieldSpec) string {
	var cols int = 7 // Number of columns to display FieldSpec
	var colWidth int = m.contentWidth / cols
	var maxContentLength int = colWidth - 3 // Make space for border (2 chars)
	var rows []string

	var stringItems []string = []string{
		"Column", "Type", "Nullable", "Default", "Read", "Write", "Mutable",
	}

	for _, item := range items {
		var defaultVal string = ""
		if item.Default != nil {
			defaultVal = *item.Default
		}

		stringItem := []string{
			item.ColumnName, item.Type, strconv.FormatBool(item.Nullable), defaultVal, strconv.FormatBool(item.Read), strconv.FormatBool(item.Write), strconv.FormatBool(item.Mutable),
		}

		stringItems = append(stringItems, stringItem...)
	}
	m.selectedResourceTableCellLength = len(stringItems)

	for i := 0; i < len(stringItems); i += cols {
		end := min(i+cols, len(stringItems))
		var cells []string
		for j := i; j < end; j++ {
			cells = append(cells, renderCell(stringItems[j], j == m.selectedResourceTableCell, colWidth, maxContentLength))
		}

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderCell(content string, selected bool, width, maxContentLength int) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Height(3).
		Width(width-1).
		Align(lipgloss.Center, lipgloss.Center)

	if selected {
		style = style.BorderForeground(lipgloss.Color("#0087ff"))
	}
	if content == "" {
		content = "nil"
	}
	if len(content) > maxContentLength {
		return style.Render(content[:maxContentLength-3] + "...")
	}

	return style.Render(content)
}
