package main

import (
	"fmt"
	"os/exec"
	"testing"
)

// Test function which checks build process
func TestAwsjack(t *testing.T) {
	cmd := exec.Command("go", "build", "main.go")

	output, err := cmd.Output()
	// Check the output for errors
	if err != nil {
		t.Error("Failure Occcured=", err)
	} else {
		outputString := fmt.Sprintf("succes = %s", output)
		fmt.Println("Success=", outputString)

	}

}
