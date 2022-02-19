package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/jghiloni/route53-ddns/dns"
)

var cli struct {
	CheckInterval time.Duration `help:"If greater than zero, the record set will be updated on this interval" env:"CHECK_INTERVAL" default:"0s"`
	HostName      string        `help:"The hostname to be updated" env:"HOST_NAME"`
	HostedZoneID  string        `help:"The Route53 Hosted Zone ID holding the hostname record" env:"HOSTED_ZONE_ID"`
}

func main() {
	kong.Parse(&cli)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan
		cancel()
	}()

	route53Client := route53.New(session.Must(session.NewSession()))
	clientIPResolver := dns.NewClientIPResolver()

	var t *time.Ticker
	if cli.CheckInterval > 0 {
		t = time.NewTicker(cli.CheckInterval)
	}

	for {
		if _, err := dns.UpdateDNSRecord(cli.HostedZoneID, cli.HostName, clientIPResolver, route53Client); err != nil {
			log.Fatal(err)
		}

		if t == nil {
			os.Exit(0)
		}

		select {
		case <-t.C:
		case <-ctx.Done():
			log.Fatal("interrupted")
		}
	}
}
