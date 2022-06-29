package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/desertbit/grumble"
)

var session2Grumble *Session
var app = grumble.New(&grumble.Config{
	Name:        "config",
	Description: "",

	Flags: func(f *grumble.Flags) {
		// f.String("d", "directory", "DEFAULT", "set an alternative directory path")
		// f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func initInteractive() {
	app.AddCommand(&grumble.Command{
		Name:    "resume",
		Help:    "resume scanning",
		Aliases: []string{"r"},
		Run: func(c *grumble.Context) error {
			c.Stop()
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name:    "stop",
		Help:    "resume scanning",
		Aliases: []string{"s", "exit", "quit", "kill"},
		Run: func(c *grumble.Context) error {
			c.Stop()
			os.Exit(0)
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "fe",
		Help: "filter extensions",

		Args: func(a *grumble.Args) {
			a.String("extensions", "extensions to filter,use +[ext] to append to current", grumble.Default(strings.Join(session2Grumble.extFilters, ",")))
		},

		Run: func(c *grumble.Context) error {
			exts := c.Args.String("extensions")
			if exts != "" {
				if strings.HasPrefix(exts, "+") {
					tmp := strings.Join(session2Grumble.extFilters, ",") + "," + strings.TrimPrefix(exts, "+")
					session2Grumble.extFilters = parseExtFilters(tmp)
				} else {
					session2Grumble.extFilters = parseExtFilters(exts)
				}
			}
			app.Println("Using: " + strings.Join(session2Grumble.extFilters, ","))
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "fc",
		Help: "filter status code",

		Args: func(a *grumble.Args) {
			a.String("statuscode", "statuscode to filter,use +[code] to append to current", grumble.Default(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(session2Grumble.statusFilters)), ","), "[]")))
		},

		Run: func(c *grumble.Context) error {
			codes := c.Args.String("statuscode")
			if codes != "" {
				if strings.HasPrefix(codes, "+") {
					tmp := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(session2Grumble.statusFilters)), ","), "[]") + "," + codes
					session2Grumble.statusFilters = parseStatusFilters(tmp)
				} else {
					session2Grumble.statusFilters = parseStatusFilters(codes)
				}
			}
			app.Println("Using: " + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(session2Grumble.statusFilters)), ","), "[]"))
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "fr",
		Help: "filter url regex",

		Args: func(a *grumble.Args) {
			a.String("regex", "url regex filter,use +[regex] to append to current", grumble.Default(session2Grumble.urlFilter))
		},

		Run: func(c *grumble.Context) error {
			codes := c.Args.String("regex")
			if codes != "" {
				if strings.HasPrefix(codes, "+") {
					session2Grumble.urlFilter = session2Grumble.urlFilter + "|" + strings.TrimPrefix(codes, "+")
				} else {
					session2Grumble.urlFilter = codes
				}
			}
			app.Println("Using: " + session2Grumble.urlFilter)
			return nil
		},
	})
}
