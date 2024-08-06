package esbuild

import (
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"regexp"
)

var hashbangPattern = regexp.MustCompile(`^#!.*`)

func StripHashbangPlugin() Plugin {
	return NewPluginBuilder("strip-hashbang").
		Setup(func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: ".*"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				contents, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				strippedContents := string(hashbangPattern.ReplaceAll(contents, []byte{}))

				return api.OnLoadResult{
					Contents: &strippedContents,
					Loader:   api.LoaderJS,
				}, nil
			})
		}).
		Build()
}
