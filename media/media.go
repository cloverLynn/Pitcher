package media

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os/exec"
)

var style = lipgloss.NewStyle().Inline(true).Width(100)

type Media struct {
	Name    string
	Trailer string
	Poster  string
	Width   int
}

func (m *Media) View() string {
	if m.Poster == "" {
		return "No poster available"
	}

	cmd := exec.Command("catimg", "-w", "60", m.Poster)
	out, err := cmd.Output()
	if err != nil {
		return "Error rendering poster: " + err.Error()
	}

	name := style.Render(m.Name)
	return lipgloss.JoinVertical(lipgloss.Top, name, string(out))
}

func (m *Media) PlayTrailer() {
	app := "mpv"
	arg0 := "-fs"
	arg1 := m.Trailer
	cmd := exec.Command(app, arg0, arg1)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Error playing trailer:", err)
	}
}
