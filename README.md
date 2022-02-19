# route53-ddns

A small, sharp tool to update a [Route 53](https://aws.amazon.com/route53/) domain with the public IP of the client host.

## Arguments
```
Flags:
  -h, --help                     Show context-sensitive help.
      --check-interval=0s        If greater than zero, the record set will be updated on this interval ($CHECK_INTERVAL)
      --host-name=STRING         The hostname to be updated ($HOST_NAME)
      --hosted-zone-id=STRING    The Route53 Hosted Zone ID holding the hostname record ($HOSTED_ZONE_ID)
```

### Notes
* The `--check-interval` flag should be in the format of a Go `time.Duration`, as understood by [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration)
* The `--host-name` field should be the fully qualified domain name of the record, not the subdomain (that is, pass `vpn.myexample.com` rather than just `vpn`)
* This tool is dependent on a valid AWS environment. This can be achieved via `AWS_*` environment variables or by running `aws configure` using the CLI to set up a `default` profile.

## License

This application is licensed under the Apache 2.0 License

## Contributing

Pull requests are welcome! Please open an issue first so that it can be discussed.