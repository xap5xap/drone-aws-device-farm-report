package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var build string

func main() {
	app := cli.NewApp()
	app.Name = "AWS device farm drone plugin"
	app.Usage = "AWS device farm drone plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.0+%s", build)
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "access-key",
			Usage:  "aws access key",
			EnvVar: "PLUGIN_ACCESS_KEY,AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  "aws secret key",
			EnvVar: "PLUGIN_SECRET_KEY,AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "region",
			Usage:  "aws region",
			Value:  "us-west-2",
			EnvVar: "PLUGIN_REGION",
		},
		cli.StringFlag{
			Name:   "test-project",
			Usage:  "Name of the AWS Device farm project where you want to upload the app, tests, and schedule the run",
			EnvVar: "PLUGIN_TEST_PROJECT",
		},
		cli.StringFlag{
			Name:   "run-name",
			Usage:  "Name of the Run to check the status",
			EnvVar: "PLUGIN_RUN_NAME",
		},
		cli.BoolTFlag{
			Name:   "yaml-verified",
			Usage:  "Ensure the yaml was signed",
			EnvVar: "DRONE_YAML_VERIFIED",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func run(c *cli.Context) error {
	plugin := Plugin{
		Key:          c.String("access-key"),
		Secret:       c.String("secret-key"),
		Region:       c.String("region"),
		TestProject:  c.String("test-project"),
		RunName:      c.String("run-name"),
		YamlVerified: c.BoolT("yaml-verified"),
	}

	return plugin.Exec()
}
