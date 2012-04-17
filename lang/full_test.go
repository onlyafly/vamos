package lang

import (
	"../util"
	"io/ioutil"
	"strings"
	"testing"
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
		if strings.HasSuffix(name, ".in") {
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

	expected, errOut := util.ReadFile(outFilePath)
	if errOut != nil {
		t.Errorf("Error reading file <" + outFilePath + ">: " + errOut.Error())
		return
	}

	// Remove any carriage return line endings from .out file
	expected = strings.Replace(expected, "\r", "", -1)

	nodes, errors := Parse(input)
	if errors.Len() != 0 {
		verify(t, testNumber, input, expected, errors.String())
	} else {
		e := NewMapEnv()

		var result Node
		for _, n := range nodes {
			result = Eval(e, n)
		}

		actual := result.String()
		verify(t, testNumber, input, expected, actual)
	}
}

func verify(t *testing.T, testNumber, input, expected, actual string) {
	if expected != actual {
		t.Errorf(
			"TEST CASE #%s FAILED...\n"+
				"<<<<<INPUT>>>>>\n%v\n"+
				"<<<<<EXPECTED>>>>>\n%v\n"+
				"<<<<<ACTUAL>>>>>\n%v\n"+
				"<<<<<END>>>>>\n",
			testNumber,
			input,
			expected,
			actual)
	}
}
