package main

import (
	"fmt"
	bugyoclient "github.com/tomtwinkle/bugyo-client-go"
	"github.com/tomtwinkle/bugyo-client-go/bugyocli"
	"github.com/urfave/cli"
	"log"
	"os"
)

var version = "unknown"
var revision = "unknown"

func main() {
	bcli := bugyocli.NewCLI()
	app := cli.NewApp()
	app.Name = "Bugyo Client CLI for Go"
	app.Usage = "奉行クラウドCLI"
	app.Author = "tomtwinkle"
	app.Version = fmt.Sprintf("bugyo-client-go cli version %s.rev-%s", version, revision)
	app.Commands = []cli.Command{
		{
			Name:      "punchmark",
			ShortName: "pm",
			Usage: "タイムレコーダー記録を行う",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "type, t",
					Usage:    "出勤: --type in or -t in" +
						"\n\t退勤: --type out or -t out" +
						"\n\t外出: --type go or -t go" +
						"\n\t再入: --type return or -t return",
					Required: true,
					Value:    "",
				},
			},
			Action: func(c *cli.Context) error {
				switch c.String("type") {
				case "in":
					return bcli.PunchMark(bugyoclient.ClockTypeClockIn)
				case "out":
					return bcli.PunchMark(bugyoclient.ClockTypeClockOut)
				case "go":
					return bcli.PunchMark(bugyoclient.ClockTypeGoOut)
				case "return":
					return bcli.PunchMark(bugyoclient.ClockTypeReturned)
				default:
					return cli.ShowSubcommandHelp(c)
				}
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
