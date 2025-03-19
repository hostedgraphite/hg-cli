package pipeline

import (
	"context"
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

	ctx     context.Context
	postRun func(context.Context) error
}

func (p *Pipe) execPostRun() error {
	if p.OutErr != nil || p.postRun == nil {
		return p.OutErr
	}

	if p.ctx == nil {
		p.ctx = context.Background()
	}
	p.ctx = context.WithValue(p.ctx, "output", p.Output)

	err := p.postRun(p.ctx)
	return err
}

func NewPipe(name string, cmd *exec.Cmd) *Pipe {
	return &Pipe{
		Name: name,
		Cmd:  cmd,
	}
}

func (p *Pipe) Context(ctx context.Context) *Pipe {
	p.ctx = ctx
	return p
}

func (p *Pipe) PostRun(postRun func(context.Context) error) *Pipe {
	p.postRun = postRun
	return p
}

func (p *Pipe) Run() (string, error) {
	startTime := time.Now()
	output, err := p.Cmd.Output()
	p.Output = string(output)
	p.Executed = true
	p.OutErr = err
	p.OutErr = p.execPostRun()
	p.Success = p.OutErr == nil
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
	Curr      *Pipe
	LastRun   *Pipe
	OutputLog []string
	Err       error

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
		p.Curr = pipe
		output, err = pipe.Run()
		p.LastRun = pipe

		p.OutputLog = append(p.OutputLog, output)
		if err != nil {
			p.failed = true
			p.Err = err
			break
		}
	}
	p.Curr = nil
	p.isRunning = false
	p.completed = true
	return err
}

func (p *Pipeline) IsCompleted() bool {
	return p.completed
}

func (p *Pipeline) IsRunning() bool {
	return p.isRunning
}

func (p *Pipeline) Failed() bool {
	return p.failed
}

func (p *Pipeline) Success() bool {
	return p.completed && !p.failed
}

func (p *Pipeline) GetDuration() time.Duration {
	var duration time.Duration

	for _, pipe := range p.Pipes {
		duration += pipe.Duration
	}

	return duration
}
