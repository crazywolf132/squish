package esbuild

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"path/filepath"
	"squish/internal/config"
	"squish/internal/utils"
)

type BundlerConfig struct {
	SrcDir           string
	DistDir          string
	Minify           bool
	Target           []string
	TsconfigPath     string
	Env              map[string]string
	ExportConditions []string
	Sourcemap        string
	CleanDist        bool
}

type Bundler struct {
	config   *BundlerConfig
	pkg      *config.PackageJSON
	buildCtx *api.BuildContext
}

func NewBundler(config *BundlerConfig, pkg *config.PackageJSON) *Bundler {
	return &Bundler{
		config: config,
		pkg:    pkg,
	}
}

func (b *Bundler) Bundle() error {
	if b.config.CleanDist {
		if err := utils.CleanDirectory(b.config.DistDir); err != nil {
			return fmt.Errorf("failed to clean dist directory: %w", err)
		}
	}

	entries, err := b.pkg.GetExportEntries()
	if err != nil {
		return err
	}

	// Generate TypeScript declaration files
	//if err := utils.RunTSC(b.config.SrcDir, b.config.DistDir, b.config.TsconfigPath); err != nil {
	//	return err
	//}

	for _, entry := range entries {
		sourcePath, err := utils.GetSourcePath(entry, b.config.SrcDir, b.config.DistDir)
		if err != nil {
			return fmt.Errorf("error resolving source path: %w", err)
		}

		if err := b.bundleEntry(sourcePath, entry); err != nil {
			return fmt.Errorf("failed to bundle entry: %s, %w", entry.OutputPath, err)
		}
	}

	return nil
}

func (b *Bundler) bundleEntry(sourcePath *utils.SourcePathResult, entry config.ExportEntry) error {
	outfile := filepath.Join(b.config.DistDir, entry.OutputPath)

	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outfile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	//isEsm := b.getFormat(entry.Type) == api.FormatESModule

	plugins := []api.Plugin{
		createEsbuildPlugin(CreateRequirePlugin()),
		//createEsbuildPlugin(IsFormatEsmPlugin(isEsm)),
		createEsbuildPlugin(ExternalizeNodeBuiltinsPlugin(b.config.Target)),
		//createEsbuildPlugin(StripHashbangPlugin()),
	}

	if entry.IsExecutable {
		plugins = append(plugins, createEsbuildPlugin(PatchBinaryPlugin([]string{entry.OutputPath})))
	}

	buildOptions := api.BuildOptions{
		EntryPoints:       []string{sourcePath.Input},
		Outfile:           outfile,
		Bundle:            true,
		Write:             true,
		Format:            b.getFormat(entry.Type),
		Target:            b.getEsbuildTarget(),
		Platform:          api.PlatformNode,
		External:          b.getExternalDependencies(),
		Define:            b.getDefine(),
		Sourcemap:         b.getSourcemap(),
		MinifyWhitespace:  b.config.Minify,
		MinifyIdentifiers: b.config.Minify,
		MinifySyntax:      b.config.Minify,
		Plugins:           plugins,
		TreeShaking:       api.TreeShakingTrue,
	}

	ctx, err := api.Context(buildOptions)
	if err != nil {
		return err
	}

	if b.config.TsconfigPath != "" {
		buildOptions.Tsconfig = b.config.TsconfigPath
	}

	if len(b.config.ExportConditions) > 0 {
		buildOptions.Conditions = b.config.ExportConditions
	}

	var result api.BuildResult

	if b.buildCtx != nil {
		result = ctx.Rebuild()
	} else {
		result = api.Build(buildOptions)
		b.buildCtx = &ctx
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("build failed for %s: %v", entry.OutputPath, result.Errors)
	}

	return nil
}

func (b *Bundler) getFormat(packageType config.PackageType) api.Format {
	switch packageType {
	case config.PackageTypeModule:
		return api.FormatESModule
	case config.PackageTypeCommonJS:
		return api.FormatCommonJS
	default:
		return api.FormatESModule
	}
}

func (b *Bundler) getEsbuildTarget() api.Target {
	targets := make([]api.Target, 0, len(b.config.Target))
	for _, t := range b.config.Target {
		switch t {
		case "es2022":
			targets = append(targets, api.ES2022)
		case "es2021":
			targets = append(targets, api.ES2021)
		case "es2020":
			targets = append(targets, api.ES2020)
		// Add more cases as needed
		default:
			targets = append(targets, api.ES2022) // Default to ES2022
		}
	}
	return targets[0] // esbuild only accepts a single target, so we use the first one
}

func (b *Bundler) getSourcemap() api.SourceMap {
	switch b.config.Sourcemap {
	case "inline":
		return api.SourceMapInline
	case "":
		return api.SourceMapNone
	default:
		return api.SourceMapLinked
	}
}

func (b *Bundler) getDefine() map[string]string {
	define := make(map[string]string)
	for key, value := range b.config.Env {
		define[fmt.Sprintf("process.env.%s", key)] = fmt.Sprintf("\"%s\"", value)
	}
	return define
}

func (b *Bundler) getExternalDependencies() []string {
	externals := make([]string, 0)
	for dep := range b.pkg.Dependencies {
		externals = append(externals, dep)
	}
	for dep := range b.pkg.PeerDependencies {
		externals = append(externals, dep)
	}
	return externals
}

func createEsbuildPlugin(p Plugin) api.Plugin {
	return api.Plugin{
		Name: p.Hooks().Name,
		Setup: func(build api.PluginBuild) {
			if p.Hooks().Setup != nil {
				p.Hooks().Setup(build)
			}
			if p.Hooks().OnStart != nil {
				build.OnStart(p.Hooks().OnStart)
			}
			if p.Hooks().OnEnd != nil {
				build.OnEnd(p.Hooks().OnEnd)
			}
			if p.Hooks().OnResolve != nil {
				build.OnResolve(api.OnResolveOptions{Filter: ".*"}, p.Hooks().OnResolve)
			}
			if p.Hooks().OnLoad != nil {
				build.OnLoad(api.OnLoadOptions{Filter: ".*"}, p.Hooks().OnLoad)
			}
		},
	}
}
