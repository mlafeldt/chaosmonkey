package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/simpledb"
)

// AutoScalingGroup describes an AWS auto scaling group.
type AutoScalingGroup struct {
	Name            string
	CurrentSize     int
	DesiredCapacity int
	MinSize         int
	MaxSize         int
}

func autoScalingGroups() ([]AutoScalingGroup, error) {
	var groups []AutoScalingGroup
	svc := autoscaling.New(session.New())
	err := svc.DescribeAutoScalingGroupsPages(nil, func(out *autoscaling.DescribeAutoScalingGroupsOutput, last bool) bool {
		for _, g := range out.AutoScalingGroups {
			groups = append(groups, AutoScalingGroup{
				Name:            aws.StringValue(g.AutoScalingGroupName),
				CurrentSize:     len(g.Instances),
				DesiredCapacity: int(aws.Int64Value(g.DesiredCapacity)),
				MinSize:         int(aws.Int64Value(g.MinSize)),
				MaxSize:         int(aws.Int64Value(g.MaxSize)),
			})
		}
		return !last
	})
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func deleteSimpleDBDomain(domainName string) error {
	var domainExists bool
	svc := simpledb.New(session.New())
	err := svc.ListDomainsPages(nil, func(out *simpledb.ListDomainsOutput, last bool) bool {
		for _, n := range out.DomainNames {
			if aws.StringValue(n) == domainName {
				domainExists = true
			}
		}
		return !last
	})
	if !domainExists {
		return fmt.Errorf("SimpleDB domain %q does not exist", domainName)
	}
	_, err = svc.DeleteDomain(&simpledb.DeleteDomainInput{
		DomainName: aws.String(domainName),
	})
	return err
}
