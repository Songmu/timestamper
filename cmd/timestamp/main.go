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

func main() {
	err := run(os.Args[1:])
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("timestamp", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `timestamp - sample application of Songmu/timestamper
Usage:`)
		fs.PrintDefaults()
	}
	var (
		utc    = fs.Bool("utc", false, "use utc timestamp")
		layout = fs.String("layout", "", "custom timestamp layout")
	)
	if err := fs.Parse(args); err != nil {
		return err
	}
	var opts []timestamper.Option
	if *utc {
		opts = append(opts, timestamper.UTC())
	}
	if *layout != "" {
		opts = append(opts, timestamper.Layout(*layout))
	}
	_, err := io.Copy(os.Stdout,
		transform.NewReader(os.Stdin,
			timestamper.New(opts...)))
	return err
}
