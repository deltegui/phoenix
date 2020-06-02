package phoenix

import "fmt"

type PhoenixConfig struct {
	// name of your project
	projectName string

	// Version of your software
	projectVersion string

	// enable static server in /static folder
	enableStaticServer bool

	// logo file leave empty if you dont want an ASCII logo
	logoFile string
}

func (config *PhoenixConfig) SetProjectInfo(name string, version string) *PhoenixConfig {
	config.projectName = name
	config.projectVersion = version
	return config
}

func (config *PhoenixConfig) EnableStaticServer() *PhoenixConfig {
	config.enableStaticServer = true
	return config
}

func (config *PhoenixConfig) EnableLogoFile() *PhoenixConfig {
	config.logoFile = "logo"
	return config
}

func (config PhoenixConfig) formatProjectInfo() string {
	return fmt.Sprintf("%s v%s\n", config.projectName, config.projectVersion)
}

func (config PhoenixConfig) isStaticServerEnabled() bool {
	return config.enableStaticServer
}

func (config PhoenixConfig) isLogoFileEnabled() bool {
	return config.logoFile != ""
}

func (config PhoenixConfig) getLogoFilename() string {
	return config.logoFile
}
