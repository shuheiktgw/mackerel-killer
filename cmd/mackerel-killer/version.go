package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/tcnksm/go-latest"
)

// The name of this repository
const RepoName = "mackerel-killer"

// The current version of mackerel-killer
const Version = "0.0.1"

// The owner of mackerel-killer
const RepoOwner = "shuheiktgw"

// outputVersion outputs current version of mackerel-killer. It also checks
// the latest release and adds a warning to update mackerel-killer
func outputVersion() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s current version v%s\n", Name, Version)

	// Get the latest release
	verCheckCh := make(chan *latest.CheckResponse)
	go func() {
		githubTag := &latest.GithubTag{
			Owner:      RepoOwner,
			Repository: RepoName,
		}

		res, err := latest.Check(githubTag, Version)

		// Ignore the error
		if err != nil {
			return
		}

		verCheckCh <- res
	}()

	select {
	case <-time.After(2 * time.Second):
	case res := <-verCheckCh:
		if res.Outdated {
			fmt.Fprintf(&b, "The latest version is v%s, please update\n", res.Current)
		}
	}

	return b.String()
}
