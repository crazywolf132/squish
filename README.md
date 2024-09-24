# Squish 🍊

Squish is a minimalistic package bundler for TypeScript, built with Go. It's designed to be fast, efficient, and incredibly easy to use, with zero configuration required to get started.

![Squish Logo](https://via.placeholder.com/150x150.png?text=Squish)

## Features

- 🚀 Lightning-fast bundling
- 📦 TypeScript support out of the box
- 🔧 Zero configuration needed to start
- 🎛️ Customizable when you need it
- 👀 Watch mode for development
- 🔍 Source map generation
- 🧹 Clean dist directory option
- 🔌 Plugin system for extensibility

## Why Squish?

Squish stands out from other bundlers by prioritizing simplicity and ease of use. With Squish, you can start building your TypeScript project immediately, without the need for complex configuration files or setup processes.

### Zero Configuration

Squish works out of the box with zero configuration. It automatically reads your `package.json` file to determine:

- Entry points
- Output formats
- Package type (CommonJS or ES Module)
- TypeScript configuration (using `tsconfig.json` if present)

This means you can focus on writing code, not configuring your build tool.

## Installation

To install Squish, you need to have Go installed on your system. Then, run:

```bash
go get -u github.com/crazywolf132/squish
```

## Quick Start

1. Navigate to your TypeScript project directory.
2. Ensure your `package.json` file is set up with the appropriate `main`, `module`, `types`, and/or `exports` fields.
3. Run Squish:

```bash
squish
```

That's it! Squish will automatically bundle your TypeScript files based on your `package.json` configuration.

## Usage

While Squish works without configuration, you can customize its behavior when needed:

```
squish [flags]
```

### Flags

- `--src string`: Source directory (default "./src")
- `--dist string`: Output directory (default "./dist")
- `--minify`: Minify output
- `--watch, -w`: Watch mode
- `--target stringSlice`: Environments to support (default [es2022])
- `--tsconfig string`: Custom tsconfig.json file path
- `--env stringSlice`: Compile-time environment variables (e.g., --env NODE_ENV=production)
- `--export-condition stringSlice`: Export conditions for resolving dependency export and import maps
- `--sourcemap string`: Sourcemap generation. Provide 'inline' for inline sourcemap
- `--clean-dist`: Clean dist before bundling

## Configuration

Squish is designed to work without a dedicated configuration file. Instead, it intelligently reads your project's `package.json` file to determine the entry points and output formats. It supports various package.json fields including `main`, `module`, `types`, and `exports`.

This approach allows you to manage your project configuration in one place, reducing complexity and potential conflicts.

## Plugin System

While Squish aims for simplicity, it also provides a flexible plugin system for when you need to extend its functionality. Built-in plugins include:

- Create Require Plugin
- Externalize Node Builtins Plugin
- Patch Binary Plugin
- Strip Hashbang Plugin

## Contributing

We welcome contributions to Squish! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

Squish is [MIT licensed](LICENSE).

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub issue tracker](https://github.com/yourusername/squish/issues).

## Acknowledgements

Squish is built with the awesome [esbuild](https://github.com/evanw/esbuild) under the hood. Many thanks to the esbuild team and all our contributors!