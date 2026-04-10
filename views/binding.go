package views

import (
	"fmt"
	"image/color"
	"resource-manager/internal/contracts"
	"resource-manager/internal/state"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *model) GeneratePathListContent() string {
	var paths []string
	for key := range state.AppState.ApiContract.Paths {
		paths = append(paths, key)
	}
	sort.Strings(paths)
	if len(paths) == 0 || m.menuIndex < 0 || m.menuIndex >= len(paths) {
		return ""
	}
	m.selectedPath = paths[m.menuIndex]

	return generatePathListContentStrings(m.selectedPath, m.contentWidth)
}

func generatePathListContentStrings(path string, width int) string {
	pathItem := state.AppState.ApiContract.Paths[path]

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0087ff")).
		Align(lipgloss.Center).
		Width(width)

	var methodCards []string

	if pathItem.GET != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("GET", pathItem.GET, width))
	}
	if pathItem.POST != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("POST", pathItem.POST, width))
	}
	if pathItem.PUT != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("PUT", pathItem.PUT, width))
	}
	if pathItem.PATCH != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("PATCH", pathItem.PATCH, width))
	}
	if pathItem.DELETE != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("DELETE", pathItem.DELETE, width))
	}
	if pathItem.HEAD != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("HEAD", pathItem.HEAD, width))
	}
	if pathItem.OPTIONS != nil {
		methodCards = append(methodCards, generateOpenApiOperationCardStrings("OPTIONS", pathItem.OPTIONS, width))
	}

	methodColumn := lipgloss.JoinVertical(lipgloss.Center, methodCards...)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(path),
		methodColumn,
	)
}

func generateOpenApiOperationCardStrings(method string, operation *contracts.OpenApiOperation, width int) string {
	cardWidth := max(24, width/2)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(cardWidth).
		Align(lipgloss.Left)

	methodStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(methodColor(method)).
		Width(cardWidth - 2).
		Align(lipgloss.Center)

	labelStyle := lipgloss.NewStyle().
		Bold(true)

	contentStyle := lipgloss.NewStyle().
		Width(cardWidth - 2).
		Align(lipgloss.Left)

	var rows []string

	rows = append(rows, methodStyle.Render(method))

	if operation.OperationID != "" {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Operation ID: ")+operation.OperationID,
		))
	}

	if operation.Summary != "" {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Summary: ")+operation.Summary,
		))
	}

	if operation.Description != "" {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Description: ")+operation.Description,
		))
	}

	if len(operation.Tags) > 0 {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Tags: ")+strings.Join(operation.Tags, ", "),
		))
	}

	if operation.XResource != nil {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render(fmt.Sprintf("Resource: %s.%s", operation.XResource.ResourceName, operation.XResource.Table)),
		))
	}

	cardColumn := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return cardStyle.Render(cardColumn)
}

func generateOpenApiOperationStrings(method string, operation *contracts.OpenApiOperation, width int) string {
	methodStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(methodColor(method)).
		Width(width).
		Align(lipgloss.Left).
		Padding(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true)

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Padding(0, 1, 0, 1)

	var rows []string

	rows = append(rows, methodStyle.Render(method))

	if operation.OperationID != "" {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Operation ID: ")+operation.OperationID,
		))
	}

	if operation.Summary != "" {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Summary: ")+operation.Summary,
		))
	}

	if operation.Description != "" {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Description: ")+operation.Description,
		))
	}

	if len(operation.Tags) > 0 {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render("Tags: ")+strings.Join(operation.Tags, ", "),
		))
	}

	if operation.XResource != nil {
		rows = append(rows, contentStyle.Render(
			labelStyle.Render(fmt.Sprintf("Resource: %s.%s", operation.XResource.ResourceName, operation.XResource.Table)),
		))
	}

	cardColumn := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return cardColumn
}

func methodColor(method string) color.Color {
	switch method {
	case "GET":
		return lipgloss.Color("#00D787") // green
	case "POST":
		return lipgloss.Color("#FFD700") // yellow
	case "PUT":
		return lipgloss.Color("#0087FF") // blue
	case "PATCH":
		return lipgloss.Color("#00BFFF") // cyan
	case "DELETE":
		return lipgloss.Color("#FF5F5F") // red
	case "HEAD":
		return lipgloss.Color("#888888") // gray
	case "OPTIONS":
		return lipgloss.Color("#AF5FFF") // purple
	default:
		return lipgloss.Color("#FFFFFF")
	}
}

