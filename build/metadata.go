package build

import "os"

type (
	Metadata struct {
		Project Project `toml:"project"`
	}
	Project struct {
		Name   string `toml:"name"`
		Team   string `toml:"team"`
		Domain string `toml:"domain"`
	}
	ExtraMetadata struct {
		CommitMessage string
		Version       string
	}
)

func loadExtraMetadata() ExtraMetadata {
	version := os.Getenv(versionEnv)
	if version == "" {
		version = unknown
	}

	commit := os.Getenv(comitEnv)
	if commit == "" {
		commit = unknown
	}

	return ExtraMetadata{
		CommitMessage: commit,
		Version:       version,
	}
}
