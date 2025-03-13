package pipeline

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipelineRun(t *testing.T) {
	running := make(chan *Pipe)

	pipeline := Pipeline{
		Pipes: []*Pipe{
			{
				Name: "Say hello",
				Cmd:  exec.Command("echo", "hello"),
			},
			{
				Name: "Say world",
				Cmd:  exec.Command("echo", "world"),
			},
		},
		Running: running,
	}

	var err error
	t.Log("Starting pipeline")
	ctx := PipelineRunner(&pipeline)

	t.Log("Waiting for Updates")
	func(ctx context.Context) {
		for {
			select {
			case pipe := <-running:
				t.Log("Running: " + pipe.Name)
			case <-ctx.Done():
				t.Log("Context done")
				return
			}
		}
	}(ctx)

	t.Log("Pipeline Full Output")
	t.Logf("Full Pipeline Output: \n%s", strings.Join(pipeline.OutputLog, ""))

	t.Log("Pipeline Full Error")
	t.Logf("Full Pipeline Error: \n%s", strings.Join(pipeline.ErrLog, ""))

	require.NoError(t, err)
}
