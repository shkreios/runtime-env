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
	"github.com/shkreios/runtime-env/pkg/config"
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

func generateJSConfig(config map[string]string) (string, error) {
	res, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("window.__RUNTIME_CONFIG__ = %s", res), nil
}

func KeysString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, fmt.Sprintf("\t\t\t%s: string;", k))
	}
	return strings.Join(keys, "\n")
}

func generateTSConfig(config map[string]string) (string, error) {
	return fmt.Sprintf(`/* eslint-disable */
/* ignore jslint start */
// tslint:disable
// jscs:disable
// jshint ignore: start
// prettier-ignore

export {};

declare global {
	interface Window {
		__RUNTIME_CONFIG__: {
%s		
		};
	}
}`, KeysString(config)), nil
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

	var envFile string
	var prefix string
	var output string
	var typeDeclarationsFile string
	var removePrefix bool
	var noEnvs bool

	app := &cli.App{
		Version:              config.Version,
		EnableBashCompletion: true,
		Name:                 "runtime-env",
		Usage:                "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env-file",
				Destination: &envFile,
				Usage:       "The .env file to be parsed",
				Aliases:     []string{"f"},
			},
			&cli.StringFlag{
				Name:        "prefix",
				Destination: &prefix,
				Usage:       "The env prefix to matched",
				Aliases:     []string{"p"},
			},
			&cli.StringFlag{
				Name:        "output",
				Destination: &output,
				Usage:       "Output file path",
				Value:       "./env.js",
				Aliases:     []string{"o"},
			},
			&cli.StringFlag{
				Name:        "type-declarations-file",
				Destination: &typeDeclarationsFile,
				Usage:       "Output file path for the typescript declaration file",
				Aliases:     []string{"dts"},
			},
			&cli.BoolFlag{
				Name:        "remove-prefix",
				Destination: &removePrefix,
				Value:       false,
				Usage:       "Remove the prefix from the env",
			}, &cli.BoolFlag{
				Name:        "no-envs",
				Destination: &noEnvs,
				Value:       false,
				Usage:       "Only read envs from file not from environment variables",
			},
		},
		Action: func(c *cli.Context) error {
			envs, err := load(envFile, prefix, removePrefix, noEnvs)
			if err != nil {
				return err
			}
			js, err := generateJSConfig(envs)
			if err != nil {
				return err
			}

			if typeDeclarationsFile != "" {
				ts, err := generateTSConfig(envs)
				if err != nil {
					return err
				}
				err = writeFile(typeDeclarationsFile, ts)
				if err != nil {
					return err
				}
			}

			err = writeFile(output, js)
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
