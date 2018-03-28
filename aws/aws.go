// Package aws provides access to Amazon Web Services (AWS).
// AWS credentials need to be passed via environment variables.
package aws

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/simpledb"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Client is a client to the AWS API.
type Client struct {
	Region string
}

// NewClient returns a new Client.
func NewClient(region string) *Client {
	return &Client{Region: region}
}

// AutoScalingGroup describes an AWS auto scaling group.
type AutoScalingGroup struct {
	Name               string
	InstancesInService int
	DesiredCapacity    int
	MinSize            int
	MaxSize            int
}

// AutoScalingGroups returns a list of all auto scaling groups.
func (c *Client) AutoScalingGroups() ([]AutoScalingGroup, error) {
	sess, err := c.newSession()
	if err != nil {
		return nil, err
	}
	svc := autoscaling.New(sess)

	var groups []AutoScalingGroup
	err = svc.DescribeAutoScalingGroupsPages(nil, func(out *autoscaling.DescribeAutoScalingGroupsOutput, last bool) bool {
		for _, g := range out.AutoScalingGroups {
			inService := 0
			for _, i := range g.Instances {
				if aws.StringValue(i.LifecycleState) == autoscaling.LifecycleStateInService {
					inService++
				}
			}
			groups = append(groups, AutoScalingGroup{
				Name:               aws.StringValue(g.AutoScalingGroupName),
				InstancesInService: inService,
				DesiredCapacity:    int(aws.Int64Value(g.DesiredCapacity)),
				MinSize:            int(aws.Int64Value(g.MinSize)),
				MaxSize:            int(aws.Int64Value(g.MaxSize)),
			})
		}
		return !last
	})
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// DeleteSimpleDBDomain deletes an existing SimpleDB domain.
func (c *Client) DeleteSimpleDBDomain(domainName string) error {
	sess, err := c.newSession()
	if err != nil {
		return err
	}
	svc := simpledb.New(sess)

	var domainExists bool
	err = svc.ListDomainsPages(nil, func(out *simpledb.ListDomainsOutput, last bool) bool {
		for _, n := range out.DomainNames {
			if aws.StringValue(n) == domainName {
				domainExists = true
			}
		}
		return !last
	})
	if err != nil {
		return err
	}
	if !domainExists {
		return fmt.Errorf("SimpleDB domain %q does not exist", domainName)
	}
	_, err1 := svc.DeleteDomain(&simpledb.DeleteDomainInput{
		DomainName: aws.String(domainName),
	})
	return err1
}

func (c *Client) newSession() (*session.Session, error) {
	config := &aws.Config{
		Region:     aws.String(c.Region),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	if role := os.Getenv("AWS_ROLE"); role != "" {
		if err := assumeRole(role, config); err != nil {
			return nil, err
		}
	}

	return session.NewSession(config)
}

func assumeRole(role string, config *aws.Config) error {
	svc := sts.New(session.New(config))
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(role),
		RoleSessionName: aws.String("chaosmonkey"),
		DurationSeconds: aws.Int64(900),
	}
	out, err := svc.AssumeRole(params)
	if err != nil {
		return err
	}

	config.Credentials = credentials.NewStaticCredentials(
		aws.StringValue(out.Credentials.AccessKeyId),
		aws.StringValue(out.Credentials.SecretAccessKey),
		aws.StringValue(out.Credentials.SessionToken),
	)

	return nil
}
