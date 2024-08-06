package esbuild

import (
	"github.com/evanw/esbuild/pkg/api"
	"strconv"
	"strings"
)

var nodeBuiltins = []string{
	"assert", "buffer", "child_process", "cluster", "crypto", "dgram", "dns", "domain", "events", "fs", "http", "https", "net",
	"os", "path", "punycode", "querystring", "readline", "stream", "string_decoder", "tls", "tty", "url", "util", "v8", "vm", "zlib",
}

func ExternalizeNodeBuiltinsPlugin(target []string) Plugin {
	stripNodeProtocol := shouldStripNodeProtocol(target)

	return NewPluginBuilder("externalize-node-builtins").
		Setup(func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: ".*"}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				if strings.HasPrefix(args.Path, "node:") {
					if stripNodeProtocol {
						return api.OnResolveResult{
							Path:     strings.TrimPrefix(args.Path, "node:"),
							External: true,
						}, nil
					}
					return api.OnResolveResult{External: true}, nil
				}

				for _, builtin := range nodeBuiltins {
					if args.Path == builtin {
						return api.OnResolveResult{External: true}, nil
					}
				}

				return api.OnResolveResult{}, nil
			})
		}).
		Build()
}

func shouldStripNodeProtocol(target []string) bool {
	for _, t := range target {
		if strings.HasPrefix(t, "node") {
			version := strings.TrimPrefix(t, "node")
			parts := strings.Split(version, ".")
			major, _ := strconv.Atoi(parts[0])
			minor, _ := strconv.Atoi(parts[1])

			if (major == 12 && minor >= 20) || major >= 14 {
				return false
			}
		}
	}
	return true
}
