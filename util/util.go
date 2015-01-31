package util

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
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
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
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
