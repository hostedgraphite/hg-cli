package pipeline

import (
	"context"
	"fmt"
	"io"
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

func NewRunner(pipeline *Pipeline, daemonize bool, updates chan *Pipe) *runner {
	spin := spinner.New()
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).PaddingRight(1)

	r := &runner{
		Pipeline:  pipeline,
		Daemonize: daemonize,
		Updates:   updates,

		spinner: spin,
	}
	return r
}

type runner struct {
	Pipeline  *Pipeline
	Daemonize bool
	Updates   chan *Pipe

	spinner spinner.Model
	ctx     context.Context
	logger  *log.Logger
}

func (r *runner) Init() tea.Cmd {
	logger.Infof("Running Pipeline: %s", r.Pipeline.Name)
	r.ctx = PipelineRunner(r.Pipeline)

	return tea.Batch(
		r.spinner.Tick,
		r.nextPipelineMsg,
	)
}

type pipelineUpdate struct {
	update *Pipe
}

type pipelineFinished struct {
	finished bool
}

func (r *runner) nextPipelineMsg() tea.Msg {
	return pipelineUpdate{<-r.Updates}
}

func (r *runner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
	case pipelineUpdate:
		output := msg.update.Name
		logger.Info(output)
		cmds = tea.Batch(cmds, r.nextPipelineMsg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		r.spinner, cmd = r.spinner.Update(msg)
		cmds = tea.Batch(cmds, cmd)
	}

	if r.Pipeline.IsCompleted() {
		return r, tea.Quit
	}

	return r, cmds
}

func (r *runner) View() string {
	s := "\n" + r.spinner.View() + fmt.Sprintf(r.Pipeline.Name)
	s += "\n\n"

	for _, pipe := range r.Pipeline.Pipes {
		if pipe.Executed {
			if pipe.Success {
				s += checkMark.Render("") + pipe.Name + fmt.Sprintf(" | finished in %dms", time.Duration(pipe.Duration))
			} else {
				s += crossMark.Render("") + pipe.Name + fmt.Sprintf(" | failed after %dms", time.Duration(pipe.Duration))
			}
		} else if pipe == r.Pipeline.Curr {
			s += r.spinner.View() + pipe.Name
		}
		s += "\n"
	}

	return s
}

func (r *runner) Run() error {
	var opts []tea.ProgramOption

	if r.Daemonize {
		opts = []tea.ProgramOption{tea.WithoutRenderer()}
	} else {
		logger.SetOutput(io.Discard)
	}

	_, err := tea.NewProgram(r, opts...).Run()
	if err != nil {
		fmt.Printf("Runner Failed")
		return err
	}

	if r.Pipeline.failed {
		fmt.Printf("\n\nFailed '%s' on cmd '%s'\n", r.Pipeline.Name, r.Pipeline.LastRun.Name)
		fmt.Printf("Error: %s\n", r.Pipeline.Err)
		return r.Pipeline.Err
	}

	fmt.Printf("\n%s Completed\n", r.Pipeline.Name)
	return err
}
