package stats

import (
	"fmt"
	"sort"
	"strings"

	recipeDomain "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/model"
)

// Service defines contracts for recipe service.
type Stats interface {
	CalculateStats(ch <-chan recipeDomain.Data, filterByPostcodeAndTime recipeDomain.DeliveriesByPostcodeAndTime, filterByWords []string) recipeDomain.ExpectedOutput
}

// service implements Service interface.
type stats struct{}

// NewService returns a new instance of service.
func NewStatsCalculator() Stats {
	return &stats{}
}

// CalculateStats will calculate and return stats as ExpectedOutput.
func (s *stats) CalculateStats(ch <-chan recipeDomain.Data, filterByPostcodeAndTime recipeDomain.DeliveriesByPostcodeAndTime, filterByWords []string) recipeDomain.ExpectedOutput {
	var uniqueRecipeCountMap = make(map[string]uint)
	var countPerPostcodeMap = make(map[string]uint)
	var countDeliveriesPerPostcode = make(map[string]uint)
	var filteredRecipeNames []string

	for data := range ch {

		// populate uniqueRecipeCountMap
		if _, ok := uniqueRecipeCountMap[data.Recipe]; !ok {
			uniqueRecipeCountMap[data.Recipe] = 1
		} else {
			uniqueRecipeCountMap[data.Recipe] += 1
		}

		// populate countPerPostcodeMap
		if _, ok := countPerPostcodeMap[data.Postcode]; !ok {
			countPerPostcodeMap[data.Postcode] = 1
		} else {
			countPerPostcodeMap[data.Postcode] += 1
		}

		// populate countDeliveriesPerPostcode
		if data.Postcode == filterByPostcodeAndTime.Postcode && data.InDeliveryTime(filterByPostcodeAndTime) {
			countDeliveriesPerPostcode[data.Postcode] += 1
		}

		// populate filteredRecipeNames
		for _, word := range filterByWords {
			if strings.Contains(data.Recipe, word) && !func() bool {
				for _, filteredRecipeName := range filteredRecipeNames {
					if filteredRecipeName == data.Recipe {
						return true
					}
				}
				return false
			}() {
				filteredRecipeNames = append(filteredRecipeNames, data.Recipe)
				break
			}
		}
	}
	return s.buildExpectedOutput(uniqueRecipeCountMap, countPerPostcodeMap, countDeliveriesPerPostcode, filterByPostcodeAndTime, filteredRecipeNames)
}

// buildExpectedOutput will build final ExpectedOutput.
func (s *stats) buildExpectedOutput(uniqueRecipeCountMap, countPerPostcodeMap, countDeliveriesPerPostcode map[string]uint,
	filterByPostcodeAndTime recipeDomain.DeliveriesByPostcodeAndTime, filteredRecipeNames []string) recipeDomain.ExpectedOutput {

	var out recipeDomain.ExpectedOutput

	// set UniqueRecipeCount
	out.UniqueRecipeCount = s.createUniqueRecipeCount(uniqueRecipeCountMap)
	// set CountPerRecipe
	out.CountPerRecipe = s.createSortedCountPerRecipe(uniqueRecipeCountMap)
	// set BusiestPostcode
	out.BusiestPostcode = s.createBusiestPostcode(countPerPostcodeMap)
	// set CountPerPostcodeAndTime
	out.CountPerPostcodeAndTime = s.createCountPerPostcodeAndTime(countDeliveriesPerPostcode, filterByPostcodeAndTime)
	// set MatchByName
	sort.Strings(filteredRecipeNames)
	out.MatchByName = filteredRecipeNames

	return out
}

// createUniqueRecipeCount will create UniqueRecipeCount from given uniqueRecipeCountMap.
func (s *stats) createUniqueRecipeCount(uniqueRecipeCountMap map[string]uint) uint {
	var out uint

	for _, count := range uniqueRecipeCountMap {
		if count >= 1 {
			out += 1
		}
	}

	return out
}

// createSortedCountPerRecipe will create sorted CountPerRecipe from given uniqueRecipeCountMap.
func (s *stats) createSortedCountPerRecipe(uniqueRecipeCountMap map[string]uint) []recipeDomain.CountPerRecipe {
	// create list of recipe names
	var recipeList = make([]string, 0, len(uniqueRecipeCountMap))
	for recipe := range uniqueRecipeCountMap {
		recipeList = append(recipeList, recipe)
	}

	// sort recipe names
	sort.Strings(recipeList)

	var out []recipeDomain.CountPerRecipe

	// create and append CountPerRecipe to out
	for _, recipe := range recipeList {
		out = append(out, recipeDomain.CountPerRecipe{
			Recipe: recipe,
			Count:  uniqueRecipeCountMap[recipe],
		})
	}

	return out
}

// createBusiestPostcode will create BusiestPostcode from given countPerPostcodeMap
// the BusiestPostcode is the one with most delivered recipes.
func (s *stats) createBusiestPostcode(countPerPostcodeMap map[string]uint) recipeDomain.BusiestPostcode {
	// create list of BusiestPostcode
	var postcodeList []recipeDomain.BusiestPostcode
	for postcode, count := range countPerPostcodeMap {
		postcodeList = append(postcodeList, recipeDomain.BusiestPostcode{
			Postcode:      postcode,
			DeliveryCount: count,
		})
	}

	// sort BusiestPostcode in postcodeList based on DeliveryCount
	sort.Slice(postcodeList, func(i, j int) bool {
		return postcodeList[i].DeliveryCount > postcodeList[j].DeliveryCount
	})
	if len(postcodeList) == 0 {
		return recipeDomain.BusiestPostcode{}
	}
	return postcodeList[0]
}

// createCountPerPostcodeAndTime will create CountPerPostcodeAndTime from given countDeliveriesPerPostcode.
func (s *stats) createCountPerPostcodeAndTime(countDeliveriesPerPostcode map[string]uint, filterByPostcodeAndTime recipeDomain.DeliveriesByPostcodeAndTime) recipeDomain.CountPerPostcodeAndTime {
	return recipeDomain.CountPerPostcodeAndTime{
		Postcode:      filterByPostcodeAndTime.Postcode,
		FromAM:        fmt.Sprintf("%vAM", filterByPostcodeAndTime.From),
		ToPM:          fmt.Sprintf("%vPM", filterByPostcodeAndTime.To),
		DeliveryCount: countDeliveriesPerPostcode[filterByPostcodeAndTime.Postcode],
	}
}
