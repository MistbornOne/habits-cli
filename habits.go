
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	habits []string
	cursor int
	selected map[int]bool
}

func initialModel() model {
	return model{
		habits: []string{"💧 Water Before Coffee", "☀️ Morning Pages", "✝️ Read Bible", "😍 Gratitude Practice", "👨🏼‍💻 Coding", "🇯🇵 Japanese", "📚 Read"},

		selected: make(map[int]bool),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.habits)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = true
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "💪🏼🔥 Habit Tracker\n\n"

	for i, habits := range m.habits {
		cursor := " "

		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, habits)
	}

	s += "\nPress q to quit. \n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)

		os.Exit(1)
	}
}
