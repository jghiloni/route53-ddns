package dns

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"golang.org/x/net/context"
)

type UpdateResult struct {
	FQDN string
	IP   string
}

func UpdateDNSRecord(hostedZoneID string, hostname string, clientResolver ClientIPResolver, r53 Route53Client) (UpdateResult, error) {
	clientIP, err := clientResolver.ResolveClientIP()
	if err != nil {
		return UpdateResult{}, fmt.Errorf("could not resolve client IP: %w", err)
	}

	if err != nil {
		return UpdateResult{}, fmt.Errorf("could not get hosted zone: %w", err)
	}

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &hostedZoneID,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(route53.ChangeActionUpsert),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Type: aws.String("A"),
						Name: aws.String(hostname),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(clientIP),
							},
						},
						TTL: aws.Int64(60),
					},
				},
			},
		},
	}

	output, err := r53.ChangeResourceRecordSets(input)
	if err != nil {
		return UpdateResult{}, fmt.Errorf("could not modify A record: %w", err)
	}

	if aws.StringValue(output.ChangeInfo.Status) == route53.ChangeStatusInsync {
		return UpdateResult{
			FQDN: hostname,
			IP:   clientIP,
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	in := &route53.GetChangeInput{
		Id: output.ChangeInfo.Id,
	}

	for {
		out, err := r53.GetChange(in)
		if err != nil {
			return UpdateResult{}, fmt.Errorf("failed to get change status: %w", err)
		}

		if aws.StringValue(out.ChangeInfo.Status) == route53.ChangeStatusInsync {
			return UpdateResult{
				FQDN: hostname,
				IP:   clientIP,
			}, nil
		}

		select {
		case <-time.After(5 * time.Second):
		case <-ctx.Done():
			return UpdateResult{}, errors.New("operation timed out or cancelled")
		}
	}
}
