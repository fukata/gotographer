package main

import (
	"os"
	"runtime"
	"fmt"
	"log"
	"strings"
	"sync"
	"github.com/codegangsta/cli"
)

type Option struct {
	From string
	To string
	PhotoDir string
	VideoDir string
	Recursive bool
	DryRun bool
	Excludes []string
	Concurrency int
}

func main() {
	app := cli.NewApp()
	app.Name = "ghoto"
	app.Usage = "Transfer photo(video)"
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "from",
			Value: "/path/to/src",
			Usage: "Source directory",
		},
		cli.StringFlag {
			Name: "to",
			Value: "/path/to/dst",
			Usage: "Destination directory",
		},
		cli.StringFlag {
			Name: "photo-dir, P",
			Value: "photo",
			Usage: "Destination photo directory",
		},
		cli.StringFlag {
			Name: "video-dir, V",
			Value: "video",
			Usage: "Destination video directory",
		},
		cli.StringFlag {
			Name: "exclude, x",
			Value: "",
			Usage: "Exclude dir/file separate comma.",
		},
		cli.IntFlag {
			Name: "concurrency, c",
			Value: runtime.NumCPU(),
			Usage: "Concurrency num.",
		},
		cli.BoolFlag {
			Name: "recursive, r",
			Usage: "Resursive",
		},
		cli.BoolFlag {
			Name: "dry-run",
			Usage: "Dry Run",
		},
	}
	app.Action = func(c *cli.Context) {
		// options
		option := &Option{
			c.String("from"),
			c.String("to"),
			c.String("photo-dir"),
			c.String("video-dir"),
			c.Bool("recursive"),
			c.Bool("dry-run"),
			strings.Split(c.String("exclude"), ","),
			c.Int("concurrency"),
		}

		// check path
		isDir, err := IsDirectory(option.From)
		if err != nil {
			log.Fatal(err)
		}

		if isDir != false {
			fmt.Errorf("%s is not found.", option.From)
		}

		isDir, err = IsDirectory(option.To)
		if err != nil {
			log.Fatal(err)
		}

		if isDir != false {
			fmt.Errorf("%s is not found.", option.To)
		}

		// move
		var wg sync.WaitGroup
		ch := make(chan int, option.Concurrency)
		MoveFiles(&wg, ch, option.From, option)
		wg.Wait()
	}
	app.Run(os.Args)
}
