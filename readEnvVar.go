package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func readEnvVar(varName string) string {
	return os.Getenv(varName)
}

func readAndCheckEnvVar(varName string) (varVal string) {
	if varVal = readEnvVar(varName); varVal == "" {
		err := fmt.Errorf("%s environment variable not set properly", varName)
		log.Fatal(err)
	}

	return
}

func readAndCheckIntEnvVar(varName string) (varVal string) {
	varVal = readAndCheckEnvVar(varName)
	if _, err := strconv.Atoi(varVal); err != nil {
		err := fmt.Errorf("Value of %s environment variable has to be an integer", varName)
		log.Fatal(err)
	}

	return
}
