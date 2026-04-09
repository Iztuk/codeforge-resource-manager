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
