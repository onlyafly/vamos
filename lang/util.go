package lang

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

////////// Utility Functions

func ReadLine() string {
	bufferedReader := bufio.NewReader(os.Stdin)
	line, _ := bufferedReader.ReadString('\n')
	return strings.TrimRight(line, "\n")
}

func ReadFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	content := bytes.NewBuffer(data).String()
	return content, nil
}

func WriteFile(fileName string, data string) error {
	f, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func checkEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}

func checkEqualInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}

func checkEqualFloat(t *testing.T, expected, actual float64) {
	if expected != actual {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}