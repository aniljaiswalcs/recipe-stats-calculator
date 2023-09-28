package stats

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Data represents a single recipe data obj in hf_test_calculation_fixtures.json file.
type Data struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

// InDeliveryTime will check given time in Data.Delivery
// d.Delivery has following format: {weekday} {h}AM - {h}PM.
func (d *Data) InDeliveryTime(filterByPostcodeAndTime DeliveriesByPostcodeAndTime) bool {

	regex := regexp.MustCompile(`(\d{0,2})AM\s-\s(\d{0,2})PM`)
	matches := regex.FindStringSubmatch(d.Delivery)

	from, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		return false
	}

	to, err := strconv.Atoi(matches[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		return false
	}

	return uint(from) >= filterByPostcodeAndTime.From && uint(to) <= filterByPostcodeAndTime.To
}

func (d *Data) ValidateDelivery() bool {

	deliveryDay := strings.Split(d.Delivery, " ")
	saturdayPattern := `^(Satur|Sun)day$` //only Saturday and Sunday not consider as weekday
	satmatch, _ := regexp.MatchString(saturdayPattern, deliveryDay[0])
	if satmatch {
		/* errorString := fmt.Sprintf("%q is not defined day as a Weekday", deliveryDay[0])
		fmt.Fprintf(os.Stderr, "%q\n", errorString)
		*/
		return false
	}

	pattern := `^(Mon|Tues|Wednes|Thurs|Fri)day$` //only weekday Allowed
	match, err := regexp.MatchString(pattern, deliveryDay[0])
	if !match {
		errorString := fmt.Sprintf("%s is not defined day as a Weekday, skipping parsing data: %q", deliveryDay[0], d)
		fmt.Fprintf(os.Stderr, "%s\n", errorString)
		return false
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		return false
	}

	regex := regexp.MustCompile(`\s(0?[1-9]|1[0-2])AM\s-\s(0?[1-9]|1[0-2])PM$`)
	matches := regex.FindStringSubmatch(d.Delivery)

	if len(matches) == 0 {
		errortimeString := fmt.Sprintf("time is not correctly format: %q", deliveryDay[1:])
		errString := fmt.Sprintf("Error: %s, skipping parsing data: %s", errortimeString, d)
		fmt.Fprintf(os.Stderr, "%s\n\n", errString)
		return false
	}

	_, err = strconv.Atoi(matches[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		return false
	}

	_, err = strconv.Atoi(matches[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		return false
	}

	return true
}

// IsDigitsOnly checks if a string contains only digits.
func IsDigitsOnly(input string) bool {
	for _, char := range input {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

// Rule 2: Validate the number of distinct postcodes
func (d *Data) ValidatePostcodeLen() bool {
	// Assuming that the length of postcode should not exceed 10 characters
	return len(d.Postcode) > 0 && len(d.Postcode) <= 10 && IsDigitsOnly(d.Postcode)
}

// Rule 3: Validate the number of distinct recipe names
func (d *Data) ValidateRecipeNameLen() bool {
	// Assuming that the length of recipe name should not exceed 100 characters
	return len(d.Recipe) > 0 && len(d.Recipe) <= 100
}
