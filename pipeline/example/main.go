package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	pipeline "github.com/hostedgraphite/hg-cli/pipeline"
)

func fileWriter(ctx context.Context) error {
	output := ctx.Value("output").(string)
	path := ctx.Value("path").(string)
	err := os.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileDeleter(ctx context.Context) error {
	path := ctx.Value("path").(string)
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	fileCtx := context.WithValue(context.Background(), "path", "./text.txt")

	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("cmd1", exec.Command("sleep", "0.4s")),
		pipeline.NewPipe("cmd2", exec.Command("sleep", "0.65s")),
		pipeline.NewPipe("cmd3", exec.Command("sleep", "1.5s")),
		pipeline.NewPipe("cmd4", exec.Command("sleep", "0.6s")),
		pipeline.NewPipe("cmd5", exec.Command("echo", "hello")).Context(fileCtx).PostRun(fileWriter),
		pipeline.NewPipe("cmd6", exec.Command("sleep", "0.323s")),
		pipeline.NewPipe("cmd7", exec.Command("sleep", "1.143s")).Context(fileCtx).PostRun(fileDeleter),
	}

	updates := make(chan *pipeline.Pipe)

	// Create a new pipeline
	testpipeline := pipeline.NewPipeline("Example Pipeline", pipes, updates)

	runner := pipeline.NewRunner(
		&testpipeline,
		true,
		updates,
	)

	err := runner.Run()
	if err != nil {
		fmt.Println("Error running pipeline:", err)
		os.Exit(1)
	}

}
