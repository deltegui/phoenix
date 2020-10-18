package phoenix

import (
	"fmt"
	"net/http"
)

// PhoenixConfig stores all project configurations.
// It's a struct instead of global variables because
// it's stored in app instance.
type PhoenixConfig struct {
	// name of your project
	projectName string

	// Version of your software
	projectVersion string

	// enable static server in /static folder
	enableStaticServer bool

	// logo file leave empty if you dont want an ASCII logo
	logoFile string

	// function to be called when server is shutting down
	onStop func()

	// function to be called when server should start. Use it to customize
	// calling ListenAndServer like function
	onStart func(*http.Server) error
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

func (config *PhoenixConfig) StopHook(onStop func()) *PhoenixConfig {
	config.onStop = onStop
	return config
}

func (config *PhoenixConfig) StartHook(onStart func(*http.Server) error) *PhoenixConfig {
	config.onStart = onStart
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
