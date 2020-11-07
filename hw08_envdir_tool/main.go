package main

import (
	"fmt"
	"os"
)

func main() {
	envDir := os.Args[1]
	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Printf("Couldn't read environment variables from dir '%v'\n", envDir)
		os.Exit(111)
	}
	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}
