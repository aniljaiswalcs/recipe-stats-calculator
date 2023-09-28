package main

import (
	"fmt"
	"os"

	cfg "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/config"
	jsonfile "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/handler"
	recipeDomain "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/model"
	recipeStats "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/repository"
)

const configFilePath = "./config.yml"

func main() {

	// instantiate config file and file reader
	config := cfg.NewConfig(configFilePath)
	jsonReader := jsonfile.NewReader()

	// read each recipe data and send it to channel
	ch, err := jsonReader.Read(config.DataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		os.Exit(1)
	}

	// instantiate recipe service
	recipeSVC := recipeStats.NewStatsCalculator()

	// to count the number of deliveries to postcode 10120 that lie within the delivery time between 10AM and 3PM
	filterByPostcodeAndTime := recipeDomain.DeliveriesByPostcodeAndTime{
		Postcode: config.Postcode,
		From:     config.From,
		To:       config.To,
	}

	// call CalculateStats method from recipeSVC
	output := recipeSVC.CalculateStats(ch, filterByPostcodeAndTime, config.MatchedRecipes)

	err = output.GenerateOutput(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
	}

}
