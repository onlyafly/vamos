package lang

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vamos/lang/ast"
	"vamos/lang/interpretation"
	"vamos/lang/parsing"
	"vamos/util"
)

const (
	testsuiteDir = "../testsuite"
)

func TestFullSuite(t *testing.T) {
	filepath.Walk(testsuiteDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil // Can't visit this node, but continue walking elsewhere
		}
		if !!fi.IsDir() {
			return nil // Not a file, ignore.
		}

		name := fi.Name()
		matched, err := filepath.Match("*.v", name)
		if err != nil {
			return err // malformed pattern
		}

		if matched {
			testInputFile(fp, t)
		}

		return nil
	})
}

func testInputFile(sourceFilePath string, t *testing.T) {
	sourceDirPart, sourceFileNamePart := filepath.Split(sourceFilePath)
	parts := strings.Split(sourceFileNamePart, ".")
	testName := parts[0]

	outputFilePath := sourceDirPart + testName + ".out"

	input, errIn := util.ReadFile(sourceFilePath)
	if errIn != nil {
		t.Errorf("Error reading file <" + sourceFilePath + ">: " + errIn.Error())
		return
	}

	expectedRaw, errOut := util.ReadFile(outputFilePath)
	if errOut != nil {
		t.Errorf("Error reading file <" + outputFilePath + ">: " + errOut.Error())
		return
	}

	// Remove any carriage return line endings from .out file
	expectedWithUntrimmed := strings.Replace(expectedRaw, "\r", "", -1)
	expected := strings.TrimSpace(expectedWithUntrimmed)

	nodes, errors := parsing.Parse(input)
	if errors.Len() != 0 {
		verify(t, sourceFilePath, input, expected, errors.String())
	} else {
		e := interpretation.NewTopLevelMapEnv()

		var outputBuffer bytes.Buffer

		var result ast.Node
		var evalError error
		for _, n := range nodes {
			result, evalError = interpretation.Eval(e, n, &outputBuffer)
			if evalError != nil {
				break
			}
		}

		var actual string
		actual = (&outputBuffer).String()

		if evalError == nil {
			actual = actual + result.String()
		} else {
			actual = actual + evalError.Error()
		}
		verify(t, sourceFilePath, input, expected, actual)
	}
}

func verify(t *testing.T, testCaseName, input, expected, actual string) {
	if expected != actual {
		t.Errorf(
			"\n===== TEST SUITE CASE FAILED: %s\n"+
				"===== INPUT\n%v\n"+
				"===== EXPECTED\n%v\n"+
				"===== ACTUAL\n%v\n"+
				"===== END\n",
			testCaseName,
			input,
			expected,
			actual)
	}
}
