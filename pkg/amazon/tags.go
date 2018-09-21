package amazon

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateTags - Creates tags on resources relating to a vpcpeering
func (c *AwsClient) CreateTags(resources []*string, tags []*ec2.Tag) (*ec2.CreateTagsOutput, error) {

	request := &ec2.CreateTagsInput{
		Resources: resources,
		Tags:      tags,
	}

	return c.Session.CreateTags(request)
}
