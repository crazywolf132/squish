package esbuild

import (
	"github.com/evanw/esbuild/pkg/api"
)

type PluginHooks struct {
	Name      string
	Setup     func(build api.PluginBuild)
	OnStart   func() (api.OnStartResult, error)
	OnEnd     func(result *api.BuildResult) (api.OnEndResult, error)
	OnResolve func(args api.OnResolveArgs) (api.OnResolveResult, error)
	OnLoad    func(args api.OnLoadArgs) (api.OnLoadResult, error)
}

type Plugin interface {
	Hooks() PluginHooks
}

type PluginBuilder struct {
	hooks PluginHooks
}

func NewPluginBuilder(name string) *PluginBuilder {
	return &PluginBuilder{
		hooks: PluginHooks{Name: name},
	}
}

func (pb *PluginBuilder) Setup(fn func(build api.PluginBuild)) *PluginBuilder {
	pb.hooks.Setup = fn
	return pb
}

func (pb *PluginBuilder) OnStart(fn func() (api.OnStartResult, error)) *PluginBuilder {
	pb.hooks.OnStart = fn
	return pb
}

func (pb *PluginBuilder) OnEnd(fn func(result *api.BuildResult) (api.OnEndResult, error)) *PluginBuilder {
	pb.hooks.OnEnd = fn
	return pb
}

func (pb *PluginBuilder) OnResolve(fn func(args api.OnResolveArgs) (api.OnResolveResult, error)) *PluginBuilder {
	pb.hooks.OnResolve = fn
	return pb
}

func (pb *PluginBuilder) OnLoad(fn func(args api.OnLoadArgs) (api.OnLoadResult, error)) *PluginBuilder {
	pb.hooks.OnLoad = fn
	return pb
}

func (pb *PluginBuilder) Build() Plugin {
	return &pluginImpl{hooks: pb.hooks}
}

type pluginImpl struct {
	hooks PluginHooks
}

func (p *pluginImpl) Hooks() PluginHooks {
	return p.hooks
}
