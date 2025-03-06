package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"
)

var agents = []string{"telegraf"}

func ShowAvailableAgents() {
	fmt.Println("Available agent: ")
	for _, agent := range agents {
		fmt.Println("- " + agent)
	}
}
func ValidateAgent(agent string) bool {
	return slices.Contains(agents, agent)
}
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

func ReadFile(filePath string) (string, error) {
	var err error
	cmd := exec.Command("sudo", "cat", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error reading the file :%v", err)
	}
	return out.String(), err
}

func WriteFile(filePath, updatedContent string) error {
	writeCmd := exec.Command("sudo", "tee", filePath)
	writeCmd.Stdin = bytes.NewBufferString(updatedContent)
	var writeOut bytes.Buffer
	writeCmd.Stdout = &writeOut
	writeCmd.Stderr = &writeOut

	if err := writeCmd.Run(); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func UpdateConfigBlock(fullConfig, confBlock string, updates map[string]string) (string, error) {
	configRegex := regexp.MustCompile(confBlock)
	configBlock := configRegex.FindString(fullConfig)

	updatedBlock := configBlock

	if configBlock == "" {
		return "", fmt.Errorf("error: no matching graphite configuration found")
	}

	for regexPattern, replacement := range updates {
		re := regexp.MustCompile(regexPattern)
		updatedBlock = re.ReplaceAllString(updatedBlock, replacement)
	}

	updatedConfig := configRegex.ReplaceAllString(fullConfig, updatedBlock)

	return updatedConfig, nil
}
