package dns

import "github.com/aws/aws-sdk-go/service/route53"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Route53Client
type Route53Client interface {
	ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
	GetChange(input *route53.GetChangeInput) (*route53.GetChangeOutput, error)
	GetHostedZone(input *route53.GetHostedZoneInput) (*route53.GetHostedZoneOutput, error)
}
