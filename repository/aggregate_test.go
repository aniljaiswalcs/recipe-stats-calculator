package stats_test

import (
	"testing"

	recipeDomain "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/model"
	recipeStats "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/repository"
	"github.com/stretchr/testify/assert"
)

var emptySlice []string

func TestCalculateStats(t *testing.T) {
	cases := []struct {
		description         string
		data                []recipeDomain.Data
		deliverPostcodeTIme recipeDomain.DeliveriesByPostcodeAndTime
		wordfilter          []string
		output              recipeDomain.ExpectedOutput
	}{
		{
			description: "Positive test case",
			data: []recipeDomain.Data{
				{Postcode: "10115", Recipe: "Potato Soup", Delivery: "Wednesday 1AM - 7PM"},
				{Postcode: "10115", Recipe: "Mushroom Risotto", Delivery: "Wednesday 1AM - 7PM"},
				{Postcode: "10115", Recipe: "Veggie Pizza", Delivery: "Wednesday 5AM - 8PM"},
				{Postcode: "10115", Recipe: "Potato Soup", Delivery: "Wednesday 3AM - 10PM"},
				{Postcode: "10115", Recipe: "Mushroom Risotto", Delivery: "Wednesday 1AM - 7PM"},
				{Postcode: "10115", Recipe: "Veggie Pizza", Delivery: "Wednesday 10AM - 12PM"},
				{Postcode: "10243", Recipe: "Mushroom Risotto", Delivery: "Wednesday 10AM - 4PM"},
				{Postcode: "10243", Recipe: "Potato Soup", Delivery: "Wednesday 1AM - 1PM"},
				{Postcode: "10243", Recipe: "Mushroom Risotto", Delivery: "Wednesday 1AM - 7PM"},
				{Postcode: "10243", Recipe: "Veggie Pizza", Delivery: "Wednesday 1AM - 7PM"},
			},
			deliverPostcodeTIme: recipeDomain.DeliveriesByPostcodeAndTime{
				Postcode: "10115",
				From:     10,
				To:       12,
			},
			wordfilter: []string{"Potato", "Veggie", "Mushroom"},
			output: recipeDomain.ExpectedOutput{
				UniqueRecipeCount: 3,
				CountPerRecipe: []recipeDomain.CountPerRecipe{
					{Recipe: "Mushroom Risotto", Count: 4},
					{Recipe: "Potato Soup", Count: 3},
					{Recipe: "Veggie Pizza", Count: 3},
				},
				BusiestPostcode: recipeDomain.BusiestPostcode{
					Postcode: "10115", DeliveryCount: 6},
				CountPerPostcodeAndTime: recipeDomain.CountPerPostcodeAndTime{
					Postcode:      "10115",
					FromAM:        "10AM",
					ToPM:          "12PM",
					DeliveryCount: 1,
				},
				MatchByName: []string{"Mushroom Risotto", "Potato Soup", "Veggie Pizza"},
			},
		},

		{
			description: "wordfilter with different dishes than input, result empty slice",
			data: []recipeDomain.Data{
				{Postcode: "10115", Recipe: "Potato Soup", Delivery: "Wednesday 1AM - 7PM"},
				{Postcode: "10115", Recipe: "Mushroom Risotto", Delivery: "Wednesday 1AM - 7PM"},
			},
			deliverPostcodeTIme: recipeDomain.DeliveriesByPostcodeAndTime{
				Postcode: "10115",
				From:     10,
				To:       12,
			},
			wordfilter: []string{"New dish", "Old dish"},
			output: recipeDomain.ExpectedOutput{
				UniqueRecipeCount: 2,
				CountPerRecipe: []recipeDomain.CountPerRecipe{
					{Recipe: "Mushroom Risotto", Count: 1},
					{Recipe: "Potato Soup", Count: 1},
				},
				BusiestPostcode: recipeDomain.BusiestPostcode{
					Postcode:      "10115",
					DeliveryCount: 2,
				},
				CountPerPostcodeAndTime: recipeDomain.CountPerPostcodeAndTime{
					Postcode:      "10115",
					FromAM:        "10AM",
					ToPM:          "12PM",
					DeliveryCount: 0,
				},
				MatchByName: emptySlice,
			},
		},
		{
			description: "From and TO value are out of range 0 - 12, return 0 count of CountPerPostcodeAndTime ",
			data: []recipeDomain.Data{
				{Postcode: "10115", Recipe: "Potato Soup", Delivery: "Wednesday 12AM - 12PM"},
				{Postcode: "10115", Recipe: "Mushroom Risotto", Delivery: "Wednesday 11AM - 7PM"},
			},
			deliverPostcodeTIme: recipeDomain.DeliveriesByPostcodeAndTime{
				Postcode: "10115",
				From:     15,
				To:       18,
			},
			wordfilter: []string{"Potato", "Veggie", "Mushroom"},
			output: recipeDomain.ExpectedOutput{
				UniqueRecipeCount: 2,
				CountPerRecipe: []recipeDomain.CountPerRecipe{
					{Recipe: "Mushroom Risotto", Count: 1},
					{Recipe: "Potato Soup", Count: 1},
				},
				BusiestPostcode: recipeDomain.BusiestPostcode{
					Postcode:      "10115",
					DeliveryCount: 2,
				},
				CountPerPostcodeAndTime: recipeDomain.CountPerPostcodeAndTime{
					Postcode:      "10115",
					FromAM:        "15AM",
					ToPM:          "18PM",
					DeliveryCount: 0,
				},
				MatchByName: []string{"Mushroom Risotto", "Potato Soup"},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			ch := make(chan recipeDomain.Data, len(tt.data))
			go func(ch chan recipeDomain.Data) {
				for _, d := range tt.data {
					{
						ch <- d
					}
				}
				defer close(ch)
			}(ch)

			recipeSVC := recipeStats.NewStatsCalculator()

			actualOutput := recipeSVC.CalculateStats(ch, tt.deliverPostcodeTIme, tt.wordfilter)

			assert.Equal(t, tt.output, actualOutput)

		})
	}

}
