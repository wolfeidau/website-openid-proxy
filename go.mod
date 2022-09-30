module github.com/wolfeidau/website-openid-proxy

go 1.15

require (
	github.com/alecthomas/kong v0.2.12
	github.com/apex/gateway/v2 v2.0.0
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.36.26
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dghubble/sessions v0.1.0
	github.com/golang/mock v1.4.4
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/labstack/echo/v4 v4.9.0
	github.com/pquerna/cachecontrol v0.0.0-20201205024021-ac21108117ac // indirect
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/wolfeidau/echo-s3-middleware v1.2.1-0.20210114095551-db494251c0ef
	github.com/wolfeidau/lambda-go-extras v1.2.1
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
