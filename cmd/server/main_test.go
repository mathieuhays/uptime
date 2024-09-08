package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("missing port check", func(t *testing.T) {
		getEnv := func(name string) string {
			if name == "PORT" {
				return ""
			}

			return "lskjdf"
		}

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})
		err := run(getEnv, stdout, stderr)

		if !errors.Is(err, errMissingPort) {
			t.Errorf("Port error expected. Got different error: %s", err)
		}
	})
}
