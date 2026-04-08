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
		methodCards = append(methodCards, generateOpenApiOperationStrings("GET", pathItem.GET, width))
	}
	if pathItem.POST != nil {
		methodCards = append(methodCards, generateOpenApiOperationStrings("POST", pathItem.POST, width))
	}
	if pathItem.PUT != nil {
		methodCards = append(methodCards, generateOpenApiOperationStrings("PUT", pathItem.PUT, width))
	}
	if pathItem.PATCH != nil {
		methodCards = append(methodCards, generateOpenApiOperationStrings("PATCH", pathItem.PATCH, width))
	}
	if pathItem.DELETE != nil {
		methodCards = append(methodCards, generateOpenApiOperationStrings("DELETE", pathItem.DELETE, width))
	}
	if pathItem.HEAD != nil {
		methodCards = append(methodCards, generateOpenApiOperationStrings("HEAD", pathItem.HEAD, width))
	}
	if pathItem.OPTIONS != nil {
		methodCards = append(methodCards, generateOpenApiOperationStrings("OPTIONS", pathItem.OPTIONS, width))
	}

	methodColumn := lipgloss.JoinVertical(lipgloss.Center, methodCards...)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(path),
		methodColumn,
	)
}

func generateOpenApiOperationStrings(method string, operation *contracts.OpenApiOperation, width int) string {
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
