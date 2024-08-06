package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PackageType string
type BinField struct {
	Single string
	Multi  map[string]string
}

const (
	PackageTypeModule   PackageType = "module"
	PackageTypeCommonJS PackageType = "commonjs"
	PackageTypeTypes    PackageType = "types"
)

type ExportEntry struct {
	OutputPath   string
	Type         PackageType
	Platform     string
	IsExecutable bool
	From         string
}

type PackageJSON struct {
	Name             string                 `json:"name"`
	Version          string                 `json:"version"`
	Type             PackageType            `json:"type"`
	Main             string                 `json:"main"`
	Module           string                 `json:"module"`
	Types            string                 `json:"types"`
	Bin              BinField               `json:"bin"`
	Exports          map[string]interface{} `json:"exports"`
	Dependencies     map[string]string      `json:"dependencies"`
	PeerDependencies map[string]string      `json:"peerDependencies"`
	DevDependencies  map[string]string      `json:"devDependencies"`
}

func ReadPackageJSON(dir string) (*PackageJSON, error) {
	path := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	if pkg.Type == "" {
		pkg.Type = PackageTypeCommonJS
	}

	return &pkg, nil
}

func (p *PackageJSON) GetExportEntries() ([]ExportEntry, error) {
	entries := []ExportEntry{}

	// Handle main entry point
	if p.Main != "" {
		entries = append(entries, ExportEntry{
			OutputPath: p.Main,
			Type:       getFileType(p.Main, p.Type),
			From:       "main",
		})
	}

	// Handle module entry point
	if p.Module != "" {
		entries = append(entries, ExportEntry{
			OutputPath: p.Module,
			Type:       PackageTypeModule,
			From:       "module",
		})
	}

	// Handle types entry point
	if p.Types != "" {
		entries = append(entries, ExportEntry{
			OutputPath: p.Types,
			Type:       PackageTypeTypes,
			From:       "types",
		})
	}

	// Handle bin entries
	if p.Bin.Single != "" {
		entries = append(entries, ExportEntry{
			OutputPath:   p.Bin.Single,
			Type:         getFileType(p.Bin.Single, p.Type),
			IsExecutable: true,
			From:         "bin",
		})
	} else if len(p.Bin.Multi) > 0 {
		for binNames, binPath := range p.Bin.Multi {
			entries = append(entries, ExportEntry{
				OutputPath:   binPath,
				Type:         getFileType(binPath, p.Type),
				IsExecutable: true,
				From:         fmt.Sprintf("bin.%s", binNames),
			})
		}
	}

	// Handle exports
	if err := p.parseExports(p.Exports, &entries, "exports"); err != nil {
		return nil, err
	}

	return entries, nil
}

func (p *PackageJSON) parseExports(exports interface{}, entries *[]ExportEntry, from string) error {
	switch e := exports.(type) {
	case string:
		if strings.HasPrefix(e, "./") {
			*entries = append(*entries, ExportEntry{
				OutputPath: e,
				Type:       getFileType(e, p.Type),
				From:       from,
			})
		}
	case map[string]interface{}:
		for key, value := range e {
			newFrom := fmt.Sprintf("%s.%s", from, key)
			switch key {
			case "require":
				if path, ok := value.(string); ok {
					*entries = append(*entries, ExportEntry{
						OutputPath: path,
						Type:       PackageTypeCommonJS,
						From:       newFrom,
					})
				}
			case "import":
				if path, ok := value.(string); ok {
					*entries = append(*entries, ExportEntry{
						OutputPath: path,
						Type:       PackageTypeModule,
						From:       newFrom,
					})
				}
			case "types":
				if path, ok := value.(string); ok {
					*entries = append(*entries, ExportEntry{
						OutputPath: path,
						Type:       PackageTypeTypes,
						From:       newFrom,
					})
				}
			case "node":
				if path, ok := value.(string); ok {
					*entries = append(*entries, ExportEntry{
						OutputPath: path,
						Type:       getFileType(path, p.Type),
						Platform:   "node",
						From:       newFrom,
					})
				}
			case "default":
				if path, ok := value.(string); ok {
					*entries = append(*entries, ExportEntry{
						OutputPath: path,
						Type:       getFileType(path, p.Type),
						From:       newFrom,
					})
				}
			default:
				if err := p.parseExports(value, entries, newFrom); err != nil {
					return err
				}
			}
		}
	case []interface{}:
		for i, value := range e {
			newFrom := fmt.Sprintf("%s[%d]", from, i)
			if err := p.parseExports(value, entries, newFrom); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported export type for %s", from)
	}
	return nil
}

func (b *BinField) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		b.Single = s
		return nil
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err == nil {
		b.Multi = m
		return nil
	}

	return fmt.Errorf("bin field must be either a string or a map[string]string")
}

func getFileType(filePath string, defaultType PackageType) PackageType {
	switch {
	case strings.HasSuffix(filePath, ".mjs"):
		return PackageTypeModule
	case strings.HasSuffix(filePath, ".cjs"):
		return PackageTypeCommonJS
	case strings.HasSuffix(filePath, ".d.ts"):
		return PackageTypeTypes
	default:
		return defaultType
	}
}
