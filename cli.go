package main

import (
	"bouldering-auto-book/internal/persist"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/golang-module/carbon/v2"
	"github.com/urfave/cli/v2"
)

var errorOutput = color.New(color.FgRed)
var successOutput = color.New(color.FgMagenta)

var start = carbon.Now().StartOfDay()
var end = carbon.Now().AddDays(dayCount).EndOfDay()

func runCli() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "book",
				Usage: "Book a single slot from a gym",
				Action: func(c *cli.Context) error {
					if len(global.Profiles) == 0 {
						errorOutput.Println(`There isn't any profile yet. Add one with "bouldering profile add"`)
						return nil
					}
					gym := chooseGym()

					fmt.Println("Fetching slots...")
					slots := gym.getSlots(start, end)
					slotsToBook := askSlot(slots, start, end, false)
					if len(slotsToBook) == 0 {
						runCommand("clear")
						errorOutput.Println("You should choose a slot to book")
						return nil
					}
					profile := chooseProfile("Choose a profile:")
					confirmation := confirm(slotsToBook, gym, profile, "Would you like to book above slot(s)?")
					if confirmation == false {
						errorOutput.Println("Oh no...")
						return nil
					}
					result := gym.bookSingle(profile, slotsToBook[0])
					if result {
						notify(slotsToBook, gym)
					} else {
						errorOutput.Println("Booking has failed...")
					}
					return nil
				},
			},
			{
				Name:  "auto-book",
				Usage: "Auto book first slot that's available from selecteds",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "interval",
						Value: "10",
						Usage: "Interval for attemting another booking",
					},
				},
				Action: func(c *cli.Context) error {
					interval := defaultRetryIntervalSeconds
					if c.String("interval") != "" {
						interval_, err := strconv.Atoi(c.String("interval"))
						if err != nil {
							errorOutput.Println("Interval should be a number that stands for seconds")
							return nil
						}
						interval = interval_
					}

					if len(global.Profiles) == 0 {
						errorOutput.Println(`There isn't any profile yet. Add one with "bouldering profile add"`)
						return nil
					}
					gym := chooseGym()
					fmt.Println("Fetching slots...")
					slots := gym.getSlots(start, end)
					slotsToBook := askSlot(slots, start, end, true)
					if len(slotsToBook) == 0 {
						runCommand("clear")
						errorOutput.Println("You should choose a slot to book")
						return nil
					}
					fmt.Println("choose profile")
					profile := chooseProfile("Choose a profile:")
					fmt.Println(profile)

					successOutput.Printf("I will try to book every 30 seconds...")
					confirmation := confirm(slotsToBook, gym, profile, "Are these fields correct?")
					if confirmation == false {
						errorOutput.Println("Oh no...")
						return nil
					}

					autoBook(start, end, gym, profile, slotsToBook, interval)
					return nil
				},
			},
			{
				Name: "profile",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new profile to book automatically",
						Action: func(c *cli.Context) error {
							answers := Profile{}
							err := survey.Ask(qs, &answers)
							if err != nil {
								fmt.Println(err.Error())
								return nil
							}
							global.Profiles = append(global.Profiles, answers)
							persist.Save(&global)
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove an existing profile",
						Action: func(c *cli.Context) error {
							email := ""
							var options []string
							if len(global.Profiles) == 0 {
								errorOutput.Println(`There isn't any profile yet. Add one with "bouldering profile add"`)
								return nil
							}
							for _, profile := range global.Profiles {
								options = append(options, profile.Email)
							}
							prompt := &survey.Select{
								Message: "Choose a profile to remove:",
								Options: options,
							}
							survey.AskOne(prompt, &email)
							return nil
						},
					},
					{
						Name:  "list",
						Usage: "list the existing profiles",
						Action: func(c *cli.Context) error {
							if len(global.Profiles) == 0 {
								fmt.Println("No profile yet. Add one with 'bouldering profile add'")
								return nil
							}
							for _, profile := range global.Profiles {
								fmt.Printf("%s (%s)\n", profile.Name, profile.Email)
							}
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
