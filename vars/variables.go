package vars

import "fmt"

var (
	//Name of your project
	projectName string = "gotemplate"

	//Version of your software
	projectVersion string = "0.1.0"

	//EnableStaticServer in /static folder
	enableStaticServer bool = false

	//EnableTemplates system in /templtates folder
	enableTemplates bool = false

	//LogoFile leave empty if you dont want a ASCII logo
	logoFile string = ""
)

func SetProjectName(name string) {
	projectName = name
}

func SetProjectVersion(version string) {
	projectVersion = version
}

func EnableStaticServer() {
	enableStaticServer = true
}

func EnableTemplates() {
	enableTemplates = true
}

func EnableLogoFile() {
	logoFile = "logo"
}

func FormatProjectInfo() string {
	return fmt.Sprintf("%s v%s\n", projectName, projectVersion)
}

func IsStaticServerEnabled() bool {
	return enableStaticServer
}

func IsTemplatesEnabled() bool {
	return enableTemplates
}

func IsLogoFileEnabled() bool {
	return logoFile != ""
}

func GetLogoFilename() string {
	return logoFile
}
