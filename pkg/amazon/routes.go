package amazon

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gitlab.com/pickledrick/vpc-peering-operator/pkg/apis/r4/v1"
)

func (c *AwsClient) CreateRoutes(o *v1.VpcPeering) error {

	reqKey := "vpc-id"
	reqVal := o.Spec.SourceVpcId
	filter := []*ec2.Filter{
		{
			Name:   &reqKey,
			Values: []*string{&reqVal},
		},
	}

	input := ec2.DescribeRouteTablesInput{
		Filters: filter,
	}

	rtbs, err := c.Session.DescribeRouteTables(&input)
	if err != nil {
		return err
	}

	for _, rtb := range rtbs.RouteTables {
		routeInput := ec2.CreateRouteInput{
			DestinationCidrBlock:   &o.Spec.PeerCIDR,
			VpcPeeringConnectionId: o.Status.PeeringId,
			RouteTableId:           rtb.RouteTableId,
		}

		for _, r := range rtb.Routes {
			if r.DestinationCidrBlock == &o.Spec.PeerCIDR {
				err = fmt.Errorf("routes already exist")
				return err
			}
		}
		_, err := c.Session.CreateRoute(&routeInput)
		if err != nil {
			return err
		}

	}
	return err
}

func (c *AwsClient) DeleteRoutes(o *v1.VpcPeering) error {

	reqKey := "vpc-id"
	reqVal := o.Spec.SourceVpcId
	filter := []*ec2.Filter{
		{
			Name:   &reqKey,
			Values: []*string{&reqVal},
		},
	}

	input := ec2.DescribeRouteTablesInput{
		Filters: filter,
	}

	rtbs, err := c.Session.DescribeRouteTables(&input)
	if err != nil {
		return err
	}

	for _, rtb := range rtbs.RouteTables {
		routeInput := ec2.DeleteRouteInput{
			DestinationCidrBlock: &o.Spec.PeerCIDR,
			RouteTableId:         rtb.RouteTableId,
		}

		_, err := c.Session.DeleteRoute(&routeInput)
		if err != nil {
			return err
		}

	}
	return err
}
