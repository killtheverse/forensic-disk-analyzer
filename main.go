package main

import (	
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "boot_info"
	app.Usage = "Analyzes the MBR/GPT information of forensic images"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "filepath",
			Aliases: []string{"f"},
			Usage: "`path` of the forensic image file",
			Required: true,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		StoreHashes(ctx.String("filepath"))
		AnalyzeImage(ctx.String("filepath"))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
