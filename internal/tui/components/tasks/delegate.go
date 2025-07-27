package tasks

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frozenkro/mcpsequencer/internal/tui/viewmodels"
)

type TaskItemDelegate struct{}

func (d TaskItemDelegate) Height() int { return 1 }

func (d TaskItemDelegate) Spacing() int { return 0 }

func (d TaskItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskItemDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	i, ok := li.(viewmodels.TaskView)
	if !ok {
		return
	}

	var str strings.Builder

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	var icon string
	if i.IsCompleted {
		icon = "✅"
	} else if i.IsInProgress {
		icon = "🔄"
	} else {
		icon = "○"
	}

	if index == m.Index() {
		str.WriteString(selectedStyle.Render(fmt.Sprintf("%s %s", icon, i.Title())))
	} else {
		str.WriteString(normalStyle.Render(fmt.Sprintf("%s %s", icon, i.Title())))
	}

	fmt.Fprint(w, str.String())
}
