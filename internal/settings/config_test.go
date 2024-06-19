package settings

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPrintJSONSettings(t *testing.T) {
	settings := &Settings{
		global: viper.New(),
		local:  viper.New(),
	}

	// Set some test settings
	settings.global.Set("setting1", "value1")
	settings.global.Set("setting2", "value2")
	settings.local.Set("setting3", "value3")
	settings.local.Set("setting4", "value4")

	// Create a buffer to capture the output
	buf := new(bytes.Buffer)

	// Temporarily set your file as the standard output (and save the old)
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function being tested
	printJSONSettings(settings)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func(buf *bytes.Buffer) {
		_, err := io.Copy(buf, r)
		assert.NoError(t, err)
		outC <- buf.String()
	}(buf)

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	// Verify the output
	expected := `{"Global":{"setting1":"value1","setting2":"value2"},"Local":{"setting3":"value3","setting4":"value4"}}`
	assert.Equal(t, expected, out)
}
