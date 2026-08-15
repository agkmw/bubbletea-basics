package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var primaryColor = lipgloss.Color("#e35fc6")

type styles struct {
	app          lipgloss.Style
	title        lipgloss.Style
	queueItem    lipgloss.Style
	selectedItem lipgloss.Style
	help         lipgloss.Style
}

func newStyles() styles {
	return styles{
		app: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor),
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor),
		help: lipgloss.NewStyle().
			Foreground(lipgloss.Alpha(primaryColor, 0.1)),
		queueItem: lipgloss.NewStyle().Foreground(lipgloss.White),
		selectedItem: lipgloss.NewStyle().
			Background(primaryColor).
			Foreground(lipgloss.Black),
	}
}

type model struct {
	styles   styles
	songs    []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	m := model{
		styles:   newStyles(),
		songs:    []string{"Haru", "Setting Sun", "Plover"},
		selected: make(map[int]struct{}),
	}

	return m
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
			if m.cursor < len(m.songs)-1 {
				m.cursor++
			}

		case "enter", "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	header := m.styles.title.Render("Select a song to play.\n")

	queue := "\n"
	for i, song := range m.songs {
		item := fmt.Sprintf("%2d. %-15s Eve ", i+1, song)
		if i == m.cursor {
			queue += m.styles.selectedItem.MaxWidth(40).Render(item)
		} else {
			queue += m.styles.queueItem.MaxWidth(40).Render(item)
		}
		queue += "\n"
	}

	help := m.styles.help.Render("\nPress q to quit.\n")

	s := m.styles.app.Render(header, queue, help)
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
