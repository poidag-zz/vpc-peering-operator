package watcher

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-kit/kit/log"
	"github.com/pickledrick/vpc-peering-operator/pkg/amazon"
	"github.com/pickledrick/vpc-peering-operator/pkg/apis/r4/v1"
	"github.com/pickledrick/vpc-peering-operator/pkg/peering"
	"github.com/pickledrick/vpc-peering-operator/pkg/wiring"
	"time"
)

const (
	logInitWatch    = "waiting for peering to become active"
	logStatusActive = "peering became active"
	logUpdateRoutes = "updating vpc route tables"
	logTimedOut     = "timed out waiting for peering to become active"
	logCidrMismatch = "peering became active with cidr mismatch"
)

type Watcher interface {
	Watch(o *v1.VpcPeering)
}

type VpcPeeringWatcher struct {
	cfg    *wiring.Config
	logger log.Logger
}

func New(cfg *wiring.Config, log log.Logger) *VpcPeeringWatcher {
	return &VpcPeeringWatcher{
		cfg:    cfg,
		logger: log,
	}
}

func (w *VpcPeeringWatcher) Watch(o *v1.VpcPeering) {
	m := peering.New()

	interval := time.Duration(w.cfg.Poller.WaitSeconds) * time.Second

	w.logger.Log("msg", logInitWatch)

	c, err := amazon.New()
	if err != nil {
		w.logger.Log("err", err)
	}

	input := ec2.DescribeVpcPeeringConnectionsInput{
		VpcPeeringConnectionIds: []*string{o.Status.PeeringId},
	}

	r := 0
	for r < w.cfg.Poller.Retries {
		peer, err := c.Session.DescribeVpcPeeringConnections(&input)
		if err != nil {
			w.logger.Log("err", err)
		}

		if peer.VpcPeeringConnections != nil {
			p := *peer.VpcPeeringConnections[0]
			switch status := *p.Status.Code; status {
			case "active":
				w.logger.Log("msg", logStatusActive)
				if o.Spec.PeerCIDR != *p.AccepterVpcInfo.CidrBlock {
					w.logger.Log("msg", logCidrMismatch)
					m.UpdateStatus(o, "active-cidr-mismatch")
					return
				}
				if w.cfg.ManageRoutes {
					w.logger.Log("msg", logUpdateRoutes)
					err := c.CreateRoutes(o)
					if err != nil {
						w.logger.Log("err", err)
						m.UpdateStatus(o, "failed-adding-routes")
					} else {
						m.UpdateStatus(o, status)
						return
					}
				}
			default:
				m.UpdateStatus(o, status)
			}

			time.Sleep(interval)
		}
		r++
	}
	w.logger.Log("msg", logTimedOut)

	m.UpdateStatus(o, "timed-out")

}