func GeneratePathItemMethods(pathItem contracts.OpenApiPathItem) []string {
	var methods []string

	if pathItem.GET != nil {
		methods = append(methods, "GET")
	}
	if pathItem.POST != nil {
		methods = append(methods, "POST")
	}
	if pathItem.PUT != nil {
		methods = append(methods, "PUT")
	}
	if pathItem.PATCH != nil {
		methods = append(methods, "PATCH")
	}
	if pathItem.DELETE != nil {
		methods = append(methods, "DELETE")
	}
	if pathItem.HEAD != nil {
		methods = append(methods, "HEAD")
	}
	if pathItem.OPTIONS != nil {
		methods = append(methods, "OPTIONS")
	}

	return methods
}

func (m *model) GeneratePathItemContent() string {
	pathItem := state.AppState.ApiContract.Paths[m.selectedPath]
	pathItemMethods := GeneratePathItemMethods(state.AppState.ApiContract.Paths[m.selectedPath])
	if len(pathItemMethods) == 0 || m.menuIndex < 0 || m.menuIndex >= len(pathItemMethods) {
		return ""
	}

	switch pathItemMethods[m.menuIndex] {
	case "GET":
		return generateOpenApiOperationStrings("GET", pathItem.GET, m.contentWidth)
	case "POST":
		return generateOpenApiOperationStrings("POST", pathItem.POST, m.contentWidth)
	case "PUT":
		return generateOpenApiOperationStrings("PUT", pathItem.PUT, m.contentWidth)
	case "PATCH":
		return generateOpenApiOperationStrings("PATCH", pathItem.PATCH, m.contentWidth)
	case "DELETE":
		return generateOpenApiOperationStrings("DELETE", pathItem.DELETE, m.contentWidth)
	case "HEAD":
		return generateOpenApiOperationStrings("HEAD", pathItem.HEAD, m.contentWidth)
	case "OPTIONS":
		return generateOpenApiOperationStrings("OPTIONS", pathItem.OPTIONS, m.contentWidth)
	}
	return ""
}

func (m *model) RemoveResourceBinding() {
	path, ok := state.AppState.ApiContract.Paths[m.selectedPath]
	if !ok {
		return
	}

	switch m.CurrentMenuSelection() {
	case "GET":
		if path.GET != nil {
			path.GET.XResource = nil
		}
	case "POST":
		if path.POST != nil {
			path.POST.XResource = nil
		}
	case "PUT":
		if path.PUT != nil {
			path.PUT.XResource = nil
		}
	case "PATCH":
		if path.PATCH != nil {
			path.PATCH.XResource = nil
		}
	case "DELETE":
		if path.DELETE != nil {
			path.DELETE.XResource = nil
		}
	case "HEAD":
		if path.HEAD != nil {
			path.HEAD.XResource = nil
		}
	case "OPTIONS":
		if path.OPTIONS != nil {
			path.OPTIONS.XResource = nil
		}
	default:
		return
	}

	state.AppState.ApiContract.Paths[m.selectedPath] = path
	state.WriteToContractFile()
}

type BindResourceRow struct {
	Indices []int
	Values  []string
}

