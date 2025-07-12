package constants

import tea "github.com/charmbracelet/bubbletea"

// Key mappings
const (
	KeyQuit1   = "ctrl+c"
	KeyQuit2   = "q"
	KeyUp1     = "up"
	KeyUp2     = "k"
	KeyDown1   = "down"
	KeyDown2   = "j"
	KeyLeft1   = "left"
	KeyLeft2   = "h"
	KeyRight1  = "right"
	KeyRight2  = "l"
	KeySelect1 = "enter"
	KeySelect2 = " "
)

// KeyMap returns true if the key pressed matches any of the specified keys
func KeyMatch(msg tea.KeyMsg, keys ...string) bool {
	for _, key := range keys {
		if msg.String() == key {
			return true
		}
	}
	return false
}
