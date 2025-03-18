package pipeline

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stdout)

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	crossMark = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).SetString("✗")
)

// Starts the Provided Pipeline in a go routine and returns a context that will be cancelled after 240 seconds
func PipelineRunner(pipeline *Pipeline) context.Context {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

	go func(ctx context.Context) {
		defer cancelCtx()
		pipeline.Run()
	}(ctx)

	return ctx
}

func NewRunner(pipeline *Pipeline, render bool, updates chan *Pipe) *Runner {
	spin := spinner.New(spinner.WithSpinner(spinner.Dot))
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#20b9f7")).PaddingRight(1)

	r := &Runner{
		Pipeline: pipeline,
		Render:   render,
		Updates:  updates,

		spinner: spin,
	}
	return r
}

type Runner struct {
	Pipeline *Pipeline
	Render   bool
	Updates  chan *Pipe

	static  bool
	spinner spinner.Model
	ctx     context.Context
}

func (r *Runner) Init() tea.Cmd {
	r.ctx = PipelineRunner(r.Pipeline)

	return tea.Batch(
		r.spinner.Tick,
		r.nextPipelineMsg,
	)
}

type PipeUpdate struct {
	update *Pipe
}

type pipelineFinished struct {
	finished bool
}

func (r *Runner) nextPipelineMsg() tea.Msg {
	return PipeUpdate{<-r.Updates}
}

func (r *Runner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds tea.Cmd

	switch msg := msg.(type) {
	case PipeUpdate:
		cmds = tea.Batch(cmds, r.nextPipelineMsg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		r.spinner, cmd = r.spinner.Update(msg)
		cmds = tea.Batch(cmds, cmd)
	}

	if !r.static && r.Pipeline.IsCompleted() {
		return r, tea.Quit
	}

	return r, cmds
}

func (r *Runner) View() string {
	s := "\n"
	if r.Pipeline.isRunning {
		s += r.spinner.View()
	}
	s += fmt.Sprintf(r.Pipeline.Name)
	s += "\n\n"

	for _, pipe := range r.Pipeline.Pipes {
		if pipe.Executed {
			if pipe.Success {
				s += checkMark.Render("") + pipe.Name + " | " + fmt.Sprintf("finished in %dms", time.Duration(pipe.Duration))
			} else {
				s += crossMark.Render("") + pipe.Name + " | " + fmt.Sprintf("failed after %dms", time.Duration(pipe.Duration))
			}
		} else if pipe == r.Pipeline.Curr {
			s += r.spinner.View() + pipe.Name
		} else {
			s += pipe.Name
		}
		s += "\n"
	}

	if r.Pipeline.failed {
		s += fmt.Sprintf("\n\nFailed '%s' on cmd '%s'\n", r.Pipeline.Name, r.Pipeline.LastRun.Name)
		s += fmt.Sprintf("Error: %s\n", r.Pipeline.Err)
	} else if r.Pipeline.completed {
		s += fmt.Sprintf("\n%s Completed\n", r.Pipeline.Name)
	}

	return s
}

func (r *Runner) Run() error {
	var opts []tea.ProgramOption

	if !r.Render {
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	}

	_, err := tea.NewProgram(r, opts...).Run()
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) RunStatic() tea.Cmd {
	cmds := r.Init()
	r.static = true
	return cmds
}
