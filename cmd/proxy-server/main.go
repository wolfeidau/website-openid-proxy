package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	middleware "github.com/wolfeidau/echo-middleware"
	s3middleware "github.com/wolfeidau/echo-s3-middleware"
	"github.com/wolfeidau/website-openid-proxy/internal/app"
	"github.com/wolfeidau/website-openid-proxy/internal/flags"
	"github.com/wolfeidau/website-openid-proxy/internal/secrets"
	"github.com/wolfeidau/website-openid-proxy/internal/server"
	"github.com/wolfeidau/website-openid-proxy/internal/session"
)

var cfg = new(flags.ServerAPI)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	flds := map[string]interface{}{"commit": app.Commit, "buildDate": app.BuildDate, "stage": cfg.Stage, "branch": cfg.Branch}

	if err := cfg.Valid(); err != nil {
		log.Fatal().Err(err).Msg("config validation failed")
	}

	e := echo.New()

	e.Logger.SetOutput(ioutil.Discard)

	e.Use(middleware.ZeroLogWithConfig(
		middleware.ZeroLogConfig{
			Fields: flds,
		},
	))

	secretCache := secrets.NewCache(&aws.Config{})

	// session middleware is available everwhere
	sessionMiddleware, err := session.SetupMiddleware(&cfg.API, secretCache)
	if err != nil {
		log.Fatal().Err(err).Msg("session middleware setup failed")
	}

	e.Use(sessionMiddleware)

	agr := e.Group("/auth")

	login, err := server.NewAuth(&cfg.API, oidc.NewProvider)
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

	log.Info().Str("port", cfg.Port).Str("cert", cfg.CertFile).Msg("listing")

	log.Error().Err(e.StartTLS(fmt.Sprintf(":%s", cfg.Port), cfg.CertFile, cfg.KeyFile)).Msg("failed to bind port")
}
