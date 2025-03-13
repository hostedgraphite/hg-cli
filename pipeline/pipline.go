package pipeline

import (
	"fmt"
	"os/exec"
	"time"
)

type Pipe struct {
	Name     string
	Cmd      *exec.Cmd
	Output   string
	OutErr   error
	Duration time.Duration
	Executed bool
	Success  bool
}

func NewPipe(name string, cmd *exec.Cmd) *Pipe {
	return &Pipe{
		Name: name,
		Cmd:  cmd,
	}
}

func (p *Pipe) Run() (string, error) {
	startTime := time.Now()
	output, err := p.Cmd.Output()
	p.Output = string(output)
	p.Executed = true
	p.Success = err == nil
	p.OutErr = err
	p.Duration = time.Duration(time.Since(startTime).Milliseconds())
	return p.Output, p.OutErr
}

func NewPipeline(title string, pipes []*Pipe, updates chan<- *Pipe) Pipeline {
	return Pipeline{
		Name:    title,
		Pipes:   pipes,
		Running: updates,
	}
}

type Pipeline struct {
	Name      string
	Pipes     []*Pipe
	Running   chan<- *Pipe
	LastRun   *Pipe
	OutputLog []string
	ErrLog    []string

	executed  bool
	isRunning bool
	completed bool
	failed    bool
}

func (p *Pipeline) Run() error {
	if p.executed || p.failed {
		return fmt.Errorf("pipeline already executed or failed")
	}
	p.executed = true
	p.isRunning = true
	var err error
	var output string

	for _, pipe := range p.Pipes {
		p.Running <- pipe
		output, err = pipe.Run()
		p.LastRun = pipe

		p.OutputLog = append(p.OutputLog, output)
		if err != nil {
			p.failed = true
			p.ErrLog = append(p.ErrLog, err.Error())
			break
		}
	}
	p.isRunning = false
	p.completed = true
	return err
}

func (p *Pipeline) IsCompleted() bool {
	return p.completed
}

func (p *Pipeline) GetDuration() time.Duration {
	var duration time.Duration

	for _, pipe := range p.Pipes {
		duration += pipe.Duration
	}

	return duration
}
