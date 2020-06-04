package phoenix

import "fmt"

type PhoenixConfig struct {
	// name of your project
	projectName string

	// Version of your software
	projectVersion string

	// enable static server in /static folder
	enableStaticServer bool

	// enables session cookie
	enableSessions bool

	// logo file leave empty if you dont want an ASCII logo
	logoFile string

	// function to be called when server is shutting down
	onStop func()

	// enable HTTPS (TLS) throught pem files
	tlsCertFile string
	tlsKeyFile  string

	// enable HTTPS (TLS) using let's encrypt with domains
	domains []string
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

func (config *PhoenixConfig) EnableSessions() *PhoenixConfig {
	config.enableSessions = true
	return config
}

func (config *PhoenixConfig) EnableLogoFile() *PhoenixConfig {
	config.logoFile = "logo"
	return config
}

func (config *PhoenixConfig) SetStopHandler(onStop func()) *PhoenixConfig {
	config.onStop = onStop
	return config
}

func (config *PhoenixConfig) UseHTTPS(certFile, keyFile string) *PhoenixConfig {
	config.tlsCertFile = certFile
	config.tlsKeyFile = keyFile
	return config
}

func (config *PhoenixConfig) UseAutoHTTPS(domains ...string) *PhoenixConfig {
	config.domains = domains
	return config
}

func (config PhoenixConfig) formatProjectInfo() string {
	return fmt.Sprintf("%s v%s\n", config.projectName, config.projectVersion)
}

func (config PhoenixConfig) isStaticServerEnabled() bool {
	return config.enableStaticServer
}

func (config PhoenixConfig) areSessionsEnabled() bool {
	return config.enableSessions
}

func (config PhoenixConfig) isLogoFileEnabled() bool {
	return config.logoFile != ""
}

func (config PhoenixConfig) isHTTPSEnabled() bool {
	return config.tlsCertFile != "" && config.tlsKeyFile != ""
}

func (config PhoenixConfig) getHTTPSCertKeyFiles() (string, string) {
	return config.tlsCertFile, config.tlsKeyFile
}

func (config PhoenixConfig) isAutoHTTPSEnabled() bool {
	return config.domains != nil
}

func (config PhoenixConfig) getAutoHTTPSDomains() []string {
	return config.domains
}

func (config PhoenixConfig) getLogoFilename() string {
	return config.logoFile
}

func (config PhoenixConfig) getStopHandler() func() {
	return config.onStop
}
