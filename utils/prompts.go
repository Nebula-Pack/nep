package utils

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Item struct {
	TitleText, Desc string
}

func (i Item) Title() string       { return i.TitleText }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.TitleText }

// GroupedTextInput collects multiple text inputs in a single prompt
func GroupedTextInput(prompts []string) ([]string, error) {
	models := make([]textinput.Model, len(prompts))
	for i := range prompts {
		ti := textinput.New()
		ti.Placeholder = prompts[i]
		if i == 0 {
			ti.Focus()
		}
		models[i] = ti
	}

	m := groupedTextInputModel{inputs: models}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	finalTextInputModel, ok := finalModel.(groupedTextInputModel)
	if !ok {
		return nil, fmt.Errorf("could not get user input")
	}

	results := make([]string, len(finalTextInputModel.inputs))
	for i, input := range finalTextInputModel.inputs {
		results[i] = input.Value()
	}

	return results, nil
}

type groupedTextInputModel struct {
	inputs   []textinput.Model
	active   int
	quitting bool
}

func (m groupedTextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m groupedTextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.active == len(m.inputs)-1 {
				m.inputs[m.active].Blur()
				m.quitting = true
				return m, tea.Quit
			}
			m.inputs[m.active].Blur()
			m.active++
			m.inputs[m.active].Focus()
			return m, textinput.Blink
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}

	if !m.quitting {
		var cmd tea.Cmd
		m.inputs[m.active], cmd = m.inputs[m.active].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m groupedTextInputModel) View() string {
	var sb strings.Builder
	sb.WriteString("\n(press enter to move to the next field, ctrl+c or esc to quit)\n")
	for i, input := range m.inputs {
		style := blurredStyle
		if i == m.active && !m.quitting {
			style = focusedStyle
		}
		sb.WriteString(style.Render(input.View()) + "\n")
	}
	return sb.String()
}

// SelectFromList creates an interactive list for selection
func SelectFromList(title string, items []Item, itemsToShow int) (string, error) {
	const listWidth = 60

	l := list.New(toListItems(items), list.NewDefaultDelegate(), listWidth, (itemsToShow+1)*3)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	m := listModel{list: l}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	finalListModel, ok := finalModel.(listModel)
	if !ok {
		return "", fmt.Errorf("could not get user selection")
	}

	selectedItem, ok := finalListModel.list.SelectedItem().(Item)
	if !ok {
		return "", fmt.Errorf("selected item type assertion failed")
	}

	return selectedItem.FilterValue(), nil
}

func toListItems(items []Item) []list.Item {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}
	return listItems
}

type listModel struct {
	list list.Model
	err  error
}

func (m listModel) Init() tea.Cmd { return nil }

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" || msg.String() == "enter" {
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return "\n" + m.list.View()
}
