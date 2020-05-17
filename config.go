package phoenix

import "github.com/deltegui/phoenix/vars"

type PhoenixConfig struct{}

func Configure() PhoenixConfig {
	return PhoenixConfig{}
}

func (config PhoenixConfig) SetProjectInfo(name string, version string) PhoenixConfig {
	vars.SetProjectName(name)
	vars.SetProjectVersion(version)
	return config
}

func (config PhoenixConfig) EnableStaticServer() PhoenixConfig {
	vars.EnableStaticServer()
	return config
}

func (config PhoenixConfig) EnableTemplates() PhoenixConfig {
	vars.EnableTemplates()
	return config
}

func (config PhoenixConfig) EnableLogoFile() PhoenixConfig {
	vars.EnableLogoFile()
	return config
}
