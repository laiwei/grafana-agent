package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-kit/log/level"
	"github.com/grafana/agent/pkg/util"
	smoke "github.com/grafana/agent/tools/smoke/internal"
	"github.com/weaveworks/common/logging"
)

func main() {
	var (
		cfg         smoke.Config
		logLevel    logging.Level
		logFormat   logging.Format
		withTimeout time.Duration
	)

	cfg.RegisterFlags(flag.CommandLine)
	logLevel.RegisterFlags(flag.CommandLine)
	logFormat.RegisterFlags(flag.CommandLine)
	flag.DurationVar(&withTimeout, "duration", time.Duration(0), "test duration")
	flag.Parse()

	logger := util.NewLoggerFromLevel(logLevel, logFormat)

	ctx := context.Background()
	if withTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, withTimeout)
		defer cancel()
		level.Debug(logger).Log("msg", "running with duration", "duration", withTimeout.String())
	}

	level.Info(logger).Log("msg", "starting smoke test")
	smokeTest, err := smoke.New(logger, cfg)
	if err != nil {
		level.Error(logger).Log("msg", "error constructing smoke test", "err", err)
		os.Exit(1)
	}
	if err := smokeTest.Run(ctx); err != nil {
		level.Error(logger).Log("msg", "smoke test run failure", "err", err)
		os.Exit(1)
	}
}
