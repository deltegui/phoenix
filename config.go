package locomotive

import "github.com/deltegui/locomotive/vars"

type LocomotiveConfig struct {}

func Configure() LocomotiveConfig {
	return LocomotiveConfig{}
}

func(config LocomotiveConfig) SetProjectInfo(name string, version string) LocomotiveConfig {
	vars.SetProjectName(name)
	vars.SetProjectVersion(version)
	return config
}

func(config LocomotiveConfig) EnableStaticServer() LocomotiveConfig {
	vars.EnableStaticServer()
	return config
}

func(config LocomotiveConfig) EnableTemplates() LocomotiveConfig {
	vars.EnableTemplates()
	return config
}

func(config LocomotiveConfig) EnableLogoFile() LocomotiveConfig {
	vars.EnableLogoFile()
	return config
}