package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mkg0/bouldering/internal/persist"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/golang-module/carbon/v2"
	"github.com/urfave/cli/v2"
)

var errorOutput = color.New(color.FgRed)
var successOutput = color.New(color.FgMagenta)

func runCli() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "book",
				Usage: "Book a single slot from a gym",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "local",
						Usage: "skip cloud booking",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "show the booking without an actual booking",
					},
					&cli.StringFlag{
						Name:  "offset",
						Value: "0",
						Usage: "Day offset for slot selection matrix",
					},
				},
				Action: func(c *cli.Context) error {
					isLocal := isLocalBooking(c.Bool("local"))
					isDryRun := isLocalBooking(c.Bool("dry-run"))

					offset := 0
					if len(global.Profiles) == 0 {
						errorOutput.Println(`There isn't any profile yet. Add one with "bouldering profile add"`)
						return nil
					}
					if c.String("offset") != "" {
						offset_, err := strconv.Atoi(c.String("offset"))
						if err != nil {
							errorOutput.Println("offset should be a number that stands for day")
							return nil
						}
						offset = offset_
					}
					gym := chooseGym()

					fmt.Println("Fetching slots...")
					var start = carbon.Now().AddDays(offset).StartOfDay()
					var end = carbon.Now().AddDays(offset).AddDays(dayCount).EndOfDay()
					slots := gym.getSlots(start, end)
					slotsToBook := askSlot(slots, start, end, false, !isLocal)
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
					if isDryRun {
						notify(slotsToBook, gym)
						errorOutput.Println("Booking call skipped due to dry run")
						return nil
					}
					result := gym.bookSingle(profile, slotsToBook[0], isLocal)
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
				Usage: "Continuously try to book the first slot that's available from selected slots(local only) ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "interval",
						Value: "10",
						Usage: "Interval for each booking attempt",
					},
				},
				Action: func(c *cli.Context) error {
					interval := defaultRetryIntervalSeconds
					offset := 0
					if c.String("interval") != "" {
						interval_, err := strconv.Atoi(c.String("interval"))
						if err != nil {
							errorOutput.Println("Interval should be a number that stands for seconds")
							return nil
						}
						interval = interval_
					}
					if c.String("offset") != "" {
						offset_, err := strconv.Atoi(c.String("offset"))
						if err != nil {
							errorOutput.Println("offset should be a number that stands for day")
							return nil
						}
						offset = offset_
					}

					if len(global.Profiles) == 0 {
						errorOutput.Println(`There isn't any profile yet. Add one with "bouldering profile add"`)
						return nil
					}
					gym := chooseGym()
					fmt.Println("Fetching slots...")
					var start = carbon.Now().AddDays(offset).StartOfDay()
					var end = carbon.Now().AddDays(offset).AddDays(dayCount).EndOfDay()
					slots := gym.getSlots(start, end)
					slotsToBook := askSlot(slots, start, end, true, false)
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
				Name:  "enable-remote-booking",
				Usage: "Enables scheduled future bookings(later than one week) and telegram bot features",
				Action: func(c *cli.Context) error {
					value := &survey.Input{
						Message: "What is the endpoint for API?",
					}
					var result string
					err := survey.AskOne(value, &result)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					global.ApiEndpoint = result
					global.RemoteBookingEnabled = true
					persist.Save(&global)
					fmt.Println(result)
					return nil
				},
			},
			{
				Name:  "disable-remote-booking",
				Usage: "Returns back to local booking",
				Action: func(c *cli.Context) error {
					global.RemoteBookingEnabled = false
					persist.Save(&global)
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
								errorOutput.Println("No profile yet. Add one with 'bouldering profile add'")
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
			{
				Name: "config",
				Subcommands: []*cli.Command{
					{
						Name:  "show",
						Usage: "Show config content and the file path",
						Action: func(c *cli.Context) error {
							fmt.Println(persist.GetFilePath())
							PrettyPrint(global)
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
