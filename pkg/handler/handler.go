package handler

import (
	"context"

	"github.com/operator-framework/operator-sdk/pkg/sdk"

	"github.com/go-kit/kit/log"
	"github.com/pickledrick/vpc-peering-operator/pkg/amazon"
	"github.com/pickledrick/vpc-peering-operator/pkg/apis/r4/v1"
	"github.com/pickledrick/vpc-peering-operator/pkg/watcher"
	"github.com/pickledrick/vpc-peering-operator/pkg/wiring"

	"reflect"
)

const (
	logCreateSuccess = "successfully created peering"
	logDeleteSuccess = "sucessfully deleted peering"
)

type Handler interface {
	Handle(ctx context.Context, event sdk.Event) error
}

type VpcPeeringHandler struct {
	cfg    *wiring.Config
	logger log.Logger
	client *amazon.AwsClient
}

func New(cfg *wiring.Config, logger log.Logger) sdk.Handler {
	c, err := amazon.New()
	if err != nil {
		logger.Log("err", err)
	}
	return &VpcPeeringHandler{
		cfg:    cfg,
		logger: logger,
		client: c,
	}
}

func (h *VpcPeeringHandler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1.VpcPeering:
		vpcpeering := o
		eventLogger := log.With(h.logger, "namespace", vpcpeering.GetNamespace(), "vpcpeering", vpcpeering.GetName())
		if event.Deleted {
			deleteLogger := log.With(eventLogger, "action", "delete", "peering-id", vpcpeering.Status.PeeringId)
			if h.cfg.ManageRoutes && vpcpeering.Status.Status == "active" {
				err := h.client.DeleteRoutes(o)
				if err != nil {
					deleteLogger.Log("err", err)
				}
			}
			_, err := h.client.DeletePeering(o)
			if err != nil {
				deleteLogger.Log("err", err)
			}
			deleteLogger.Log("msg", logDeleteSuccess)
			return nil
		}

		createLogger := log.With(eventLogger, "action", "create")

		if vpcpeering.Status.Status == "" {
			p, err := h.client.CreatePeering(o)
			if err != nil {
				createLogger.Log("err", err)
			}

			if !reflect.DeepEqual(p.VpcPeeringConnection.VpcPeeringConnectionId, vpcpeering.Status.PeeringId) {

				vpcpeering.Status.PeeringId = p.VpcPeeringConnection.VpcPeeringConnectionId
				vpcpeering.Status.Status = "requested"
				createLogger = log.With(createLogger, "peering-id", vpcpeering.Status.PeeringId)

				err := sdk.Update(vpcpeering)
				if err != nil {
					createLogger.Log("err", err)
				}

				createLogger.Log("msg", logCreateSuccess)

				w := watcher.New(h.cfg, log.With(eventLogger, "action", "watch", "peering-id", vpcpeering.Status.PeeringId))

				go w.Watch(o)

			}
		}
	}

	return nil
}
