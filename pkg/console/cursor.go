package console

import "fmt"

type Cursor struct{}

// ClearLine Clears the line
func (cursor *Cursor) ClearLine() {
	fmt.Printf("\033[2K")
}

// Hide Hides the cursor
func (cursor *Cursor) Hide() {
	fmt.Printf("\033[?25l")
}

// MoveDown Moves down a line
func (cursor *Cursor) MoveDown(rows int) {
	fmt.Printf("\033[%dE", rows)
}

// MoveUp Moves up a line
func (cursor *Cursor) MoveUp(rows int) {
	fmt.Printf("\033[%dF", rows)
}

// Show Shows the cursor
func (cursor *Cursor) Show() {
	fmt.Printf("\033[?25h")
}
