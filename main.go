package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	flagSet *flag.FlagSet
	toYAML  bool
)

func init() {
	// workaround to avoid inheriting vendor flags
	flagSet = flag.NewFlagSet("kubemerge", flag.ExitOnError)
	flagSet.BoolVar(&toYAML, "yaml", false, "Output as YAML")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] source target\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `Examples:

  kubemerge policy.yaml deployment.yaml
  kubemerge -yaml policy.yaml deployment.yaml
	cat deployment.yaml | kubemerge policy.yaml

Options:
`)
		flagSet.PrintDefaults()
	}
}

func main() {
	var source *os.File
	var target *os.File
	var err error

	if err = flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// Read in source file as first argument
	switch name := flagSet.Arg(0); {
	case name == "":
		log.Fatal("Must supply a source file for merging")
	default:
		source, err = os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer source.Close()
	}

	// Read in target file as second argument
	switch name := flagSet.Arg(1); {
	case name == "":
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}
		// Print usage unless we already have STDIN data or incoming pipe
		if fi.Size() == 0 && fi.Mode()&os.ModeNamedPipe == 0 {
			flagSet.Usage()
			return
		}
		target = os.Stdin
	default:
		if target, err = os.Open(name); err != nil {
			log.Fatal(err)
		}
		defer target.Close()
	}

	sourceResource, err := ParseDoc(source)
	if err != nil {
		log.Fatal(err)
	}

	targetResources, err := ParseInlineDocs(target)
	if err != nil {
		log.Fatal(err)
	}

	for _, resource := range targetResources {
		var result interface{}
		switch resource.Kind {
		case "Deployment":
			result, err = resource.MergeDeployment(sourceResource.Data)
		case "DaemonSet":
			result, err = resource.MergeDaemonSet(sourceResource.Data)
		case "ReplicaSet":
			result, err = resource.MergeReplicaSet(sourceResource.Data)
		case "ReplicationController":
			result, err = resource.MergeRC(sourceResource.Data)
		default:
			result, err = resource.UnmarshalGeneric()
		}

		if err != nil {
			log.Fatal(err)
		}

		if err = printResource(result, toYAML); err != nil {
			log.Fatal(err)
		}
	}
}
