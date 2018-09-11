package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	ec2meta "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Client interface {
	DeleteRoutes()
	DeletePeering()
	CreateRoutes()
	CreatePeering()
}

type AwsClient struct {
	Session *ec2.EC2
}

func New() (*AwsClient, error) {
	var sess *ec2.EC2
	var client = AwsClient{Session: sess}
	s := session.Must(session.NewSession())
	meta := ec2meta.New(s)
	awsRegion, err := meta.Region()
	if err != nil {
		return &client, err
	}
	client.Session = ec2.New(s, &aws.Config{Region: aws.String(awsRegion)})
	return &client, err
}
