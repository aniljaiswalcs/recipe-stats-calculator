package stats_test

import (
	"testing"

	recipeData "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/model"
)

// TestData_InDeliveryTime is a unit test for recipe -> InDeliveryTime() method.
func TestData_InDeliveryTime(t *testing.T) {
	tests := []struct {
		name string
		data *recipeData.Data
		arg  recipeData.DeliveriesByPostcodeAndTime
		want bool
	}{
		{
			name: "Happy path",
			data: &recipeData.Data{
				Postcode: "10120",
				Recipe:   "Creamy Dill Chicken",
				Delivery: "Thursday 11AM - 2PM",
			},
			arg: recipeData.DeliveriesByPostcodeAndTime{
				Postcode: "10120",
				From:     10,
				To:       3,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.InDeliveryTime(tt.arg); got != tt.want {
				t.Errorf("Data.InDeliveryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePostcodeLen(t *testing.T) {
	// Test cases for ValidatePostcodeLen
	testCases := []struct {
		postcode   string
		expectPass bool
	}{
		{"12345", true},        // Valid postcode
		{"1234567890", true},   // Valid postcode (exactly 10 characters)
		{"12345678901", false}, // Invalid postcode (more than 10 characters)
		{"", false},            // Valid postcode (empty string)
		{"test", false},        // check only digit
	}

	for _, tc := range testCases {
		data := recipeData.Data{Postcode: tc.postcode}
		result := data.ValidatePostcodeLen()
		if result != tc.expectPass {
			t.Errorf("Expected ValidatePostcodeLen() to return %v for postcode '%s', but got %v", tc.expectPass, tc.postcode, result)
		}
	}
}

func TestValidateRecipeNameLen(t *testing.T) {
	// Test cases for ValidateRecipeNameLen
	testCases := []struct {
		recipeName string
		expectPass bool
	}{
		{"Delicious Recipe", true},           // Valid recipe name
		{"Lorem ipsum dolor sit amet", true}, // Valid recipe name (exactly 100 characters)
		{"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", false}, // Invalid recipe name (more than 100 characters)
		{"", false}, // Valid recipe name (empty string)
	}

	for _, tc := range testCases {
		data := recipeData.Data{Recipe: tc.recipeName}
		result := data.ValidateRecipeNameLen()
		if result != tc.expectPass {
			t.Errorf("Expected ValidateRecipeNameLen() to return %v for recipe name '%s', but got %v", tc.expectPass, tc.recipeName, result)
		}
	}
}

func Test_ValidateDeliverydata(t *testing.T) {
	// Test cases for Test_ValidateDeliverydata
	testCases := []struct {
		delivery   string
		expectPass bool
	}{
		{"Thursday 11AM - 2PM", true},      // Valid delivery
		{"Friday 1AM - 1PM", true},         // Valid delivery
		{"Thursday 11AM - 2", false},       // INValid delivery end time not correct
		{"Thursday 11AM - 2333PM.", false}, // INValid delivery end time not correct
		{"Thursday 11AM - ##PM.", false},   // INValid delivery end time not correct
		{"Thursday 11AM - 13PM.", false},   // INValid delivery end time not correct
		{"Thursday 11AM - RC PM.", false},  // INValid delivery end time not correct
		{"Thursday ##1AM - 11PM.", false},  // INValid delivery start time not correct
		{"Thursday 1222AM - 11PM.", false}, // INValid delivery start time not correct
		{"Thursday 14AM - 11PM.", false},   // INValid delivery start time not correct
		{"Thursday ST AM - 11PM.", false},  // INValid delivery start time not correct
		{"GJJJJday 11 AM - 11PM.", false},  // INValid delivery day not correct
		{"@@@@@ 11 AM - 11PM.", false},     // INValid delivery day not correct
		{"  11 AM - 11PM.", false},         // INValid delivery day not correct
		{"", false},                        // INValid delivery name (empty string)
	}

	for _, tc := range testCases {
		data := recipeData.Data{Delivery: tc.delivery}
		result := data.ValidateDelivery()
		if result != tc.expectPass {
			t.Errorf("Expected ValidateDelivery() to return %v for recipe name '%s', but got %v", tc.expectPass, tc.delivery, result)
		}
	}
}
