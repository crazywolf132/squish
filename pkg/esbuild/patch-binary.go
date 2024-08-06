package esbuild

import (
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"path/filepath"
)

func PatchBinaryPlugin(executablePaths []string) Plugin {
	return NewPluginBuilder("patch-binary").
		Setup(func(build api.PluginBuild) {
			build.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
				for _, outputFile := range result.OutputFiles {
					if isExecutable(outputFile.Path, executablePaths) {
						content := "#!/usr/bin/env node\n" + string(outputFile.Contents)
						err := os.WriteFile(outputFile.Path, []byte(content), 0755)
						if err != nil {
							return api.OnEndResult{}, err
						}
					}
				}
				return api.OnEndResult{}, nil
			})
		}).
		Build()
}

func isExecutable(path string, executablePaths []string) bool {
	for _, execPath := range executablePaths {
		if filepath.Base(path) == filepath.Base(execPath) {
			return true
		}
	}
	return false
}
