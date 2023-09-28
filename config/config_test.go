package cfg_test

import (
	"testing"

	cfg "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/config"
	"github.com/stretchr/testify/assert"
)

const testfilepath = "./../config.yml"

func TestViper(t *testing.T) {

	// Read in config file
	config := cfg.NewConfig(testfilepath)

	// Test getting a value
	expectedValue := "10120"
	actualValue := config.Postcode
	assert.Equal(t, expectedValue, actualValue)

	// Test getting slice value
	expectedReciepeValue := []string{"Potato", "Veggie", "Mushroom"}
	actualDefaultValue := config.MatchedRecipes
	for index, actualval := range actualDefaultValue {
		assert.Equal(t, expectedReciepeValue[index], actualval)
	}

	// test value which is not equal
	var expectedUintValue uint
	expectedUintValue = 10000
	actualUintValue := config.From
	assert.NotEqual(t, actualUintValue, expectedUintValue)
}
