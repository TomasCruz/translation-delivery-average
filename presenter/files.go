package presenter

import (
	"bufio"
	"fmt"
	"os"
)

// FileToStrings reads the entire file and returns its lines in slice of strings
func FileToStrings(fileName string) (textLines []string, err error) {
	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		textLines = append(textLines, scanner.Text())
	}

	return
}
