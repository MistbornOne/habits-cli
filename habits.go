package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var headerBox = lipgloss.NewStyle().
	Bold(true).
	Render("      💪🏼 Habit Tracker")

type model struct {
	habits              []string
	cursor              int
	selected            map[int]bool
	store               HabitStore
	manualOverride      bool
	promptingdate       bool
	promptingcurrent    bool
	promptinglongest    bool
	datepicker          textinput.Model
	currentstreakpicker textinput.Model
	longeststreakpicker textinput.Model
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
		"🏋🏼‍♂️ Workout",
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

	tiDate := textinput.New()
	tiDate.Placeholder = "YYYY-MM-DD"
	//tiDate.Focus()
	tiDate.CharLimit = 10
	tiDate.Width = 20

	tiCurrentStreak := textinput.New()
	tiCurrentStreak.Placeholder = "10"
	//tiCurrentStreak.Focus()
	tiCurrentStreak.CharLimit = 3
	tiCurrentStreak.Width = 20

	tiLongestStreak := textinput.New()
	tiLongestStreak.Placeholder = "10"
	//tiLongestStreak.Focus()
	tiLongestStreak.CharLimit = 3
	tiLongestStreak.Width = 20

	return model{
		habits:              habits,
		selected:            selected,
		store:               store,
		datepicker:          tiDate,
		currentstreakpicker: tiCurrentStreak,
		longeststreakpicker: tiLongestStreak,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.promptingdate {
		var cmd tea.Cmd
		m.datepicker, cmd = m.datepicker.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.logSpecificDate(strings.TrimSpace(m.datepicker.Value()))
				m.datepicker.Reset()
				m.promptingdate = false
			case "esc":
				m.datepicker.Reset()
				m.promptingdate = false
			}
		}
		return m, cmd
	} else if m.promptingcurrent {
		var cmd tea.Cmd
		m.currentstreakpicker, cmd = m.currentstreakpicker.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.logSpecificCurrent(strings.TrimSpace(m.currentstreakpicker.Value()))
				m.currentstreakpicker.Reset()
				m.promptingcurrent = false
			case "esc":
				m.currentstreakpicker.Reset()
				m.promptingcurrent = false
			}
		}
		return m, cmd

	} else if m.promptinglongest {
		var cmd tea.Cmd
		m.longeststreakpicker, cmd = m.longeststreakpicker.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.logSpecificLongest(strings.TrimSpace(m.longeststreakpicker.Value()))
				m.longeststreakpicker.Reset()
				m.promptinglongest = false
			case "esc":
				m.longeststreakpicker.Reset()
				m.promptinglongest = false
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
			m.promptingdate = true
			m.datepicker.Focus()
			m.currentstreakpicker.Blur()
			m.longeststreakpicker.Blur()
		case "c":
			m.promptingcurrent = true
			m.currentstreakpicker.Focus()
			m.datepicker.Blur()
			m.longeststreakpicker.Blur()
		case "l":
			m.promptinglongest = true
			m.longeststreakpicker.Focus()
			m.currentstreakpicker.Blur()
			m.datepicker.Blur()
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

			if !m.manualOverride {
				yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
				entry.Streak = calculateStreakFrom(yesterday, newDates)

			}

			entry.Longest = max(entry.Longest, calculateLongestStreak(newDates))
			m.store[habit] = entry
		}
	} else {
		m.selected[m.cursor] = true
		entry := m.store[habit]
		if !contains(entry.Dates, todayStr) {
			entry.Dates = append(entry.Dates, todayStr)

			if !m.manualOverride {

				entry.Streak = calculateStreakFrom(todayStr, entry.Dates)
			}
			entry.Longest = max(entry.Longest, calculateLongestStreak(entry.Dates))
			m.store[habit] = entry
		}
	}
	saveHabits(m.store)
	m.manualOverride = false
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

func (m *model) logSpecificCurrent(value string) {
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return // silently ignore or show feedback later
	}
	habit := m.habits[m.cursor]
	entry := m.store[habit]
	m.manualOverride = true
	entry.Streak = n
	//entry.Longest = n // keep logic consistent
	m.store[habit] = entry
	saveHabits(m.store)
}

func (m *model) logSpecificLongest(value string) {
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return
	}
	habit := m.habits[m.cursor]
	entry := m.store[habit]
	entry.Longest = n // just in case
	m.store[habit] = entry
	saveHabits(m.store)
}

func (m model) View() string {
	if m.promptingdate {
		return fmt.Sprintf("\nLog date for %s\n%s\n[enter] to confirm, [esc] to cancel\n", m.habits[m.cursor], m.datepicker.View())
	} else if m.promptingcurrent {
		return fmt.Sprintf("\nLog current streak for %s\n%s\n[enter] to confirm, [esc] to cancel\n", m.habits[m.cursor], m.currentstreakpicker.View())
	} else if m.promptinglongest {
		return fmt.Sprintf("\nLog longest streak for %s\n%s\n[enter] to confirm, [esc] to cancel\n", m.habits[m.cursor], m.longeststreakpicker.View())
	}

	s := "\n" + headerBox + "\n\n\n"
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
	s += "\n[Space] Toggle Today | [y] Log Yesterday | [d] Log Other Date |\n[c] Log Current Streak | [l] Log Longest Streak | [q] Quit\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
