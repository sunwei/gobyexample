package main

import "fmt"

type Mount struct {
	// relative path in source repo, e.g. "scss"
	Source string
	// relative target path, e.g. "assets/bootstrap/scss"
	Target string
	// any language code associated with this mount.
	Lang string
}

type Import struct {
	// Module path
	Path string
}

// Config holds a module config.
type Config struct {
	Mounts  []Mount
	Imports []Import
}

type Module interface {
	// Config The decoded module config and mounts.
	Config() Config
	// Owner In the dependency tree, this is the first module that defines this module
	// as a dependency.
	Owner() Module
	// Mounts Any directory remappings.
	Mounts() []Mount
}

type Modules []Module

var modules Modules

// moduleAdapter implemented Module interface
type moduleAdapter struct {
	projectMod bool
	owner      Module
	mounts     []Mount
	config     Config
}

func (m *moduleAdapter) Config() Config {
	return m.config
}
func (m *moduleAdapter) Mounts() []Mount {
	return m.mounts
}
func (m *moduleAdapter) Owner() Module {
	return m.owner
}

// happy path to easily understand
func main() {
	// project module config
	moduleConfig := Config{}

	imports := []string{"mytheme"}
	for _, imp := range imports {
		moduleConfig.Imports = append(
			moduleConfig.Imports, Import{
				Path: imp,
			})
	}

	// Need to run these after the modules are loaded, but before
	// they are finalized.
	collectHook := func(mods Modules) {
		// Apply default project mounts.
		// Default folder structure for hugo project
		ApplyProjectConfigDefaults(mods[0])
	}

	collectModules(moduleConfig, collectHook)

	for _, m := range modules {
		fmt.Printf("%#v\n", m)
	}
}

// Module folder structure
const (
	ComponentFolderArchetypes = "archetypes"
	ComponentFolderStatic     = "static"
	ComponentFolderLayouts    = "layouts"
	ComponentFolderContent    = "content"
	ComponentFolderData       = "data"
	ComponentFolderAssets     = "assets"
	ComponentFolderI18n       = "i18n"
)

// ApplyProjectConfigDefaults applies default/missing module configuration for
// the main project.
func ApplyProjectConfigDefaults(mod Module) {
	projectMod := mod.(*moduleAdapter)

	type dirKeyComponent struct {
		key          string
		component    string
		multilingual bool
	}

	dirKeys := []dirKeyComponent{
		{"contentDir", ComponentFolderContent, true},
		{"dataDir", ComponentFolderData, false},
		{"layoutDir", ComponentFolderLayouts, false},
		{"i18nDir", ComponentFolderI18n, false},
		{"archetypeDir", ComponentFolderArchetypes,
			false},
		{"assetDir", ComponentFolderAssets, false},
		{"", ComponentFolderStatic, false},
	}

	var mounts []Mount
	for _, d := range dirKeys {
		if d.multilingual {
			// based on language content configuration
			// multiple language has multiple source folders
			if d.component == ComponentFolderContent {
				mounts = append(mounts, Mount{
					Lang:   "en",
					Source: "mycontent",
					Target: d.component})
			}
		} else {
			mounts = append(mounts,
				Mount{
					Source: d.component,
					Target: d.component})
		}
	}

	projectMod.mounts = mounts
}

func collectModules(modConfig Config,
	hookBeforeFinalize func(m Modules)) {
	projectMod := &moduleAdapter{
		projectMod: true,
		config:     modConfig,
	}

	// module structure, [project, others...]
	addAndRecurse(projectMod)

	// Add the project mod on top.
	modules = append(Modules{projectMod}, modules...)

	if hookBeforeFinalize != nil {
		hookBeforeFinalize(modules)
	}
}

// addAndRecurse Project Imports -> Import imports
func addAndRecurse(owner *moduleAdapter) {
	moduleConfig := owner.Config()

	// theme may depend on other theme
	for _, moduleImport := range moduleConfig.Imports {
		tc := add(owner, moduleImport)
		if tc == nil {
			continue
		}
		// tc is mytheme with no config file
		addAndRecurse(tc)
	}
}

func add(owner *moduleAdapter,
	moduleImport Import) *moduleAdapter {

	fmt.Printf("start to create `%s` module\n",
		moduleImport.Path)
	ma := &moduleAdapter{
		owner: owner,
		// in the example, mytheme has no other import
		config: Config{},
	}
	modules = append(modules, ma)
	return ma
}
