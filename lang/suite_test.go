package lang

import (
	"io/ioutil"
	"strings"
	"testing"

	"../util"
)

const (
	testsuiteDir = "../testsuite"
)

func TestFullSuite(t *testing.T) {
	fileInfos, err := ioutil.ReadDir(testsuiteDir)
	if err != nil {
		t.Errorf("Unable to find test suite directory")
		return
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if strings.HasSuffix(name, ".v") {
			testInputFile(name, t)
		}
	}
}

func testInputFile(inFileName string, t *testing.T) {
	parts := strings.Split(inFileName, ".")
	testNumber := parts[0]

	inFilePath := testsuiteDir + "/" + inFileName
	outFilePath := testsuiteDir + "/" + testNumber + ".out"

	input, errIn := util.ReadFile(inFilePath)
	if errIn != nil {
		t.Errorf("Error reading file <" + inFilePath + ">: " + errIn.Error())
		return
	}

	expectedRaw, errOut := util.ReadFile(outFilePath)
	if errOut != nil {
		t.Errorf("Error reading file <" + outFilePath + ">: " + errOut.Error())
		return
	}

	// Remove any carriage return line endings from .out file
	expectedWithUntrimmed := strings.Replace(expectedRaw, "\r", "", -1)
	expected := strings.TrimSpace(expectedWithUntrimmed)

	nodes, errors := Parse(input)
	if errors.Len() != 0 {
		verify(t, testNumber, input, expected, errors.String())
	} else {
		e := NewTopLevelMapEnv()

		var result Node
		var evalError error
		for _, n := range nodes {
			result, evalError = Eval(e, n)
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
		verify(t, testNumber, input, expected, actual)
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
