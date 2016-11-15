package main

import (
	"os"
	"testing"
)

func TestMainWithYAML(t *testing.T) {
	os.Args = []string{
		"kubemerge",
		"fixtures/source.yml",
		"fixtures/deployment.yml",
	}

	main()
}

func TestMainWithJSON(t *testing.T) {
	os.Args = []string{
		"kubemerge",
		"fixtures/source.json",
		"fixtures/deployment.json",
	}

	main()
}
