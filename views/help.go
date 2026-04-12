package views

import "charm.land/lipgloss/v2"

func GenerateGeneralHelpView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0087ff"))

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true)

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00d7ff"))

	descStyle := lipgloss.NewStyle()

	spacer := "\n"

	var b string

	// Title
	b += titleStyle.Render("Help & Navigation Guide")
	b += "\n\n"

	b += GenerateGeneralCommands()

	b += spacer + spacer

	// Page Specific Commands
	b += sectionStyle.Render("Page-Specific Commands") + "\n\n"

	b += descStyle.Render("Resource Management & Field Permissions") + "\n"
	b += keyStyle.Render("ctrl+a") + "      - Add resource or binding\n"
	b += keyStyle.Render("ctrl+d") + "      - Remove/Delete selected item\n"
	b += keyStyle.Render("space") + "       - Toggle field permissions (read/write/mutable)\n"

	b += "\n"

	b += descStyle.Render("Add Resource Form") + "\n"
	b += keyStyle.Render("tab / down") + "      - Move to next field\n"
	b += keyStyle.Render("shift+tab / up") + "  - Move to previous field\n"
	b += keyStyle.Render("enter") + "           - Save resource\n"
	b += keyStyle.Render("esc") + "             - Cancel form\n"

	b += "\n"

	b += descStyle.Render("Resource Binding") + "\n"
	b += keyStyle.Render("ctrl+a") + "      - Attach resource binding to endpoint\n"
	b += keyStyle.Render("ctrl+d") + "      - Remove resource binding from endpoint\n"

	b += spacer + spacer

	// Notes / Descriptions
	b += sectionStyle.Render("Notes") + "\n\n"

	b += descStyle.Render("- Commands may vary depending on the current page\n")
	b += descStyle.Render("- Use navigation keys to explore available actions\n")
	b += descStyle.Render("- Selected items are highlighted in the UI\n")

	return b
}

func GenerateGeneralCommands() string {
	var b string

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true)

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00d7ff"))

	// General Commands
	b += sectionStyle.Render("General Commands") + "\n\n"

	b += keyStyle.Render("h / j / k / l") + "   - Navigate through menus and content\n"
	b += keyStyle.Render("ctrl+h / ctrl+l") + " - Switch between panels\n"
	b += keyStyle.Render("ctrl+n / ctrl+p") + " - Scroll panel content\n"
	b += keyStyle.Render("enter") + "           - Select highlighted item\n"
	b += keyStyle.Render("backspace") + "       - Go back to previous view\n"
	b += keyStyle.Render("ctrl+c") + "          - Exit application\n"

	return b
}
