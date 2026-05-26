package tasks

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frozenkro/sqncr/internal/models"
	"github.com/frozenkro/sqncr/internal/tui/viewmodels"
)

type TaskItemDelegate struct{}

func (d TaskItemDelegate) Height() int { return 1 }

func (d TaskItemDelegate) Spacing() int { return 0 }

func (d TaskItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d TaskItemDelegate) Render(w io.Writer, m list.Model, index int, li list.Item) {
	i, ok := li.(viewmodels.TaskItem)
	if !ok {
		return
	}

	var str strings.Builder

	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// Owner badge
	ownerBadge := ""
	switch i.Owner {
	case models.AiAgent:
		ownerBadge = "[A]"
	case models.Collabor:
		ownerBadge = "[C]"
	default:
		ownerBadge = "[U]"
	}

	// Blocker indicator
	blocker := ""
	if i.BlockerText != nil && *i.BlockerText != "" {
		blocker = "(!)"
	}

	// Date badge
	dateBadge := ""
	if i.ScheduledDate != nil {
		dateBadge = fmt.Sprintf(" [%s]", *i.ScheduledDate)
	}

	var icon string
	if i.Status == models.Completed {
		icon = "✅"
	} else if i.Status == models.InProgress {
		icon = "🔄"
	} else {
		icon = "○"
	}

	line := fmt.Sprintf("%s %s%s %s%s", ownerBadge, blocker, icon, i.Title(), dateBadge)

	if index == m.Index() {
		str.WriteString(selectedStyle.Render(line))
	} else {
		str.WriteString(normalStyle.Render(line))
	}

	fmt.Fprint(w, str.String())
}
