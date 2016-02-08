package main

import (
	"bufio"
	"flag"
	"fmt"

	"log"
	"os"
	"strings"

	"vamos/util"

	"vamos/lang/ast"
	"vamos/lang/interpreter"

	"vamos/Godeps/_workspace/src/github.com/peterh/liner"
)

const (
	version         = `0.1.0-alpha`
	versionDate     = `2015-02-20`
	historyFilename = "/tmp/.vamos_liner_history"
)

var (
	// TODO add functionality for these missing commands
	commandCompletions = []string{":quit" /*":load ", ":reset", ":help",*/, ":inspect "}
	// TODO wordCompletions    = []string{"def", "update!"}
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

func startLiner() *liner.State {
	line := liner.NewLiner()
	openLinerHistory(line)
	configureLiner(line)
	return line
}

func stopLiner(line *liner.State) {
	line.Close()
}

func main() {

	startupFileName := flag.String("l", "", "load a file at startup")
	showHelp := flag.Bool("help", false, "show the help")
	flag.Parse()
	exeFileName := flag.Arg(0)

	if showHelp != nil && *showHelp {
		fmt.Printf("Usage of vamos:\n")
		flag.PrintDefaults()
		return
	}

	// Setup liner

	line := startLiner()
	defer line.Close()

	standardReadLine := func() string {
		stopLiner(line) // Liner must be stopped to turn off raw mode...
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		line = startLiner() // ...and then turn it back on.
		return text
	}

	// Initialize

	topLevelEnv := interpreter.NewTopLevelMapEnv()

	if len(exeFileName) != 0 {
		loadFile("prelude.v", topLevelEnv, standardReadLine)
		loadFile(exeFileName, topLevelEnv, standardReadLine)
		return
	}

	fmt.Printf("Vamos %s (%s)\n", version, versionDate)
	loadFile("prelude.v", topLevelEnv, standardReadLine)
	fmt.Printf("(Press Ctrl+C or type :quit to exit)\n\n")

	// Loading of files

	if startupFileName != nil {
		loadFile(*startupFileName, topLevelEnv, standardReadLine)
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

		line.AppendHistory(input)
		writeLinerHistory(line)

		switch {
		case input == ":test":
			fmt.Print("Test Prompt: ")
			fmt.Println(standardReadLine())
		case input == ":quit":
			return
		case strings.HasPrefix(input, ":inspect "):
			withoutInspectPrefix := strings.Split(input, ":inspect ")[1]
			if result, err := interpreter.ParseEval(topLevelEnv, withoutInspectPrefix, standardReadLine, "REPL"); err == nil {
				inspect(result)
			} else {
				fmt.Println(err.Error())
			}
		default:
			interpreter.ParseEvalPrint(topLevelEnv, input, standardReadLine, "REPL", true)
		}
	}
}

func inspect(arg ast.Node) {
	switch val := arg.(type) {
	case *interpreter.EnvNode:
		fmt.Printf(
			"Environment\n  Name='%v'\n  Env=%v\n",
			val.Name(),
			val.Env.String())
	default:
		fmt.Printf("Don't know how to inspect: %v\n", val.String())
	}
}

func loadFile(fileName string, env interpreter.Env, standardReadLine func() string) {
	if len(fileName) > 0 {
		content, err := util.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Error while loading file <%v>: %v\n", fileName, err.Error())
		} else {
			interpreter.ParseEvalPrint(env, content, standardReadLine, fileName, false)
		}
	}
}
