package peering

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"

	"github.com/pickledrick/vpc-peering-operator/pkg/apis/r4/v1"

	"reflect"
)

type PeeringModifier interface {
	UpdateStatus(o *v1.VpcPeering, status string) error
}

type VpcPeeringModifier struct{}

func New() *VpcPeeringModifier {
	return &VpcPeeringModifier{}
}

func (h *VpcPeeringModifier) UpdateStatus(o *v1.VpcPeering, status string) error {

	var err error
	if !reflect.DeepEqual(o.Status.Status, status) || !reflect.DeepEqual(o.Status.Status, "active") {
		o.Status.Status = status
		err = sdk.Update(o)
	}
	if err != nil {
		return err
	}
	return err
}
