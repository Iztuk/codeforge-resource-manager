package views

import (
	"resource-manager/internal/contracts"
	"resource-manager/internal/state"
	"sort"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

type ContentMode int

const (
	ContentPreview ContentMode = iota
	ContentAddResource
	ContentBindResource
)

func (m *model) ResetContentSettings() {
	m.selectedResourceTableCell = 11
	m.selectedBindResourceCell = 0
}

func (m *model) GenerateContent() string {

	switch m.currentPage {
	case HomePage:
		return "Select an option from the menu."

	case ResourcesPage:
		switch m.resourceLevel {
		case ResourceLevelList:
			if m.contentMode == ContentAddResource {
				return m.renderAddResourceForm()
			}

			if m.contentMode == ContentPreview && m.deleteResourceErrors != nil {
				errStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FF0000"))
				return errStyle.Render(m.deleteResourceErrors.Error())
			}

		case ResourceLevelTables:
			tableName := m.CurrentMenuSelection()
			if tableName == "" {
				return "No table selected."
			}

			return m.GenerateTableContentGrid(m.selectedResource, tableName)
		}

	case BindResourcePage:
		switch m.bindLevel {
		case PathList:
			return m.GeneratePathListContent()
		case PathItem:
			if m.contentMode == ContentBindResource {
				return m.GenerateBindResourceOptions()
			}
			return m.GeneratePathItemContent()
		}
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
	var colWidth int = max(8, m.contentWidth/cols)
	var maxContentLength int = max(4, colWidth-3) // Make space for border (2 chars)
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
			cells = append(cells, renderCell(stringItems[j], j == m.selectedResourceTableCell, colWidth, maxContentLength, j))
		}

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderCell(content string, selected bool, width, maxContentLength, cellIndex int) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Height(3).
		Width(max(1, width-1)).
		Align(lipgloss.Center, lipgloss.Center)

	if selected {
		style = style.BorderForeground(lipgloss.Color("#0087ff"))
	}

	isPermissionCol := cellIndex%7 >= 4
	switch content {
	case "":
		content = "nil"
	case "true":
		if isPermissionCol {
			style = style.Foreground(lipgloss.Color("#008000"))
		}
	case "false":
		if isPermissionCol {
			style = style.Foreground(lipgloss.Color("#FF8080"))
		}
	}

	if maxContentLength > 3 && len(content) > maxContentLength {
		return style.Render(content[:maxContentLength-3] + "...")
	}

	return style.Render(content)
}

func (m *model) initAddResourceForm() {
	m.nameInput = textinput.New()
	m.nameInput.Placeholder = "Resource name"
	m.nameInput.Prompt = ""
	m.nameInput.SetWidth(30)

	m.addrInput = textinput.New()
	m.addrInput.Placeholder = "sqlite://path/to/db"
	m.addrInput.Prompt = ""
	m.addrInput.SetWidth(40)

	m.focusedInput = 0
	m.nameInput.Focus()
	m.addrInput.Blur()
}

func (m *model) setFocusedInput(index int) {
	m.focusedInput = index

	m.nameInput.Blur()
	m.addrInput.Blur()

	switch index {
	case 0:
		m.nameInput.Focus()
	case 1:
		m.addrInput.Focus()
	}
}

func (m *model) moveFormFocus(delta int) {
	next := m.focusedInput + delta
	if next < 0 {
		next = 0
	}
	if next > 3 {
		next = 3
	}
	m.setFocusedInput(next)
}

func (m *model) renderAddResourceForm() string {
	var b strings.Builder

	b.WriteString("Add Resource\n\n")
	b.WriteString("Name\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")

	b.WriteString("Addr\n")
	b.WriteString(m.addrInput.View())
	b.WriteString("\n\n")

	errStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))
	for i, err := range m.addResourceFormErrors {
		if err == nil {
			continue
		}
		b.WriteString("\u2022 ")
		b.WriteString(errStyle.Render(err.Error()))
		if i < len(m.addResourceFormErrors)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m model) ToggleResourceTableCell(resourceName, tableName string) {
	resource, ok := state.AppState.ResourceContract.Resources[resourceName]
	if !ok || resource.DB == nil {
		return
	}

	table, ok := resource.DB.Tables[tableName]
	if !ok {
		return
	}

	var fieldNames []string
	for name := range table.Fields {
		fieldNames = append(fieldNames, name)
	}
	sort.Strings(fieldNames)

	row := m.selectedResourceTableCell / 7
	if row <= 0 || row-1 >= len(fieldNames) {
		return
	}

	selectedFieldName := fieldNames[row-1]
	selectedField := table.Fields[selectedFieldName]

	col := m.selectedResourceTableCell % 7
	switch col {
	case 4:
		selectedField.Read = !selectedField.Read
	case 5:
		selectedField.Write = !selectedField.Write
	case 6:
		selectedField.Mutable = !selectedField.Mutable
	default:
		return
	}

	table.Fields[selectedFieldName] = selectedField
	resource.DB.Tables[tableName] = table
	state.AppState.ResourceContract.Resources[resourceName] = resource

	state.WriteToResourceFile()
}
