package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2")
		os.Exit(1)
	}

	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println("Error reading environment directory:", err)
		os.Exit(1)
	}

	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
