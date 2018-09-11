package amazon

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/pickledrick/vpc-peering-operator/pkg/apis/r4/v1"
)

// CreatePeering - Creates an AWS VPC Peering based a new CRD object creation
func (c *AwsClient) CreatePeering(o *v1.VpcPeering) (*ec2.CreateVpcPeeringConnectionOutput, error) {

	request := ec2.CreateVpcPeeringConnectionInput{
		PeerOwnerId: &o.Spec.PeerOwnerId,
		PeerVpcId:   &o.Spec.PeerVpcId,
		PeerRegion:  &o.Spec.PeerRegion,
		VpcId:       &o.Spec.SourceVpcId,
	}

	return c.Session.CreateVpcPeeringConnection(&request)
}

// DeletePeering - Deletes an AWS VPC Peering based a CRD object deletion
func (c *AwsClient) DeletePeering(o *v1.VpcPeering) (*ec2.DeleteVpcPeeringConnectionOutput, error) {

	opts := ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: o.Status.PeeringId,
	}

	return c.Session.DeleteVpcPeeringConnection(&opts)
}
