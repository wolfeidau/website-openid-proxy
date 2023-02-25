module github.com/wolfeidau/website-openid-proxy

go 1.15

require (
	github.com/alecthomas/kong v0.7.1
	github.com/apex/gateway/v2 v2.0.0
	github.com/aws/aws-lambda-go v1.37.0
	github.com/aws/aws-sdk-go v1.44.209
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dghubble/sessions v0.4.0
	github.com/golang/mock v1.6.0
	github.com/labstack/echo/v4 v4.10.2
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/rs/zerolog v1.29.0
	github.com/stretchr/testify v1.8.1
	github.com/wolfeidau/echo-s3-middleware v1.2.1-0.20210114095551-db494251c0ef
	github.com/wolfeidau/lambda-go-extras v1.5.0
	github.com/wolfeidau/lambda-go-extras/middleware/raw v1.5.0
	github.com/wolfeidau/lambda-go-extras/middleware/zerolog v1.5.0
	golang.org/x/oauth2 v0.5.0
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
