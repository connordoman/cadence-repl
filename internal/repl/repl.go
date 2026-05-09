package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/connordoman/cadence"
	"github.com/connordoman/windy"
	"github.com/ergochat/readline"
	"golang.org/x/term"
)

const replPrompt = ">>> "

var replStyle = lipgloss.NewStyle().Foreground(windy.Amber500.Glossy()).Bold(true)
var errorStyle = lipgloss.NewStyle().Foreground(windy.Red500.Glossy())
var verboseStyle = lipgloss.NewStyle().Foreground(windy.Purple500.Glossy()).Background(windy.Purple900.Glossy()).Bold(true).Padding(0, 1)

type Repl struct {
	Verbose       bool
	HumanReadable bool
}

func NewRepl(verbose bool, humanReadable bool) *Repl {
	return &Repl{
		Verbose:       verbose,
		HumanReadable: humanReadable,
	}
}

func (r *Repl) report(err error, msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)
	fmt.Println(errorStyle.Render(fmt.Sprintf("Error: %s: %v", message, err)))
}

func (r *Repl) log(msg string, args ...any) {
	if !r.Verbose {
		return
	}

	text := fmt.Sprintf(msg, args...)

	log.Printf("%s %s", verboseStyle.Render("verbose"), text)
}

func (r *Repl) formatDate(date time.Time) string {
	if r.HumanReadable {
		_, week := date.ISOWeek()
		return fmt.Sprintf("W%02d: %s", week, date.Format("Monday, 02 January 2006"))
	}
	return date.Format("2006-01-02")
}

func (r *Repl) handleLine(line string) (quit bool) {
	line = strings.TrimSpace(line)
	if line == "exit" {
		return true
	}

	results, err := cadence.Compile(line)
	if err != nil {
		r.report(err, "error compiling expression")
		return false
	}

	for _, result := range results {
		fmt.Println(r.formatDate(result))
	}
	return false
}

func (r *Repl) Run() error {
	r.log("Running in verbose mode")

	if r.HumanReadable {
		r.log("Running in human readable mode")
	}

	fmt.Println("Cadence REPL (type 'exit' to quit)")

	if term.IsTerminal(int(os.Stdin.Fd())) {
		if err := r.runInteractive(); err != nil {
			return err
		}
	} else {
		r.runLineOriented()
	}

	fmt.Println("Bye!")
	return nil
}

func (r *Repl) runInteractive() error {
	rl, err := readline.NewFromConfig(&readline.Config{
		Prompt: replStyle.Render(replPrompt),
	})
	if err != nil {
		return fmt.Errorf("repl: readline: %w", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF {
				break
			}
			if err == readline.ErrInterrupt {
				continue
			}
			r.report(err, "error reading input")
			continue
		}
		if r.handleLine(line) {
			break
		}
	}
	return nil
}

func (r *Repl) runLineOriented() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(replStyle.Render(replPrompt))
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			r.report(err, "error reading input")
			continue
		}
		if r.handleLine(line) {
			break
		}
	}
}
