package esbuild

import (
	"github.com/evanw/esbuild/pkg/api"
)

func CreateRequirePlugin() Plugin {
	const virtualModuleName = "pkgroll:create-require"
	const isEsmVariableName = "IS_ESM"

	return NewPluginBuilder("create-require").
		Setup(func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: "^" + virtualModuleName + "$"}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      virtualModuleName,
					Namespace: "create-require",
				}, nil
			})

			build.OnLoad(api.OnLoadOptions{Filter: "^" + virtualModuleName + "$", Namespace: "create-require"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				contents := `
					import { createRequire } from 'module';
					export default (
						` + isEsmVariableName + `
							? /* @__PURE__ */ createRequire(import.meta.url)
							: require
					);
				`
				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderJS,
				}, nil
			})
		}).
		Build()
}

func IsFormatEsmPlugin(isEsm bool) Plugin {
	const isEsmVariableName = "IS_ESM"

	return NewPluginBuilder("is-format-esm").
		Setup(func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: ".*"}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{}, nil
			})

			build.OnLoad(api.OnLoadOptions{Filter: ".*"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				contents := "const " + isEsmVariableName + " = " + boolToString(isEsm) + ";"
				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderJS,
				}, nil
			})
		}).
		Build()
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
