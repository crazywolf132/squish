package cli

import (
	"github.com/spf13/cobra"
	"os"
	"squish/internal/config"
	"squish/internal/utils"
	"squish/internal/watcher"
	"squish/pkg/esbuild"
	"strings"
	"time"
)

var (
	srcFlag          string
	distFlag         string
	minify           bool
	watchMode        bool
	target           []string
	tsconfigPath     string
	env              []string
	exportConditions []string
	sourcemap        string
	cleanDist        bool
	bundle           bool
)

var rootCmd = &cobra.Command{
	Use:   "squish",
	Short: "Squish is a minimalistic package bundler for TypeScript",
	Run:   run,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&srcFlag, "src", "./src", "Source directory")
	rootCmd.Flags().StringVar(&distFlag, "dist", "./dist", "Output directory")
	rootCmd.Flags().BoolVar(&minify, "minify", false, "Minify output")
	rootCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "Watch mode")
	rootCmd.Flags().StringSliceVar(&target, "target", []string{"es2022"}, "Environments to support")
	rootCmd.Flags().StringVar(&tsconfigPath, "tsconfig", "", "Custom tsconfig.json file path")
	rootCmd.Flags().StringSliceVar(&env, "env", []string{}, "Compile-time environment variables (e.g., --env NODE_ENV=production)")
	rootCmd.Flags().StringSliceVar(&exportConditions, "export-condition", []string{}, "Export conditions for resolving dependency export and import maps")
	rootCmd.Flags().StringVar(&sourcemap, "sourcemap", "", "Sourcemap generation. Provide 'inline' for inline sourcemap")
	rootCmd.Flags().BoolVar(&cleanDist, "clean-dist", false, "Clean dist before bundling")
	rootCmd.Flags().BoolVar(&bundle, "bundle", true, "Bundle all dependencies")
}

func run(cmd *cobra.Command, args []string) {
	startTime := time.Now()
	cwd, err := os.Getwd()
	if err != nil {
		utils.Log("Error getting current working directory:", err)
		os.Exit(1)
	}

	pkg, err := config.ReadPackageJSON(cwd)
	if err != nil {
		utils.Log("Error reading package.json:", err)
		os.Exit(1)
	}

	utils.Log("Bundling package:", pkg.Name)

	bundlerConfig := &esbuild.BundlerConfig{
		SrcDir:           srcFlag,
		DistDir:          distFlag,
		Minify:           minify,
		Target:           target,
		TsconfigPath:     tsconfigPath,
		Env:              parseEnvFlags(env),
		ExportConditions: exportConditions,
		Sourcemap:        sourcemap,
		CleanDist:        cleanDist,
		Bundle:           bundle,
	}

	bundler := esbuild.NewBundler(bundlerConfig, pkg)

	if watchMode {
		w := watcher.NewWatcher(bundler, srcFlag)
		if err := w.Watch(); err != nil {
			utils.Log("Error watching:", err)
			os.Exit(1)
		}
	} else {
		if err := bundler.Bundle(); err != nil {
			utils.Log("Error bundling:", err)
			os.Exit(1)
		}
		utils.Log("Bundle created successfully in ", time.Since(startTime))
	}
}

func parseEnvFlags(envFlags []string) map[string]string {
	envMap := make(map[string]string)
	for _, env := range envFlags {
		key, value, _ := strings.Cut(env, "=")
		envMap[key] = value
	}
	return envMap
}
