package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	todos   []string
	cursor  int
	checked map[int]struct{}
}

func testModel() model {
	return model{
		todos:   []string{"Buy milk", "Walk the dog", "Write code"},
		cursor:  0,
		checked: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}

		case "enter", "space":
			if _, ok := m.checked[m.cursor]; ok {
				delete(m.checked, m.cursor)
			} else {
				m.checked[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Todos: \n\n"
	for i, todo := range m.todos {
		if m.cursor == i {
			s += "> "
		} else {
			s += "  "
		}

		if _, checked := m.checked[i]; checked {
			s += "[x] " + todo + "\n"
		} else {
			s += "[ ] " + todo + "\n"
		}
	}
	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(testModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
