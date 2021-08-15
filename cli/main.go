package main

import (
	"fmt"
	"github.com/phith0n/zkar"
	"github.com/phith0n/zkar/payloads"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)


func main() {
	var app = cli.App{
		Name: "zkar",
		Usage: "A Java serialization tool",
		Commands: []*cli.Command {
			{
				Name: "generate",
				Usage: "generate Java serialization attack payloads",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "output",
						Usage: "output file path",
						Aliases: []string{"o"},
						Required: false,
						Value: "",
					},
					&cli.BoolFlag{
						Name: "list",
						Usage: "list all available gadgets",
						Aliases: []string{"l"},
						Required: false,
						Value: false,
					},
				},
				Action: func(context *cli.Context) error {
					var find = context.Args().Get(0)

					// list all payloads
					if context.Bool("list") || find == "" {
						var l = funk.Map(payloads.BuiltinGadget, func(g payloads.Gadget) string {
							return g.Name()
						})
						fmt.Printf("Available gadgets: %s\n", strings.Join(l.([]string), ", "))
						return nil
					}

					for _, gadget := range payloads.BuiltinGadget {
						if gadget.Name() != find {
							continue
						}

						var args = context.Args().Slice()
						ser, err := gadget.Generate(args[1:]...)
						if err != nil {
							return fmt.Errorf("generate payload failed, error: %v", err.Error())
						}

						if context.String("output") != "" {
							ioutil.WriteFile(context.String("output"), ser.ToBytes(), fs.FileMode(644))
						} else {
							fmt.Print(string(ser.ToBytes()))
						}
						return nil
					}

					return fmt.Errorf("gadget %v is not found", find)
				},
			},
			{
				Name: "dump",
				Usage: "parse the Java serialization streams and dump the struct",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "file",
						Aliases: []string{"f"},
						Usage: "serialization data filepath",
						Required: true,
					},
					&cli.BoolFlag{
						Name: "golang",
						Usage: "dump the Go language based struct instead of human readable information",
						Required: false,
						Value: false,
					},
				},
				Action: func(context *cli.Context) error {
					var filename = context.String("file")
					data, err := ioutil.ReadFile(filename)
					if err != nil {
						return err
					}

					ser, err := zkar.FromBytes(data)
					if err != nil {
						return nil
					}

					if context.Bool("golang") {
						zkar.DumpToGoStruct(ser)
					} else {
						fmt.Println(ser.ToString())
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] %v\n", err.Error())
	}
}