// TODO: Add some validation for rendering resource binding options based on the permissions (ex. if user selects a GET endpoint and resource table does not allow any fields to be readable, do not render)
func (m *model) GenerateBindResourceOptions() string {
	resources := state.AppState.ResourceContract.Resources

	grouped := make(map[string][]string)
	var resourceNames []string
	longest := 0

	for resName, resource := range resources {
		if resource.DB == nil {
			continue
		}

		resourceNames = append(resourceNames, resName)

		for tableName := range resource.DB.Tables {
			resourceTable := fmt.Sprintf("%s.%s", resName, tableName)
			grouped[resName] = append(grouped[resName], resourceTable)

			if len(resourceTable) > longest {
				longest = len(resourceTable)
			}
		}
	}

	total := 0
	for _, tables := range grouped {
		total += len(tables)
	}
	m.bindResourceCellLength = total

	sort.Strings(resourceNames)
	for _, resName := range resourceNames {
		sort.Strings(grouped[resName])
	}

	if longest == 0 {
		m.bindResourceRows = nil
		return ""
	}

	cols := max(1, (m.contentWidth/longest)-1)
	colWidth := max(8, (m.contentWidth/cols)-3)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0087ff")).
		MarginTop(1).
		MarginBottom(1)

	var sections []string
	selectedIndex := 0

	// reset and rebuild navigation rows every render
	m.bindResourceRows = nil

	for _, resName := range resourceNames {
		var sectionParts []string

		sectionParts = append(sectionParts, titleStyle.Render(resName))

		tables := grouped[resName]
		var rows []string

		for i := 0; i < len(tables); i += cols {
			end := min(i+cols, len(tables))
			var cells []string

			rowMeta := BindResourceRow{}

			for j := i; j < end; j++ {
				cells = append(cells,
					renderResourceBindingOptionsCell(
						tables[j],
						selectedIndex == m.selectedBindResourceCell,
						colWidth,
					),
				)

				rowMeta.Indices = append(rowMeta.Indices, selectedIndex)
				rowMeta.Values = append(rowMeta.Values, tables[j])

				selectedIndex++
			}

			m.bindResourceRows = append(m.bindResourceRows, rowMeta)
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
		}

		sectionParts = append(sectionParts, lipgloss.JoinVertical(lipgloss.Top, rows...))
		sections = append(sections, lipgloss.JoinVertical(lipgloss.Top, sectionParts...))
	}

	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

func renderResourceBindingOptionsCell(content string, selected bool, width int) string {
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Height(3).
		Width(max(1, width-1)).
		Align(lipgloss.Center, lipgloss.Center)

	if selected {
		style = style.BorderForeground(lipgloss.Color("#0087ff"))
	}

	return style.Render(content)
}

func (m *model) findBindResourcePosition() (row int, col int, ok bool) {
	for r, bindRow := range m.bindResourceRows {
		for c, idx := range bindRow.Indices {
			if idx == m.selectedBindResourceCell {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}

func (m *model) currentSelectedBindResource() string {
	for _, bindRow := range m.bindResourceRows {
		for i, idx := range bindRow.Indices {
			if idx == m.selectedBindResourceCell {
				return bindRow.Values[i]
			}
		}
	}
	return ""
}

func (m *model) BindResourceToEndpoint(selectedResource string) {
	if selectedResource == "" || !strings.Contains(selectedResource, ".") {
		return
	}
	s := strings.Split(selectedResource, ".")
	resourceName, tableName := s[0], s[1]

	resource, ok := state.AppState.ResourceContract.Resources[resourceName]
	if !ok || resource.DB == nil {
		return
	}

	_, ok = resource.DB.Tables[tableName]
	if !ok {
		return
	}

	path, ok := state.AppState.ApiContract.Paths[m.selectedPath]
	if !ok {
		return
	}

	resourceBinding := &contracts.RouteResourceBinding{
		ResourceName: resourceName,
		Table:        tableName,
	}

	method := m.CurrentMenuSelection()
	switch method {
	case "GET":
		if path.GET != nil {
			path.GET.XResource = resourceBinding
		}
	case "POST":
		if path.POST != nil {
			path.POST.XResource = resourceBinding
		}
	case "PUT":
		if path.PUT != nil {
			path.PUT.XResource = resourceBinding
		}
	case "PATCH":
		if path.PATCH != nil {
			path.PATCH.XResource = resourceBinding
		}
	case "DELETE":
		if path.DELETE != nil {
			path.DELETE.XResource = resourceBinding
		}
	case "HEAD":
		if path.HEAD != nil {
			path.HEAD.XResource = resourceBinding
		}
	case "OPTIONS":
		if path.OPTIONS != nil {
			path.OPTIONS.XResource = resourceBinding
		}
	default:
		return
	}

	state.AppState.ApiContract.Paths[m.selectedPath] = path

	state.WriteToContractFile()
}
