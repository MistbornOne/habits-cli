package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	habits     []string
	cursor     int
	selected   map[int]bool
	store      HabitStore
	prompting  bool
	datepicker textinput.Model
}

func initialModel() model {
	store, _ := loadHabits()
	selected := make(map[int]bool)
	habits := []string{
		"💧 Water Before Coffee",
		"☀️ Morning Pages",
		"✝️ Read Bible",
		"😍 Gratitude Practice",
		"👨🏼‍💻 Coding",
		"🇯🇵 Japanese",
		"📚 Read",
	}
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

	ti := textinput.New()
	ti.Placeholder = "YYYY-MM-DD"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	return model{
		habits:     habits,
		selected:   selected,
		store:      store,
		datepicker: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.prompting {
		var cmd tea.Cmd
		m.datepicker, cmd = m.datepicker.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.logSpecificDate(strings.TrimSpace(m.datepicker.Value()))
				m.datepicker.Reset()
				m.prompting = false
			case "esc":
				m.datepicker.Reset()
				m.prompting = false
			}
		}
		return m, cmd
	}

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
			m.toggleToday()
		case "y":
			m.logDate(-1) // log yesterday
		case "d":
			m.prompting = true
			m.datepicker.Focus()
		}
	}
	return m, nil
}

func (m *model) toggleToday() {
	todayStr := today()
	habit := m.habits[m.cursor]
	if _, ok := m.selected[m.cursor]; ok {
		delete(m.selected, m.cursor)
		if entry, exists := m.store[habit]; exists {
			newDates := []string{}
			for _, d := range entry.Dates {
				if d != todayStr {
					newDates = append(newDates, d)
				}
			}
			entry.Dates = newDates
			yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			entry.Streak = calculateStreakFrom(yesterday, newDates)
			entry.Longest = max(entry.Longest, calculateLongestStreak(newDates))
			m.store[habit] = entry
		}
	} else {
		m.selected[m.cursor] = true
		entry := m.store[habit]
		if !contains(entry.Dates, todayStr) {
			entry.Dates = append(entry.Dates, todayStr)
			entry.Streak = calculateStreakFrom(todayStr, entry.Dates)
			entry.Longest = max(entry.Longest, calculateLongestStreak(entry.Dates))
			m.store[habit] = entry
		}
	}
	saveHabits(m.store)
}

func (m *model) logDate(daysAgo int) {
	habit := m.habits[m.cursor]
	dateStr := time.Now().AddDate(0, 0, daysAgo).Format("2006-01-02")
	entry := m.store[habit]
	if !contains(entry.Dates, dateStr) {
		entry.Dates = append(entry.Dates, dateStr)
		entry.Streak = calculateStreakFrom(today(), entry.Dates)
		entry.Longest = calculateLongestStreak(entry.Dates)
		m.store[habit] = entry
		saveHabits(m.store)
	}
}

func (m *model) logSpecificDate(dateStr string) {
	habit := m.habits[m.cursor]
	entry := m.store[habit]
	if !contains(entry.Dates, dateStr) {
		entry.Dates = append(entry.Dates, dateStr)
		entry.Streak = calculateStreakFrom(today(), entry.Dates)
		entry.Longest = calculateLongestStreak(entry.Dates)
		m.store[habit] = entry
		saveHabits(m.store)
	}
}

func (m model) View() string {
	if m.prompting {
		return fmt.Sprintf("\nLog date for %s\n%s\n[enter] to confirm, [esc] to cancel\n", m.habits[m.cursor], m.datepicker.View())
	}

	s := "      \n💪 Habit Tracker\n\n"
	for i, habit := range m.habits {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "✔︎"
		}
		entry := m.store[habit]
		s += fmt.Sprintf("%s [%s] %s (%d 🔥 / %d 🏆)\n", cursor, checked, habit, entry.Streak, entry.Longest)
	}
	s += "\n[Space] toggle today | [y] log yesterday | [d] log other date | [q] quit\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

