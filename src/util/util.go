package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/cfmt"
	"github.com/fatih/color"
)

// DoesCommandOuputContain runs the given command
// and returns whether its output contains the given string.
func DoesCommandOuputContain(cmd []string, value string) bool {
	return strings.Contains(GetCommandOutput(cmd...), value)
}

// DoesCommandOuputContainLine runs the given command
// and returns whether its output contains the given string as an entire line.
func DoesCommandOuputContainLine(cmd []string, value string) bool {
	list := strings.Split(GetCommandOutput(cmd...), "\n")
	return DoesStringArrayContain(list, value)
}

// DoesStringArrayContain returns whether the given string slice
// contains the given string.
func DoesStringArrayContain(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

// ExitWithErrorMessage prints the given error message and terminates the application.
func ExitWithErrorMessage(messages ...string) {
	PrintError(messages...)
	os.Exit(1)
}

// GetCommandOutput runs the given command and returns its output.
func GetCommandOutput(cmd ...string) string {
	output, err := GetFullCommandOutput(cmd...)
	exit.IfWrapf(err, "Command: %s\nOutput: %s", strings.Join(cmd, " "), output)
	return strings.TrimSpace(output)
}

// GetFullCommandOutput runs the given command and returns its output and error
func GetFullCommandOutput(cmd ...string) (string, error) {
	subProcess := exec.Command(cmd[0], cmd[1:]...)
	output, err := subProcess.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

var openBrowserCommands = []string{
	"xdg-open",
	"open",
	"cygstart",
	"x-www-browser",
	"firefox",
	"opera",
	"mozilla",
	"netscape",
}
var missingOpenBrowserCommandMessages = []string{
	"Cannot open a browser.",
	"If you think this is a bug,",
	"please open an issue at https://github.com/Originate/git-town/issues",
	"and mention your OS and browser.",
}

// GetOpenBrowserCommand returns the command to run on the console
// to open the default browser.
func GetOpenBrowserCommand() string {
	if runtime.GOOS == "windows" {
		// NOTE: the "explorer" command cannot handle special characters
		//       like "?" and "=", so we are using "start" here.
		//       "?" can be escaped via "\", but "=" cannot.
		return "start"
	}
	for _, command := range openBrowserCommands {
		output, err := GetFullCommandOutput("which", command)
		if err == nil && output != "" {
			return command
		}
	}
	ExitWithErrorMessage(missingOpenBrowserCommandMessages...)
	return ""
}

var inputReader = bufio.NewReader(os.Stdin)

// GetUserInput reads input from the user and returns it.
func GetUserInput() string {
	text, err := inputReader.ReadString('\n')
	exit.IfWrap(err, "Error getting user input")
	return strings.TrimSpace(text)
}

// Indent outputs the given string with the given level of indentation
// on each line. Each level of indentation is two spaces.
func Indent(message string, level int) string {
	prefix := strings.Repeat("  ", level)
	return prefix + strings.Replace(message, "\n", "\n"+prefix, -1)
}

// Pluralize outputs the count and the word. The word is made plural
// if the count isn't one
func Pluralize(count, word string) string {
	result := count + " " + word
	if count != "1" {
		result = result + "s"
	}
	return result
}

// PrintError prints the given error message to the console.
func PrintError(messages ...string) {
	errHeaderFmt := color.New(color.Bold).Add(color.FgRed)
	errMessageFmt := color.New(color.FgRed)
	fmt.Println()
	_, err := errHeaderFmt.Println("  Error")
	exit.If(err)
	for _, message := range messages {
		_, err = errMessageFmt.Println("  " + message)
		exit.If(err)
	}
	fmt.Println()
}

// PrintLabelAndValue prints the label bolded and underlined
// the value indented on the next line
// followed by an empty line
func PrintLabelAndValue(label, value string) {
	labelFmt := color.New(color.Bold).Add(color.Underline)
	_, err := labelFmt.Println(label + ":")
	exit.If(err)
	cfmt.Println(Indent(value, 1))
	fmt.Println()
}

// RemoveStringFromSlice returns a new string slice which is the given string slice
// with the given string removed
func RemoveStringFromSlice(list []string, value string) (result []string) {
	for _, element := range list {
		if element != value {
			result = append(result, element)
		}
	}
	return
}
