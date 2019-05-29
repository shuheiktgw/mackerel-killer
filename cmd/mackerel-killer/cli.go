package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/shuheiktgw/mackerel-killer/pkg/mkk"

	"github.com/mackerelio/mackerel-client-go"
)

const (
	ExitCodeOK = iota
	ExitCodeError
	ExitCodeParseFlagError
	ExitCodeInvalidFlagError
)

const Name = "mkk"

const EnvMackerelToken = "MACKEREL_API_TOKEN"

var (
	token   string
	hosts   string
	filters string
	dryRun  bool
	quiet   bool
	debug   bool
	version bool
)

type cli struct {
	outStream, errStream io.Writer
	debug                bool
}

func (c *cli) run(args []string) int {
	if err := c.parseFlags(args); err != nil {
		c.printErrorf("Error occurred while parsing flags: %s", err)
		return ExitCodeParseFlagError
	}

	c.setupOutput()

	if version {
		fmt.Fprint(c.outStream, outputVersion())
		return ExitCodeOK
	}

	if err := validateFlags(); err != nil {
		c.printErrorf("Flag validation fails: %s", err)
		return ExitCodeInvalidFlagError
	}

	c.printDebugf("Raw hosts flag: %v", hosts)

	param, err := parseHosts(hosts)
	if err != nil {
		c.printErrorf("Error occurred while parsing hosts query parameters: %s\n", err)
		return ExitCodeInvalidFlagError
	}

	c.printDebugf("Parsed mackerel.FindHostsParam: %v", param)
	c.printDebugf("Raw filters flag: %v", filters)

	fs, err := parseFilters(filters)
	if err != nil {
		c.printErrorf("Error occurred while parsing filters: %s\n", err)
		return ExitCodeInvalidFlagError
	}

	if c.debug {
		for i, f := range fs {
			c.printDebugf("Parsed filter #%d: %v", i, f)
		}
	}

	client := mkk.NewMkk(token)

	c.printInfof("Finding hosts...")
	hs, err := client.FindHosts(param, fs)
	if err != nil {
		c.printErrorf("Error occurred while finding hosts: %s\n", err)
		return ExitCodeError
	}

	if len(hs) > 0 {
		c.printInfof("%d hosts found", len(hs))

		if c.debug {
			for i, h := range hs {
				c.printDebugf("Found host #%d: %v", i, h)
			}
		}
	} else {
		c.printInfof("No hosts found with the specified query parameters and filters")
		return ExitCodeOK
	}

	if dryRun {
		c.printInfof("Running in Dry Run mode")
		c.printInfof("Hosts below will be retired without --dry-run flag\n")

		for i, h := range hs {
			c.printInfof("#%d id: %v, name: %v", i, h.ID, h.Name)
		}

		return ExitCodeOK
	}

	c.printInfof("Retiring hosts...")
	for i, h := range hs {
		if err := client.Kill(h); err != nil {
			c.printErrorf("Error occurred while retiring a host: id: %v, name: %v: %s", h.ID, h.Name, err)
			return ExitCodeError
		}

		c.printInfof("#%v Retired: id: %v, name: %v", i, h.ID, h.Name)
	}

	return ExitCodeOK
}

func (c *cli) parseFlags(args []string) error {
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, usage)
	}

	flags.StringVar(&token, "token", os.Getenv(EnvMackerelToken), "")
	flags.StringVar(&token, "t", os.Getenv(EnvMackerelToken), "")

	flags.StringVar(&hosts, "hosts", "", "")
	flags.StringVar(&hosts, "H", "", "")

	flags.StringVar(&filters, "filters", "", "")
	flags.StringVar(&filters, "F", "", "")

	flags.BoolVar(&dryRun, "dry-run", false, "")
	flags.BoolVar(&dryRun, "d", false, "")

	flags.BoolVar(&quiet, "quiet", false, "")

	flags.BoolVar(&debug, "debug", false, "")

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	return flags.Parse(args[1:])
}

func (c *cli) setupOutput() {
	if quiet {
		c.errStream = ioutil.Discard
	}

	if debug {
		c.debug = true
		c.printDebugf("Running in DEBUG mode")
	}
}

func validateFlags() error {
	if len(token) == 0 {
		return fmt.Errorf("missing Mackerel API token\n"+
			"Please set it via `%s` environment variable or `-t` option\n", EnvMackerelToken)
	}

	if len(filters) == 0 {
		return fmt.Errorf("missing filters\n" +
			"Please set it via `-F` option\n")
	}

	return nil
}

func parseHosts(hosts string) (*mackerel.FindHostsParam, error) {
	var p mackerel.FindHostsParam

	if len(hosts) == 0 {
		return &p, nil
	}

	if err := json.Unmarshal([]byte(hosts), &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func parseFilters(filters string) ([]mkk.Filter, error) {
	var arr map[string][]json.RawMessage
	if err := json.Unmarshal([]byte(filters), &arr); err != nil {
		return nil, err
	}

	var fs []mkk.Filter
	for k, v := range arr {
		switch k {
		case "GracePeriodFilter":
			for i, attr := range v {
				var f mkk.GracePeriodFilter
				if err := json.Unmarshal(attr, &f); err != nil {
					return nil, errors.Wrapf(err, "error occurred while unmarshaling %dth attribute of %s", i, k)
				}

				fs = append(fs, &f)
			}
		case "HostFilter":
			for i, attr := range v {
				var f mkk.HostFilter
				if err := json.Unmarshal(attr, &f); err != nil {
					return nil, errors.Wrapf(err, "error occurred while unmarshaling %dth attribute of %s", i, k)
				}

				fs = append(fs, &f)
			}
		case "MetricAbsenceFilter":
			for i, attr := range v {
				var f mkk.MetricAbsenceFilter
				if err := json.Unmarshal(attr, &f); err != nil {
					return nil, errors.Wrapf(err, "error occurred while unmarshaling %dth attribute of %s", i, k)
				}

				fs = append(fs, &f)
			}
		default:
			return nil, fmt.Errorf("filter named `%s` does not exist", k)
		}
	}

	return fs, nil
}

func (c *cli) printDebugf(format string, args ...interface{}) {
	if c.debug {
		fmt.Fprintf(c.outStream, fmt.Sprintf("[mkk][DEBUG] %s\n", format), args...)
	}
}

func (c *cli) printErrorf(format string, args ...interface{}) {
	fmt.Fprintf(c.errStream, fmt.Sprintf("[mkk][ERROR] %s\n", format), args...)
}

func (c *cli) printInfof(format string, args ...interface{}) {
	fmt.Fprintf(c.outStream, fmt.Sprintf("[mkk] %s\n", format), args...)
}

var usage = `mkk - Retire inactive Mackerel hosts

Synopsis:
  $ mkk --hosts '{"name":"hostName"}' --filters '{"MetricExistenceFilter":[{"name":"loadavg5","from":155891000,"to":155895000}]}'

Options:
  --debug        prints debug message
  --dry-run, -d  runs mkk without actually retiring the hosts  
  --filters, -F  specifies filters and its attributes in JSON
  --help, -h     prints help
  --hosts, -H    specifies query parameters to find hosts in JSON
  --quiet        stops printing messages to stdout
  --token, -t    specifies Mackerel API token
  --version, -v  prints the current version

`
