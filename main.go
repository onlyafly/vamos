package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"vamos/lang/ast"
	"vamos/lang/interpretation"
	"vamos/lang/parsing"
	"vamos/util"

	"github.com/peterh/liner"
)

const (
	version         = `0.2.0`
	versionDate     = `2015-01-30`
	historyFilename = "/tmp/.vamos_liner_history"
)

var (
	// TODO add functionality for these missing commands
	commandCompletions = []string{":quit", ":load", ":reset", ":help"}
	// TODO wordCompletions    = []string{"def", "set!"}
)

func configureLiner(linerState *liner.State) {
	linerState.SetCtrlCAborts(true)

	linerState.SetCompleter(func(line string) (c []string) {
		for _, n := range commandCompletions {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	/* TODO
	// WordCompleter takes the currently edited line with the cursor position and
	// returns the completion candidates for the partial word to be completed. If
	// the line is "Hello, wo!!!" and the cursor is before the first '!',
	// ("Hello, wo!!!", 9) is passed to the completer which may returns
	// ("Hello, ", {"world", "Word"}, "!!!") to have "Hello, world!!!".
	linerState.SetWordCompleter(func(line string, pos int) (head string, completions []string, tail string) {
		for _, n := range wordCompletions {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})
	*/
}

func openLinerHistory(line *liner.State) {
	if f, err := os.Open(historyFilename); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func writeLinerHistory(line *liner.State) {
	if f, err := os.Create(historyFilename); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func main() {

	fileName := flag.String("l", "", "load a file at startup")
	flag.Parse()

	// Setup liner

	line := liner.NewLiner()
	defer line.Close()
	openLinerHistory(line)
	configureLiner(line)

	// Initialize

	fmt.Printf("Vamos %s (%s)\n", version, versionDate)
	fmt.Printf("Press Ctrl+C or type :quit to exit\n\n")

	env := interpretation.NewTopLevelMapEnv()

	// Loading of files

	if fileName != nil && len(*fileName) > 0 {
		content, _ := util.ReadFile(*fileName)
		parseEval(env, content)
	}

	// REPL

	for {
		input, err := line.Prompt("> ")

		if err != nil {
			if err.Error() == "prompt aborted" {
				fmt.Printf("Quiting...\n")
			} else {
				fmt.Printf("Prompt error: %s\n", err)
			}
			return
		}

		if input == ":quit" {
			return
		}

		line.AppendHistory(input)
		writeLinerHistory(line)

		parseEval(env, input)
	}
}

func parseEval(env interpretation.Env, input string) {
	defer func() {
		// Some non-application triggered panic has occurred
		if e := recover(); e != nil {
			fmt.Printf("Host environment error: %v\n", e)
		}
	}()

	nodes, parseErrors := parsing.Parse(input)

	if parseErrors != nil {
		fmt.Println(parseErrors.String())
	}

	var result ast.Node
	var evalError error
	for _, n := range nodes {
		result, evalError = interpretation.Eval(env, n)
		if evalError != nil {
			break
		}
	}

	var actual string
	if evalError == nil {
		actual = result.String()
	} else {
		actual = evalError.Error()
	}

	fmt.Println(actual)
}
