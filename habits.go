
package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	habits []string
	cursor int
	selected map[int]bool
	store HabitStore
}

func initialModel() model {
	store, _ := loadHabits()

	selected := make(map[int]bool)
	habits := []string{"💧 Water Before Coffee", "☀️ Morning Pages", "✝️ Read Bible", "😍 Gratitude Practice", "👨🏼‍💻 Coding", "🇯🇵 Japanese", "📚 Read",}

	todayStr := today()
  
  for i, habit := range habits {
		if entry, ok := store[habit]; ok {
			for _, d := range entry.Dates {
				if d == todayStr {
					selected[i] = true
					break
				}
			}
		}
	}

	return model{
    habits: habits,
	  selected: selected,	
		store: store,

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
			todayStr := today()
			habit := m.habits[m.cursor]

			if ok {
				delete(m.selected, m.cursor)

				if entry, exists := m.store[habit]; exists {
					newDates := []string{}
					for _, d := range entry.Dates {
						if d != todayStr {
							newDates = append(newDates, d)
						}
					}

					entry.Dates = newDates
					entry.Streak = calculateStreak(newDates)
					m.store[habit] =  entry
				}
				
			} else {
				m.selected[m.cursor] = true

				entry := m.store[habit]
				if !contains(entry.Dates, todayStr) {
					entry.Dates = append(entry.Dates, todayStr)
					entry.Streak = calculateStreak(entry.Dates)
					m.store[habit] = entry
				}
			}
		saveHabits(m.store)
		}
	}

	return m, nil
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func calculateStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	dateMap := map[string]bool{}
	for _, d := range dates {
		dateMap[d] = true
	}

	streak := 0
	t := time.Now()
	for {
		ds := t.Format("2006-01-02")
		if dateMap[ds] {
			streak++
			t = t.AddDate(0,0, -1)
		} else{
			break
		}
	}
	return streak
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

		streak := m.store[habits].Streak
		s += fmt.Sprintf("%s [%s] %s (%d🔥)\n", cursor, checked, habits, streak)
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
