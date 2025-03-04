package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func RunCommand(cmd string, args []string, updates chan<- string) error {
	command := exec.Command(cmd, args...)

	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}

	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	if err := command.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	go streamOutput(stdoutPipe, updates)
	go streamOutput(stderrPipe, updates)

	if err := command.Wait(); err != nil {
		return fmt.Errorf("error running command: %v", err)
	}

	return nil
}

func streamOutput(pipe io.ReadCloser, updates chan<- string) {
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			// Ignore error if the pipe is closed
			if strings.HasPrefix(err.Error(), "read |0: file already closed") {
				break
			}
			updates <- "error reading command output: " + err.Error()
			break
		}
		updates <- strings.TrimRight(line, "\n")
		time.Sleep(300 * time.Millisecond)
	}
}
