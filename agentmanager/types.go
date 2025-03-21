package agentmanager

import (
	"github.com/hostedgraphite/hg-cli/pipeline"
)

type Agent interface {
	InstallPipeline(chan *pipeline.Pipe) (*pipeline.Pipeline, error)
	UninstallPipeline(chan *pipeline.Pipe) (*pipeline.Pipeline, error)
	UpdateApiKeyPipeline(chan *pipeline.Pipe) (*pipeline.Pipeline, error)
}
