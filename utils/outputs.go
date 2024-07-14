package utils

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	tableCellStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	tableBorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Table struct {
	Headers        []string
	Rows           [][]string
	cursor         int
	viewportHeight int
}

func (t Table) Init() tea.Cmd {
	return nil
}

func (t Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "enter":
			return t, tea.Quit
		case "up":
			if t.cursor > 0 {
				t.cursor--
			}
		case "down":
			if t.cursor < len(t.Rows)-1 {
				t.cursor++
			}
		}
	}
	return t, nil
}

func (t Table) View() string {
	var sb strings.Builder

	// Calculate column widths
	columnWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		columnWidths[i] = len(header)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	// Create top border
	sb.WriteString(createBorder(columnWidths, "┌", "┬", "┐"))
	sb.WriteString("\n")

	// Write headers
	for i, header := range t.Headers {
		sb.WriteString("│ ")
		sb.WriteString(tableHeaderStyle.Render(padRight(header, columnWidths[i])))
		sb.WriteString(" ")
	}
	sb.WriteString("│\n")

	// Create header-content separator
	sb.WriteString(createBorder(columnWidths, "├", "┼", "┤"))
	sb.WriteString("\n")

	// Write content
	start := max(0, t.cursor-t.viewportHeight/2)
	end := min(len(t.Rows), start+t.viewportHeight)
	for i, row := range t.Rows[start:end] {
		for j, cell := range row {
			sb.WriteString("│ ")
			if i+start == t.cursor {
				sb.WriteString(tableCellStyle.Reverse(true).Render(padRight(cell, columnWidths[j])))
			} else {
				sb.WriteString(tableCellStyle.Render(padRight(cell, columnWidths[j])))
			}
			sb.WriteString(" ")
		}
		sb.WriteString("│\n")
	}

	// Create bottom border
	sb.WriteString(createBorder(columnWidths, "└", "┴", "┘"))

	return sb.String()
}

func createBorder(widths []int, left, middle, right string) string {
	var parts []string
	for _, w := range widths {
		parts = append(parts, strings.Repeat("─", w+2))
	}
	return tableBorderStyle.Render(fmt.Sprintf("%s%s%s", left, strings.Join(parts, middle), right))
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func DisplayTable(headers []string, rows [][]string) error {
	table := Table{
		Headers:        headers,
		Rows:           rows,
		viewportHeight: 10, // Adjust this value to change the number of visible rows
	}
	p := tea.NewProgram(table)
	_, err := p.Run()
	return err
}
