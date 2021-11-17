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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version              = "v1.2.0"
	nolog                bool
	envFile              string
	prefix               string
	output               string
	typeDeclarationsFile string
	removePrefix         bool
	noEnvs               bool
	globalKey            string
	watch                bool
)

func printf(format string, a ...interface{}) (n int, err error) {
	if nolog {
		return 0, nil
	}
	return fmt.Printf(format, a...)
}

func load() (map[string]string, error) {
	if noEnvs {
		os.Clearenv()
	}
	if envFile != "" {
		err := godotenv.Load(envFile)
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

func generateJSConfig(config map[string]string, globalKey string) (string, error) {
	res, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("window.%s = %s", globalKey, res), nil
}

func keysString(m map[string]string, template string, delimiter string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, fmt.Sprintf(template, k))
	}
	sort.Strings(keys)
	return strings.Join(keys, delimiter)
}

func generateTSConfig(config map[string]string) (string, error) {
	return fmt.Sprintf(`/* eslint-disable */
/* ignore jslint start */
// tslint:disable
// jscs:disable
// jshint ignore: start 

// prettier-ignore
export {};

// prettier-ignore
declare global {
	interface Window {
		__RUNTIME_CONFIG__: {
%s		
		};
	}
}`, keysString(config, "\t\t\t%s: string;", "\n")), nil
}

func writeFile(name string, contents string) error {
	path := filepath.Dir(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		printf("Non existent folders on path '%s' have been created\n", path)
		if err != nil {
			return err
		}

	}
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(contents)
	return nil
}

func run() error {
	envs, err := load()
	if err != nil {
		return err
	}

	printf("Following envs have been loaded: %s\n", keysString(envs, "%s", ", "))

	js, err := generateJSConfig(envs, globalKey)
	if err != nil {
		return err
	}
	err = writeFile(output, js)
	if err != nil {
		return err
	}
	printf("ENVs have been writtem to %s\n", output)

	if typeDeclarationsFile != "" {
		ts, err := generateTSConfig(envs)
		if err != nil {
			return err
		}
		err = writeFile(typeDeclarationsFile, ts)
		if err != nil {
			return err
		}
		printf("Typescript declarations have been writtem to %s\n", typeDeclarationsFile)
	}

	return err
}

func main() {

	app := &cli.App{
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "Simon Hessel",
				Email: "simon.hessel@kreios.lu",
			},
		},
		UsageText:              "runtime-env [global options]",
		Copyright:              "Copyright © 2020 Simon Hessel",
		EnableBashCompletion:   true,
		Name:                   "runtime-env",
		Usage:                  "runtime envs for SPAs",
		UseShortOptionHandling: true,

		HideHelpCommand: true,
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
			&cli.StringFlag{
				Name:        "global-key",
				Destination: &globalKey,
				Usage:       "Customize the key on which the envs will be set on window object",
				Aliases:     []string{"key"},
				Value:       "__RUNTIME_CONFIG__",
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
			}, &cli.BoolFlag{
				Name:        "disable-logs",
				Destination: &nolog,
				Value:       false,
				Aliases:     []string{"no-logs"},
				Usage:       "Disable logging output",
			}, &cli.BoolFlag{
				Name:        "watch",
				Destination: &watch,
				Value:       false,
				Aliases:     []string{"w"},
				Usage:       "Watch .env file",
			},
		},
		Action: func(c *cli.Context) error {
			if watch && envFile == "" {
				return errors.New("the watch falg can only be used with a .env file")
			}

			err := run()
			if watch {
				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					return err
				}
				defer watcher.Close()

				done := make(chan bool)
				go func() {
					for {
						select {
						case event, ok := <-watcher.Events:
							if !ok {
								return
							}
							if event.Op&fsnotify.Write == fsnotify.Write {
								run()
							}
						case err, ok := <-watcher.Errors:
							if !ok {
								return
							}
							log.Fatal(err)
						}
					}
				}()

				err = watcher.Add(envFile)
				if err != nil {
					return err
				}
				<-done
			}
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
