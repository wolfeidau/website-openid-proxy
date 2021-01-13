package flags

import "github.com/alecthomas/kong"

// API api related flags passing in env variables
type API struct {
	Version          kong.VersionFlag
	AppName          string `help:"Stage the name of the service." env:"APP_NAME"`
	Stage            string `help:"Stage the software is deployed." env:"STAGE"`
	Branch           string `help:"Branch used to build software." env:"BRANCH"`
	ClientID         string `help:"The client identifier for the openid client." env:"CLIENT_ID"`
	ClientSecret     string `help:"The client secret for the openid client" env:"CLIENT_SECRET"`
	Issuer           string `help:"The openid issuer." env:"ISSUER"`
	RedirectURL      string `help:"The redirect URL used for callbacks." env:"REDIRECT_URL"`
	SessionSecretArn string `help:"The ARN of the secret used to sign sessions." env:"SESSION_SECRET_ARN"`
	WebsiteBucket    string `help:"The name of the website S3 bucket holding content to be served." env:"WEBSITE_BUCKET"`
}
