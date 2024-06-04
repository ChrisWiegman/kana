package console

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsole_Blue(t *testing.T) {
	console := &Console{}
	output := console.Blue("Hello, World!")
	expected := "\x1b[34mHello, World!\x1b[0m"
	assert.Equal(t, expected, output)
}
func TestConsole_Bold(t *testing.T) {
	console := &Console{}
	output := console.Bold("Hello, World!")
	expected := "\x1b[1mHello, World!\x1b[0m"
	assert.Equal(t, expected, output)
}

func TestConsole_Green(t *testing.T) {
	console := &Console{}
	output := console.Green("Hello, World!")
	expected := "\x1b[32mHello, World!\x1b[0m"
	assert.Equal(t, expected, output)
}

func TestConsole_Yellow(t *testing.T) {
	console := &Console{}
	output := console.Yellow("Hello, World!")
	expected := "\x1b[33mHello, World!\x1b[0m"
	assert.Equal(t, expected, output)
}
