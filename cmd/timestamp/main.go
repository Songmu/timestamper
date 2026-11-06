package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Songmu/timestamper"
	"golang.org/x/text/transform"
)

var (
	version  = "dev"
	revision = "HEAD"
)

func main() {
	err := run(os.Args[1:], os.Stdout)
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("timestamp", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `timestamp - sample application of Songmu/timestamper
Usage:`)
		fs.PrintDefaults()
	}
	var (
		utc         = fs.Bool("utc", false, "use utc timestamp")
		layout      = fs.String("layout", "", "custom timestamp layout")
		showVersion = fs.Bool("version", false, "show version and revision")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *showVersion {
		_, err := fmt.Fprintf(out, "timestamp %s (rev: %s)\n", version, revision)
		return err
	}
	var opts []timestamper.Option
	if *utc {
		opts = append(opts, timestamper.UTC())
	}
	if *layout != "" {
		opts = append(opts, timestamper.Layout(*layout))
	}
	_, err := io.Copy(out,
		transform.NewReader(os.Stdin,
			timestamper.New(opts...)))
	return err
}
