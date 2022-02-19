package dns_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/jghiloni/route53-ddns/dns"
	"github.com/jghiloni/route53-ddns/dns/dnsfakes"
	. "github.com/onsi/gomega"
)

type localhostResolver struct{}

func (l localhostResolver) ResolveClientIP() (string, error) {
	return "127.0.0.1", nil
}

func TestUpdateRecord(t *testing.T) {
	RegisterTestingT(t)

	ipResolver := localhostResolver{}
	r53client := &dnsfakes.FakeRoute53Client{}
	r53client.ChangeResourceRecordSetsReturns(&route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53.ChangeInfo{
			Status: aws.String(route53.ChangeStatusInsync),
		},
	}, nil)

	r53client.GetHostedZoneReturns(&route53.GetHostedZoneOutput{
		HostedZone: &route53.HostedZone{
			Name: aws.String("mynetwork.com"),
		},
	}, nil)

	result, err := dns.UpdateDNSRecord("1234", "vpn.mynetwork.com", ipResolver, r53client)
	Expect(err).NotTo(HaveOccurred())
	Expect(result.FQDN).To(Equal("vpn.mynetwork.com"))
	Expect(result.IP).To(Equal("127.0.0.1"))
}

func TestUpdateRecordWithLongResolution(t *testing.T) {
	RegisterTestingT(t)

	ipResolver := localhostResolver{}
	r53client := &dnsfakes.FakeRoute53Client{}
	r53client.ChangeResourceRecordSetsReturns(&route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53.ChangeInfo{
			Status: aws.String(route53.ChangeStatusPending),
		},
	}, nil)

	r53client.GetHostedZoneReturns(&route53.GetHostedZoneOutput{
		HostedZone: &route53.HostedZone{
			Name: aws.String("mynetwork.com"),
		},
	}, nil)

	r53client.GetChangeReturnsOnCall(0, &route53.GetChangeOutput{
		ChangeInfo: &route53.ChangeInfo{
			Status: aws.String(route53.ChangeStatusPending),
		},
	}, nil)
	r53client.GetChangeReturnsOnCall(1, &route53.GetChangeOutput{
		ChangeInfo: &route53.ChangeInfo{
			Status: aws.String(route53.ChangeStatusPending),
		},
	}, nil)
	r53client.GetChangeReturnsOnCall(2, &route53.GetChangeOutput{
		ChangeInfo: &route53.ChangeInfo{
			Status: aws.String(route53.ChangeStatusInsync),
		},
	}, nil)

	start := time.Now()
	result, err := dns.UpdateDNSRecord("1234", "vpn.mynetwork.com", ipResolver, r53client)
	dur := time.Since(start)

	Expect(err).NotTo(HaveOccurred())
	Expect(result.FQDN).To(Equal("vpn.mynetwork.com"))
	Expect(result.IP).To(Equal("127.0.0.1"))
	Expect(dur > (10 * time.Second)).To(BeTrue())
}
