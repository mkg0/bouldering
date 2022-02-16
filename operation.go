package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/golang-module/carbon/v2"
)

var qs = []*survey.Question{
	{
		Name:      "Name",
		Prompt:    &survey.Input{Message: "Firstname?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "Surname",
		Prompt:    &survey.Input{Message: "Surname?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:     "DateOfBirth",
		Prompt:   &survey.Input{Message: "Date of birth?(1989-06-06)"},
		Validate: survey.Required,
	},
	{
		Name:     "Address",
		Prompt:   &survey.Input{Message: "Address?"},
		Validate: survey.Required,
	},
	{
		Name:     "PostCode",
		Prompt:   &survey.Input{Message: "Postcode?"},
		Validate: survey.Required,
	},
	{
		Name:      "City",
		Prompt:    &survey.Input{Message: "City?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:     "Phone",
		Prompt:   &survey.Input{Message: "Phone?"},
		Validate: survey.Required,
	},
	{
		Name:      "Email",
		Prompt:    &survey.Input{Message: "Email?"},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name:     "USCid",
		Prompt:   &survey.Input{Message: "Urban Sports Club membership id?"},
		Validate: survey.Required,
	},
}

func fetch(url string, response interface{}) {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode((&response))
}

func getLongestArr(thing [][]Slot) []Slot {
	var l []Slot
	for _, item := range thing {
		if len(item) > len(l) {
			l = item
		}
	}
	return l
}

func chooseGym() BoulderStudio {
	var options []string
	for _, studio := range boulderStudios {
		options = append(options, studio.name)
	}
	prompt := &survey.Select{
		Message: "Choose a gym:",
		Options: options,
	}
	studioName := ""
	survey.AskOne(prompt, &studioName)
	if studioName == Bouldergarten.name { // lazyyyyy
		return Bouldergarten
	}
	return Boulderklub
}

func chooseProfile(message string) Profile {
	runCommand("clear")
	if len(global.Profiles) == 0 {
		return Profile{}
	}
	var options []string
	for _, p := range global.Profiles {
		options = append(options, fmt.Sprintf("%s (%s)", p.Name, p.Email))
	}
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	profileLabel := ""
	survey.AskOne(prompt, &profileLabel)
	var selectedProfile Profile
	for _, item := range global.Profiles {
		if profileLabel == fmt.Sprintf("%s (%s)", item.Name, item.Email) {
			selectedProfile = item
		}
	}
	return selectedProfile
}

func confirm(slots []Slot, gym BoulderStudio, profile Profile, message string) bool {
	runCommand("clear")
	fmt.Println("==================================================")
	fmt.Printf("   Gym: %s\n", gym.name)
	fmt.Printf("   Profile: %s(%s)\n", profile.Name, profile.Email)
	if len(slots) > 1 {
		for i, slot := range slots {
			duration := getCarbonFromUnix(slot.DateList[0].Start).DiffInMinutesWithAbs(getCarbonFromUnix(slot.DateList[0].End))
			fmt.Printf("   Slot%v: %s(%s) ⏱  %v mins\n", i+1, getCarbonFromUnix(slots[0].DateList[0].Start).ToDateTimeString(carbon.Berlin), slot.State, duration)
		}
	} else {
		fmt.Printf("   State: %s\n", slots[0].State)
		fmt.Printf("   Time: %s\n", getCarbonFromUnix(slots[0].DateList[0].Start).ToDateTimeString(carbon.Berlin))
		duration := getCarbonFromUnix(slots[0].DateList[0].Start).DiffInMinutesWithAbs(getCarbonFromUnix(slots[0].DateList[0].End))
		fmt.Printf("   ⏱  %v mins\n", duration)
	}
	fmt.Println("==================================================")
	result := false
	prompt := &survey.Confirm{
		Message: message,
		Default: true,
	}
	survey.AskOne(prompt, &result)
	return result
}

func runCommand(command string) {
	c := exec.Command(command)
	c.Stdout = os.Stdout
	c.Run()

}

func notify(slots []Slot, gym BoulderStudio) {
	runCommand("clear")
	time := getCarbonFromUnix(slots[0].DateList[0].Start).Format("l(M j) H:i", carbon.Berlin) + "-" + getCarbonFromUnix(slots[0].DateList[0].End).Format("H:i", carbon.Berlin)
	successOutput.Printf("Booked %s on %s\n", gym.name, time)
	color.New(color.FgGreen).Println("Enjoy your climbing...")
}

func (gym *BoulderStudio) bookSingle(profile Profile, slot Slot) bool {
	var shiftSelector = fmt.Sprintf(`[[0,null,true],%v,%v,%v,%v,%v,%v,%v,"%s"]`, gym.shiftModelId, slot.Selector[2].(float64), slot.Selector[3].(float64), slot.Selector[4].(float64), slot.Selector[5].(float64), int(slot.Selector[6].(float64)), int(slot.Selector[7].(float64)), slot.Selector[8].(string)) // fuck it
	var jsonStr = fmt.Sprintf(`{"clientId":%[1]v,"shiftModelId":%[2]v,"shiftSelector":%[12]s,"desiredDate":null,"dateOfBirthString":"%[3]s","streetAndHouseNumber":"%[4]s","postalCode":"%[5]s","city":"%[6]s","phoneMobile":"%[7]s","type":"booking","participants":[{"isBookingPerson":true,"tariffId":%[13]v,"dateOfBirthString":"%[3]v","firstName":"%[8]s","lastName":"%[9]s","email":"%[10]s","additionalFieldValue":"%[11]s","dateOfBirth":"%[3]v"}],"firstName":"%[8]s","lastName":"%[9]s","email":"%[10]s","dateOfBirth":"%[3]v","wantsNewsletter":false}`, gym.clientId, gym.shiftModelId, profile.DateOfBirth, profile.Address, profile.PostCode, profile.City, profile.Phone, profile.Name, profile.Surname, profile.Email, profile.USCid, shiftSelector, gym.tariffId)
	fmt.Println(jsonStr)
	client := http.Client{}
	req, _ := http.NewRequest("POST", "https://backend.dr-plano.com/bookable", strings.NewReader(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	var res map[string]interface{}
	json.NewDecoder(response.Body).Decode((&res))
	return response.StatusCode == 200
}

func autoBook(start, end carbon.Carbon, gym BoulderStudio, profile Profile, selectedSlots []Slot, retryIntervalSeconds int) bool {
	booked := false
	times := 0
	for booked != true {
		recieved := gym.getSlots(start, end)
		for _, selected := range selectedSlots {
			for _, item := range recieved {
				if item.DateList[0].Start == selected.DateList[0].Start && item.DateList[0].End == selected.DateList[0].End && item.State == "BOOKABLE" {
					result := gym.bookSingle(profile, item)
					if result {
						notify([]Slot{item}, gym)
						booked = true
						return true
					} else {
						fmt.Println("Booking has failed...")
					}
				}
			}
		}
		runCommand("clear")
		times++
		fmt.Printf("Couldn't find any available slot. Trying in %v seconds...\n Total attempts: %v\n", retryIntervalSeconds, times)
		time.Sleep(time.Second * time.Duration(retryIntervalSeconds))
	}
	return false
}

func (gym *BoulderStudio) getSlots(start, end carbon.Carbon) []Slot {
	url := fmt.Sprintf("https://backend.dr-plano.com/courses_dates?id=%v&start=%v&end=%v", gym.shiftModelId, strconv.FormatInt(start.Carbon2Time().UnixMilli(), 10), strconv.FormatInt(end.Carbon2Time().UnixMilli(), 10))
	fmt.Println(url)
	response := []Slot{}
	fetch(url, &response)
	return response
}

func ensureProfile() {
	if len(global.Profiles) == 0 {
		errorOutput.Println(`There isn't any profile yet. Add one with "bouldering profile add"`)
	}
}
