package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type Console struct {
	Debug, JSON bool
}

type Message struct {
	Status, Message string
}

// Blue outputs the requested text as blue.
func (c *Console) Blue(output string) string {
	if c.JSON {
		return output
	}

	return aurora.Blue(output).String()
}

// Bold outputs the requested text as bold.
func (c *Console) Bold(output string) string {
	if c.JSON {
		return output
	}

	return aurora.Bold(output).String()
}

// Error displays the error message and a panic if needed.
func (c *Console) Error(err error) {
	if c.JSON {
		message := Message{
			Status:  "Error",
			Message: err.Error(),
		}

		str, _ := json.Marshal(message)

		fmt.Println(string(str))
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", aurora.Bold(aurora.Red("[Error]")), err)

		if c.Debug {
			c.Println("")
			panic(err)
		}
	}

	os.Exit(1)
}

// Green outputs the requested text as green.
func (c *Console) Green(output string) string {
	if c.JSON {
		return output
	}

	return aurora.Green(output).String()
}

// Printf is a temporary wrapper on fmt.Printf.
func (c *Console) Printf(format string, a ...any) {
	if c.JSON {
		message := Message{
			Status:  "Info",
			Message: fmt.Sprintf(format, a...),
		}

		str, _ := json.Marshal(message)

		fmt.Println(string(str))
	} else {
		fmt.Printf(format, a...)
	}
}

// Println is a temporary wrapper on fmt.Println.
func (c *Console) Println(output string) {
	if c.JSON {
		message := Message{
			Status:  "Info",
			Message: output,
		}

		str, _ := json.Marshal(message)

		fmt.Println(string(str))
	} else {
		fmt.Println(output)
	}
}

// PromptConfirm asks the user to confirm output.
func (c *Console) PromptConfirm(promptText string, def bool) bool {
	if c.JSON {
		return def
	}

	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", promptText, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

// Success displays a formatted success message on successful completion of the command.
func (c *Console) Success(output string) {
	if c.JSON {
		message := Message{
			Status:  "Success",
			Message: output,
		}

		str, _ := json.Marshal(message)

		fmt.Println(string(str))
	} else {
		fmt.Printf("%s %s\n", aurora.Bold(aurora.Green("[Success]")), output)
	}
}

// Warn displays a formatted warning message.
func (c *Console) Warn(output string) {
	if c.JSON {
		message := Message{
			Status:  "Warning",
			Message: output,
		}

		str, _ := json.Marshal(message)

		fmt.Println(string(str))
	} else {
		fmt.Printf("%s %s\n", aurora.Bold(aurora.Yellow("[Warning]")), output)
	}
}

// Yellow outputs the requested text as yellow.
func (c *Console) Yellow(output string) string {
	if c.JSON {
		return output
	}

	return aurora.Yellow(output).String()
}
