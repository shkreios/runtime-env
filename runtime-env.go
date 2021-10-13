/*
Copyright © 2020 Simon Hessel

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func load(envfile string, prefix string, removePrefix bool, noenvs bool) (map[string]string, error) {
	if noenvs {
		os.Clearenv()
	}
	if envfile != "" {
		err := godotenv.Load(envfile)
		if err != nil {
			return nil, err
		}
	}

	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		if prefix != "" {
			if strings.HasPrefix(key, prefix) {
				if removePrefix {
					key = strings.Replace(key, prefix, "", 1)
				}
				envs[key] = os.Getenv(pair[0])

			}
		} else {
			envs[key] = os.Getenv(pair[0])
		}
	}

	return envs, nil
}

func generateConfig(config map[string]string) (string, error) {
	res, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("window.__RUNTIME_CONFIG__ = %s", res), nil
}

func writeFile(name string, contents string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(contents)
	return nil
}

func main() {

	app := &cli.App{
		Name:  "runtime-env",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env-file",
				Usage:   "The .env file to be parsed",
				Aliases: []string{"f"},
			},
			&cli.StringFlag{
				Name:    "prefix",
				Usage:   "The env prefix to matched",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "output",
				Usage:   "Output file path",
				Value:   "./env.js",
				Aliases: []string{"o"},
			},
			&cli.BoolFlag{
				Name:  "remove-prefix",
				Value: false,
				Usage: "Remove the prefix from the env",
			}, &cli.BoolFlag{
				Name:  "no-envs",
				Value: false,
				Usage: "Only read envs from file not from environment variables",
			},
		},
		Action: func(c *cli.Context) error {
			envs, err := load(c.String("env-file"), c.String("prefix"), c.Bool("remove-prefix"), c.Bool("no-envs"))
			if err != nil {
				return err
			}
			res, err := generateConfig(envs)
			if err != nil {
				return err
			}

			err = writeFile(c.String("output"), res)
			if err != nil {
				return err
			}
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
