package json_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	fsJson "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/handler"
	recipeData "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/model"
	"github.com/stretchr/testify/assert"
)

// Test_reader_Read is a unit test for json -> Read() method.
func Test_Read(t *testing.T) {
	file, err := ioutil.TempFile("", "test.json")

	if err != nil {
		t.Errorf("Failed to create temporary file: %v", err)
	}

	defer os.Remove(file.Name())

	// Write JSON data to the temporary file
	data := []recipeData.Data{
		{Postcode: "10120", Recipe: "Potato Curry", Delivery: "Wednesday 1AM - 7PM"},
		{Postcode: "10121", Recipe: "Veggie Pasta", Delivery: "Friday 10AM - 4PM"},
		{Postcode: "10122", Recipe: "Mushroom Risotto", Delivery: "Tuesday 2AM - 5PM"},
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal JSON data: %v", err)
	}
	err = ioutil.WriteFile(file.Name(), bytes, 0644)
	if err != nil {
		t.Errorf("Failed to write JSON data to file: %v", err)
	}

	tests := []struct {
		name string
		arg  string
		want []recipeData.Data
	}{
		{
			name: "Matching recipe data",
			arg:  file.Name(),
			want: []recipeData.Data{
				{
					Postcode: "10120",
					Recipe:   "Potato Curry",
					Delivery: "Wednesday 1AM - 7PM"},
				{
					Postcode: "10121",
					Recipe:   "Veggie Pasta",
					Delivery: "Friday 10AM - 4PM"},
				{
					Postcode: "10122",
					Recipe:   "Mushroom Risotto",
					Delivery: "Tuesday 2AM - 5PM"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := fsJson.NewReader()
			got, err := r.Read(tt.arg)
			assert.NoError(t, err, "Error opening file")
			assert.NotNil(t, file, "File was not opened successfully")
			if err != nil {
				t.Errorf("reader.Read() error = %v", err)
				return
			}
			for _, val := range tt.want {
				assert.Equal(t, <-got, val)
			}
		})
	}
}
