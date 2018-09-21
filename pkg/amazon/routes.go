package amazon

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pickledrick/vpc-peering-operator/pkg/apis/r4/v1"
	"github.com/pickledrick/vpc-peering-operator/pkg/wiring"
	"strings"
)

var (
	routeKey = "vpc-peering-operator.r4.vc/"
)

func (c *AwsClient) CreateRoutes(o *v1.VpcPeering, cfg *wiring.Config) error {

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

	var resources []*string

	for _, rtb := range rtbs.RouteTables {
		resources = append(resources, aws.String(*rtb.RouteTableId))
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

	tags := []*ec2.Tag{
		{
			Key:   aws.String(routeKey + strings.ToLower(o.Namespace) + "-" + strings.ToLower(o.Name)),
			Value: aws.String(o.Spec.PeerCIDR),
		},
	}

	_, err = c.CreateTags(resources, tags)
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
