package main

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/apex/gateway/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dghubble/sessions"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/aws-openid-proxy/internal/app"
	"github.com/wolfeidau/aws-openid-proxy/internal/cookie"
	"github.com/wolfeidau/aws-openid-proxy/internal/echosessions"
	"github.com/wolfeidau/aws-openid-proxy/internal/flags"
	"github.com/wolfeidau/aws-openid-proxy/internal/secrets"
	"github.com/wolfeidau/aws-openid-proxy/internal/server"
	"github.com/wolfeidau/aws-openid-proxy/internal/state"
	s3middleware "github.com/wolfeidau/echo-s3-middleware"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

var cfg = new(flags.API)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	secretCache := secrets.NewCache(&aws.Config{})

	sessionSecret, err := secretCache.GetValue(cfg.SessionSecretArn)
	if err != nil {
		log.Fatal().Err(err).Msg("login config failed")
	}

	sessionStore := sessions.NewCookieStore([]byte(sessionSecret), nil)

	e := echo.New()

	agr := e.Group("/auth")

	login, err := server.NewAuth(&server.AuthConfig{
		Issuer:       cfg.Issuer,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		StateStore:   state.NewCookieStore(cookie.DefaultCookieConfig),
	}, sessionStore)
	if err != nil {
		log.Fatal().Err(err).Msg("login config failed")
	}

	login.RegisterRoutes(agr)

	fs := s3middleware.New(s3middleware.FilesConfig{
		SPA:     true,
		Index:   "index.html",
		Skipper: server.LoginSkipper("/auth"),
		Summary: func(ctx context.Context, data map[string]interface{}) {
			log.Ctx(ctx).Info().Fields(data).Msg("processed s3 request")
		},
		OnErr: func(ctx context.Context, err error) {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to process s3 request")
		},
	})

	sessionMiddleware := echosessions.MiddlewareWithConfig(echosessions.Config{
		Skipper: server.LoginSkipper("/auth"),
		Store:   sessionStore,
	})

	e.Use(fs.StaticBucket(cfg.WebsiteBucket), sessionMiddleware)

	gw := gateway.NewGateway(e)

	flds := lmw.FieldMap{"commit": app.Commit, "buildDate": app.BuildDate, "stage": cfg.Stage, "branch": cfg.Branch}

	ch := lmw.New(
		zlog.New(zlog.Fields(flds)), // build a logger and inject it into the context
	)

	if cfg.Stage == "dev" {
		ch.Use(raw.New(raw.Fields(flds))) // raw event logger used during development
	}

	h := ch.Then(gw)

	lambda.StartHandler(h)
}
