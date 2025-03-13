package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	pipeline "github.com/hostedgraphite/hg-cli/pipeline"
)

func main() {
	var daemonMode bool
	flag.BoolVar(&daemonMode, "d", false, "run as a daemon")
	flag.Parse()

	pipes := []*pipeline.Pipe{
		pipeline.NewPipe("cmd1", exec.Command("sleep", "0.2s")),
		pipeline.NewPipe("cmd2", exec.Command("sleep", "0.55s")),
		pipeline.NewPipe("cmd3", exec.Command("sleep", "0.5s")),
		pipeline.NewPipe("cmd4", exec.Command("sleep", "0.1s")),
		pipeline.NewPipe("cmd5", exec.Command("sleep", "0.3s")),
		pipeline.NewPipe("cmd6", exec.Command("sleep", "0.323s")),
		pipeline.NewPipe("cmd7", exec.Command("sleep", "0.143s")),
	}

	updates := make(chan *pipeline.Pipe)

	// Create a new pipeline
	testpipeline := pipeline.NewPipeline("Example Pipeline", pipes, updates)

	runner := pipeline.NewRunner(
		&testpipeline,
		daemonMode,
		updates,
	)

	err := runner.Run()
	if err != nil {
		fmt.Println("Error running pipeline:", err)
		os.Exit(1)
	}

}
