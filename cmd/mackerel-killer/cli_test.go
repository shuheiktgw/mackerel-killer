package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestCLI_Run(t *testing.T) {
	cases := []struct {
		command           string
		expectedOutStream string
		expectedErrStream string
		expectedExitCode  int
	}{
		{
			command:           "mkk -v",
			expectedOutStream: fmt.Sprintf("mkk current version v%s\n", Version),
			expectedErrStream: "",
			expectedExitCode:  ExitCodeOK,
		},
		{
			command:           `mkk -t aqbc -H {}`,
			expectedOutStream: "",
			expectedErrStream: "missing filters",
			expectedExitCode:  ExitCodeInvalidFlagError,
		},
		{
			command:           `mkk -t aqbc -H {} -F {"UnknownFilter":[{"name":"loadavg5"}]}`,
			expectedOutStream: "",
			expectedErrStream: "filter named `UnknownFilter` does not exist",
			expectedExitCode:  ExitCodeInvalidFlagError,
		},
	}

	for i, tc := range cases {
		outStream := new(bytes.Buffer)
		errStream := new(bytes.Buffer)

		cli := cli{outStream: outStream, errStream: errStream}
		args := strings.Split(tc.command, " ")

		if got, want := cli.run(args), tc.expectedExitCode; got != want {
			t.Errorf("#%v invalid exit code got: %v, want %v", i, got, want)
		}

		if got, want := outStream.String(), tc.expectedOutStream; !strings.Contains(got, want) {
			t.Errorf("#%v invalid outStream: got: %v, want: %v", i, got, want)
		}

		if got, want := errStream.String(), tc.expectedErrStream; !strings.Contains(got, want) {
			t.Errorf("#%v invalid errStream: got: %v, want: %v", i, got, want)
		}
	}
}
