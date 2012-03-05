package lang

import (
	"strings"
	"os"
	"bufio"
)

////////// Utility Functions

func ReadLine() string {
	bufferedReader := bufio.NewReader(os.Stdin)
	line, _ := bufferedReader.ReadString('\n')
	return strings.TrimRight(line, "\n")
}
