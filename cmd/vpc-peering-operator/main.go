package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"github.com/pickledrick/vpc-peering-operator/pkg/handler"
	"github.com/pickledrick/vpc-peering-operator/pkg/wiring"
	"github.com/vrischmann/envconfig"
	"os"
)

const component = "vpc-peering-operator"

var (
	cfg    wiring.Config
	logger log.Logger
)

var (
	version  = "0.0.1"
	resource = "r4.vc/v1"
	kind     = "VpcPeering"
)

func main() {

	sdk.ExposeMetricsPort()

	watchNamespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		level.Error(logger).Log("msg", "Error fetching watch namespaces", "err", err.Error())
	}

	var resyncPeriod time.Duration
	if cfg.WatchAllNamespaces {
		watchNamespace = ""
	}

	sdk.Watch(resource, kind, watchNamespace, resyncPeriod)
	sdk.Handle(handler.New(&cfg, logger))
	sdk.Run(context.TODO())

}

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger.Log("msg", "Starting VPC Peering Operator version", "version", version)

	if err := envconfig.Init(&cfg); err != nil {
		level.Error(logger).Log("msg", "Error loading config: %s", "err", err.Error())
	}

}
