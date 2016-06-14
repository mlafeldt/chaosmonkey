package main

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/simpledb"
)

func autoScalingGroups() ([]string, error) {
	var groups []string
	svc := autoscaling.New(session.New())
	err := svc.DescribeAutoScalingGroupsPages(nil, func(out *autoscaling.DescribeAutoScalingGroupsOutput, last bool) bool {
		for _, g := range out.AutoScalingGroups {
			groups = append(groups, aws.StringValue(g.AutoScalingGroupName))
		}
		return !last
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(groups)
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
