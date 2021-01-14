package main

import (
	"canwegoyet/alexa"
	"canwegoyet/besttime"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

var (
	client                 *http.Client
	apiKeyPrivate          string
	apiKeyPublic           string
	bestTimeNewForecastUrl = "https://BestTime.app/api/v1/forecasts?"
	bestTimeQuietHoursUrl  = "https://BestTime.app/api/v1/forecasts/quiet?"
)

func Handler(request alexa.Request) (alexa.Response, error) {
	return IntentDispatcher(request), nil
}

func init() {
	client = &http.Client{}
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	apiKeyPrivate = os.Getenv("BESTTIMEPRIVATEKEY")
	apiKeyPublic = os.Getenv("BESTTIMEPUBLICKEY")
}

func main() {
	lambda.Start(Handler)
}

func IntentDispatcher(request alexa.Request) alexa.Response {
	var response alexa.Response
	switch request.Body.Intent.Name {
	case "CurrentTimeIntent":
		response = HandleCurrentTimeIntent(request)
	case "QuietHoursIntent":
		response = HandleQuietHoursIntent(request)
	case alexa.HelpIntent:
		response = HandleHelpIntent(request)
	case "AboutIntent":
		response = HandleAboutIntent(request)
	default:
		response = HandleAboutIntent(request)
	}
	return response
}

func HandleCurrentTimeIntent(request alexa.Request) alexa.Response {
	var builder alexa.SSMLBuilder
	weekday := int(time.Now().Weekday()) - 1
	venueNameSlot := request.Body.Intent.Slots["VenueName"]
	venueAddressSlot := request.Body.Intent.Slots["VenueAddress"]
	response, err := CurrentDensity(venueNameSlot.Value, venueAddressSlot.Value, weekday)
	if err != nil {
		response = "Error retrieving the current population density at " + venueNameSlot.Value
	}
	builder.Say(response)
	return alexa.NewSSMLResponse("Current population density at a venue", builder.Build())
}

func CurrentDensity(venueName string, venueAddress string, weekday int) (string, error) {
	currHour := time.Now().Hour()
	// fmt.Println(currHour)
	params := url.Values{
		"venue_name":      {venueName},
		"venue_address":   {venueAddress},
		"api_key_private": {apiKeyPrivate},
	}
	// fmt.Println("api_key_private = " + apiKeyPrivate)
	url := bestTimeNewForecastUrl + params.Encode()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var forecastResponse besttime.NewForecastResponse
	err = json.Unmarshal(body, &forecastResponse)
	// fmt.Println(forecastResponse)
	responseString := "The current population density at " + venueName + " is " + forecastResponse.Analysis[weekday].HourAnalysis[currHour].Intensity
	return responseString, nil
}

func GetVenueId(venueName string, venueAddress string) (string, error) {
	params := url.Values{
		"venue_name":      {venueName},
		"venue_address":   {venueAddress},
		"api_key_private": {apiKeyPrivate},
	}
	// fmt.Println("api_key_private = " + apiKeyPrivate)
	url := bestTimeNewForecastUrl + params.Encode()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var forecastResponse besttime.NewForecastResponse
	err = json.Unmarshal(body, &forecastResponse)
	// fmt.Println(forecastResponse)
	return forecastResponse.VenueInfo.VenueId, nil
}

func QuietHours(venueId string) (string, error) {
	params := url.Values{
		"venue_id":       {venueId},
		"api_key_public": {apiKeyPublic},
	}
	// fmt.Println("api_key_private = " + apiKeyPrivate)
	url := bestTimeQuietHoursUrl + params.Encode()
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	var quietHoursResponse besttime.QuietHoursResponse
	err = json.Unmarshal(body, &quietHoursResponse)
	if err != nil {
		fmt.Println("Error unmarshaling", err)
		return "", err
	}
	fmt.Println(quietHoursResponse)
	responseString := "The quiet hours are "
	for _, hour := range quietHoursResponse.Analysis.QuietHoursList[:len(quietHoursResponse.Analysis.QuietHoursList)-1] {
		responseString += hour + ", "
	}
	responseString += "and " + quietHoursResponse.Analysis.QuietHoursList[len(quietHoursResponse.Analysis.QuietHoursList)-1]
	return responseString, nil
}

func HandleQuietHoursIntent(request alexa.Request) alexa.Response {
	var builder alexa.SSMLBuilder
	var response string
	venueNameSlot := request.Body.Intent.Slots["VenueName"]
	venueAddressSlot := request.Body.Intent.Slots["VenueAddress"]
	venueId, err := GetVenueId(venueNameSlot.Value, venueAddressSlot.Value)
	if err != nil {
		fmt.Println("Error retrieving venueId", err)
		response = "Error retrieving the quiet hours at " + venueNameSlot.Value
	}
	fmt.Println("venueId = " + venueId)
	response, err = QuietHours(venueId)
	if err != nil {
		fmt.Println("Error retrieving quiet hours", err)
		response = "Error retrieving the quiet hours at " + venueNameSlot.Value
	}
	builder.Say(response)
	return alexa.NewSSMLResponse("Quiet hours at a venue", builder.Build())
}

func HandleHelpIntent(request alexa.Request) alexa.Response {
	var builder alexa.SSMLBuilder
	builder.Say("Here are some of the things you can ask:")
	builder.Pause("1000")
	builder.Say("How is it like at a place at this address?")
	builder.Pause("1000")
	builder.Say("What are the quiet hours at a place at this address")
	return alexa.NewSSMLResponse("CanWeGoYet Help", builder.Build())
}

func HandleAboutIntent(request alexa.Request) alexa.Response {
	return alexa.NewSimpleResponse("About", "CanWeGoYet is an Alexa Skill that can give an estimate of the current population density and the quiet hours of a place.")
}
