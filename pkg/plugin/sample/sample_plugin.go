package sample

import (
	"context"
	"encoding/json"

	"github.com/mfojtik/fob/pkg/artifact"
	"github.com/mfojtik/fob/pkg/plugin"
)

type Metadata struct {
	InfraCommit   string            `json:"infra-commit"`
	JobVersion    string            `json:"job-version"`
	Pod           string            `json:"pod"`
	Repo          string            `json:"repo"`
	RepoCommit    string            `json:"repo-commit"`
	Repos         map[string]string `json:"repos"`
	WorkNamespace string            `json:"work-namespace"`
}

type SamplePlugin struct{}

func (s *SamplePlugin) Register(manager plugin.PluginManager) error {
	manager.Add(s)
	return nil
}

func (s *SamplePlugin) Execute(ctx context.Context, o plugin.PluginOptions) error {
	metadataBytes, err := artifact.Get(ctx, o.JobUrl, "/artifacts/metadata.json")
	if err != nil {
		return err
	}

	metadata := Metadata{}
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return err
	}

	o.Printf("Job %q running in %q namespace", metadata.Repo, metadata.WorkNamespace)

	return nil
}
