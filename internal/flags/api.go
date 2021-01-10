package flags

import "github.com/alecthomas/kong"

// API api related flags passing in env variables
type API struct {
	Version kong.VersionFlag
	AppName string `help:"Stage the name of the service." env:"APP_NAME"`
	Stage   string `help:"Stage the software is deployed." env:"STAGE"`
	Branch  string `help:"Branch used to build software." env:"BRANCH"`
}
