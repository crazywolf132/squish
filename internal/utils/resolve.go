package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"squish/internal/config"
	"strings"
)

var extensionMap = map[string][]string{
	".d.ts":  {".d.ts", ".d.mts", ".d.cts", ".ts", ".mts", ".cts"},
	".d.mts": {".d.mts", ".d.ts", ".d.cts", ".ts", ".mts", ".cts"},
	".d.cts": {".d.cts", ".d.ts", ".d.mts", ".ts", ".mts", ".cts"},
	".js":    {".js", ".ts", ".tsx", ".mts", ".cts"},
	".mjs":   {".mjs", ".js", ".cjs", ".mts", ".cts", ".ts"},
	".cjs":   {".cjs", ".js", ".mjs", ".mts", ".cts", ".ts"},
}

type SourcePathResult struct {
	Input         string `json:"input"`
	SrcExtension  string `json:"srcExtension"`
	DistExtension string `json:"distExtension"`
}

func GetSourcePath(exportEntry config.ExportEntry, source, dist string) (*SourcePathResult, error) {
	outputPath := exportEntry.OutputPath

	for distExtension, sourceExts := range extensionMap {
		if strings.HasSuffix(outputPath, distExtension) {
			pathWithoutExtension := outputPath[:len(outputPath)-len(distExtension)]
			sourcePathUnresolved := filepath.Join(source, pathWithoutExtension)

			sourcePath, err := tryExtensions(sourcePathUnresolved, sourceExts)
			if err == nil {
				return &SourcePathResult{
					Input:         sourcePath.path,
					SrcExtension:  sourcePath.extension,
					DistExtension: distExtension,
				}, nil
			}
		}
	}

	outputPathJSON, _ := json.Marshal(outputPath)
	return nil, fmt.Errorf("could not find matching source file for export path %s", string(outputPathJSON))
}

type sourcePath struct {
	path      string
	extension string
}

func tryExtensions(pathWithoutExtension string, extensions []string) (*sourcePath, error) {
	for _, extension := range extensions {
		pathWithExtension := pathWithoutExtension + extension
		if FileExists(pathWithExtension) {
			return &sourcePath{
				path:      pathWithExtension,
				extension: extension,
			}, nil
		}
	}
	return nil, fmt.Errorf("no matching file found")
}
