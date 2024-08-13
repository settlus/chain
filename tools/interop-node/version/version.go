// Package version is a convenience utility that provides interop-node
// consumers with a ready-to-use version command that
// produces apps versioning information based on flags
// passed at compile time.
//
// # Configure the version command
//
// The version command can be just added to your cobra root command.
// At build time, the variables Name, Version, Commit, and BuildTags
// can be passed as build flags as shown in the following example:
//
//	go build -X github.com/settlus/chain/tools/interop-node/version.Name=settlus \
//	 -X github.com/settlus/chain/tools/interop-node/version.AppName=settlusd \
//	 -X github.com/settlus/chain/tools/interop-node/version.Version=1.0 \
//	 -X github.com/settlus/chain/tools/interop-node/version.Commit=f0f7b7dab7e36c20b757cebce0e8f4fc5b95de60 \
//	 -X "github.com/settlus/chain/tools/interop-node/version.BuildTags=linux darwin amd64"
package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
)

// ContextKey is used to store the ExtraInfo in the context.
type ContextKey struct{}

var (
	// Name application's name
	Name = ""
	// AppName application binary name
	AppName = "<appd>"
	// Version application's version string
	Version = ""
	// Commit commit
	Commit = ""
	// BuildTags build tags
	BuildTags = ""
)

// ExtraInfo contains a set of extra information provided by apps
type ExtraInfo map[string]string

// Info defines the application version information.
type Info struct {
	Name      string     `json:"name" yaml:"name"`
	AppName   string     `json:"server_name" yaml:"server_name"`
	Version   string     `json:"version" yaml:"version"`
	GitCommit string     `json:"commit" yaml:"commit"`
	BuildTags string     `json:"build_tags" yaml:"build_tags"`
	GoVersion string     `json:"go" yaml:"go"`
	BuildDeps []buildDep `json:"build_deps" yaml:"build_deps"`
	ExtraInfo ExtraInfo  `json:"extra_info,omitempty" yaml:"extra_info,omitempty"`
}

func NewInfo() Info {
	return Info{
		Name:      Name,
		AppName:   AppName,
		Version:   Version,
		GitCommit: Commit,
		BuildTags: BuildTags,
		GoVersion: fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		BuildDeps: depsFromBuildInfo(),
	}
}

func (vi Info) String() string {
	return fmt.Sprintf(`%s: %s
git commit: %s
build tags: %s
%s`,
		vi.Name, vi.Version, vi.GitCommit, vi.BuildTags, vi.GoVersion,
	)
}

func depsFromBuildInfo() (deps []buildDep) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	for _, dep := range buildInfo.Deps {
		deps = append(deps, buildDep{dep})
	}

	return
}

type buildDep struct {
	*debug.Module
}

func (d buildDep) String() string {
	if d.Replace != nil {
		return fmt.Sprintf("%s@%s => %s@%s", d.Path, d.Version, d.Replace.Path, d.Replace.Version)
	}

	return fmt.Sprintf("%s@%s", d.Path, d.Version)
}

func (d buildDep) MarshalJSON() ([]byte, error)      { return json.Marshal(d.String()) }
func (d buildDep) MarshalYAML() (interface{}, error) { return d.String(), nil }
