package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	cfg "github.com/hellofreshdevtests/aniljaiswalcs-HFtest-backend-go/config"
)

// CountPerRecipe represents a unique recipe with the number of its occurrences in ExpectedOutput.
type CountPerRecipe struct {
	Recipe string `json:"recipe"`
	Count  uint   `json:"count"`
}

// BusiestPostcode represents the postcode with most delivered recipes in ExpectedOutput.
type BusiestPostcode struct {
	Postcode      string `json:"postcode"`
	DeliveryCount uint   `json:"delivery_count"`
}

// CountPerPostcodeAndTime represents delivery count to a postcode at certain time in ExpectedOutput.
type CountPerPostcodeAndTime struct {
	Postcode      string `json:"postcode"`
	FromAM        string `json:"from"`
	ToPM          string `json:"to"`
	DeliveryCount uint   `json:"delivery_count"`
}

// ExpectedOutput represents expected output with required stats.
type ExpectedOutput struct {
	UniqueRecipeCount       uint                    `json:"unique_recipe_count"`
	CountPerRecipe          []CountPerRecipe        `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostcode         `json:"busiest_postcode"`
	CountPerPostcodeAndTime CountPerPostcodeAndTime `json:"count_per_postcode_and_time"`
	MatchByName             []string                `json:"match_by_name"`
}

// Jsonify will convert ExpectedOutput to json string.
func (eo ExpectedOutput) Jsonify() (string, error) {
	b, err := json.MarshalIndent(eo, "", "\t")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (eo ExpectedOutput) GenerateOutput(config *cfg.Config) error {
	jsonOutput, err := eo.Jsonify()
	if err != nil {
		return err
	}

	if config.Output == "stdout" {
		fmt.Fprintf(os.Stdout, "%v\n", jsonOutput)
	} else if config.Output == "file" {
		err = eo.writeToFile(config.Output_Dir, config.Output_File)
		if err != nil {
			return err
		}
		stringFill := fmt.Sprintf(" Check for %s file for generated stats", config.Output_File)
		fmt.Fprintf(os.Stdout, "%s\n", stringFill)
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", " Please choose type in config.yml for Output")
	}
	return nil
}

// to write into file
func (eo ExpectedOutput) writeToFile(outputDir, outputFile string) error {

	b, err := eo.Jsonify()
	if err != nil {
		return err
	}
	_, err = os.ReadDir(outputDir)
	if err != nil {
		return err
	}
	os.Mkdir(outputDir, 0644)
	os.Chmod(outputDir, 0777)
	path := path.Join(outputDir, outputFile)
	err = os.WriteFile(path, []byte(b), 0777)
	if err != nil {
		return err
	}
	return nil
}
