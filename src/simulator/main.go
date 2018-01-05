package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
	"fmt"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "simulator"
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   "web_api_addr",
			EnvVar: "WEB_API_ADDR",
			Value:  "http://192.168.1.206:12301",
			Usage:  "WebAPI连接地址",
		},
		cli.IntFlag{
			Name:  "concurrence",
			Value: 1,
			Usage: "并发数",
		},

	)
	app.Commands = []cli.Command{
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action:  func(c *cli.Context) error {
				fmt.Println("complete do")
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action:  func(c *cli.Context) error {
				fmt.Println("add do")
				return nil
			},
		},
	}

	app.Action = func(ctx *cli.Context) error {
		arr := strings.Split(ctx.String("web_api_addr"), ",")
		if len(arr) == 0 {
			return errors.New("redis_addrs value not correct")
		}

		NewSimulator(ctx.String("web_api_addr")).Start(ctx.Int("concurrence"))

		return nil
	}

	app.Run(os.Args)
}
