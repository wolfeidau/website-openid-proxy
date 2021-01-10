package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/apex/gateway/v2"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo/v4"
	"github.com/wolfeidau/aws-openid-proxy/internal/flags"
	"github.com/wolfeidau/aws-openid-proxy/internal/flags/app"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

var cfg = new(flags.API)

func main() {
	kong.Parse(cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	e := echo.New()

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
