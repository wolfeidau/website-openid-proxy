package main

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/apex/gateway/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	s3middleware "github.com/wolfeidau/echo-s3-middleware"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
	"github.com/wolfeidau/website-openid-proxy/internal/app"
	"github.com/wolfeidau/website-openid-proxy/internal/flags"
	"github.com/wolfeidau/website-openid-proxy/internal/secrets"
	"github.com/wolfeidau/website-openid-proxy/internal/server"
	"github.com/wolfeidau/website-openid-proxy/internal/session"
)

var cfg = new(flags.API)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	if err := cfg.Valid(); err != nil {
		log.Fatal().Err(err).Msg("config validation failed")
	}

	e := echo.New()

	secretCache := secrets.NewCache(&aws.Config{})

	// session middleware is available everwhere
	sessionMiddleware, err := session.SetupMiddleware(cfg, secretCache)
	if err != nil {
		log.Fatal().Err(err).Msg("session middleware setup failed")
	}

	e.Use(sessionMiddleware)

	agr := e.Group("/auth")

	login, err := server.NewAuth(cfg, oidc.NewProvider)
	if err != nil {
		log.Fatal().Err(err).Msg("auth config failed")
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

	e.Use(server.CheckAuthWithConfig(server.Config{
		Skipper: server.LoginSkipper("/auth"),
	}))

	e.Use(fs.StaticBucket(cfg.WebsiteBucket))

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
