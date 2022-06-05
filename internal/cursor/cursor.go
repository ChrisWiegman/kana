package cursor

import "fmt"

type Cursor struct{}

func (cursor *Cursor) Hide() {
	fmt.Printf("\033[?25l")
}

func (cursor *Cursor) Show() {
	fmt.Printf("\033[?25h")
}

func (cursor *Cursor) MoveUp(rows int) {
	fmt.Printf("\033[%dF", rows)
}

func (cursor *Cursor) MoveDown(rows int) {
	fmt.Printf("\033[%dE", rows)
}

func (cursor *Cursor) ClearLine() {
	fmt.Printf("\033[2K")
}
