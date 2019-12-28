// codechain establishes code trust through hash chains.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/command"
	"github.com/frankbraun/codechain/util/homedir"
)

func usage() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage: %s treehash [-l]\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s keygen [-s seckey.bin]\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s keyfile [-l] -s seckey.bin [-c]\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s start -s seckey.bin\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s publish [-s seckey.bin]\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s review [-a] [-d] [-s seckey.bin] [treehash]\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s addkey [-w] pubkey signature [comment]\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s remkey pubkey\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s sigctl -m\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s createdist -f dist.tar.gz\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s apply\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s status\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s cleanslate\n", cmd)
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	argv0 := os.Args[0] + " " + os.Args[1]
	args := os.Args[2:]
	var err error
	switch os.Args[1] {
	case "treehash":
		err = command.TreeHash(argv0, args...)
	case "keygen":
		err = command.KeyGen("codechain", homedir.Codechain(), argv0, args...)
	case "keyfile":
		err = command.KeyFile("codechain", homedir.Codechain(), argv0, args...)
	case "start":
		err = command.Start(argv0, args...)
	case "publish":
		err = command.Publish(argv0, args...)
	case "review":
		err = command.Review(argv0, args...)
	case "addkey":
		err = command.AddKey(argv0, args...)
	case "remkey":
		err = command.RemKey(argv0, args...)
	case "sigctl":
		err = command.SigCtl(argv0, args...)
	case "createdist":
		err = command.CreateDist(argv0, args...)
	case "apply":
		err = command.Apply(argv0, args...)
	case "status":
		err = command.Status(argv0, args...)
	case "cleanslate":
		err = command.CleanSlate(argv0, args...)
	default:
		usage()
	}
	if err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
			os.Exit(1)
		}
		os.Exit(2)
	}
}
