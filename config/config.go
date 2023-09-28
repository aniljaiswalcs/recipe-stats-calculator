package cfg

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DataFile       string
	Output         string
	MatchedRecipes []string
	Postcode       string
	From           uint
	To             uint
	Output_Dir     string
	Output_File    string
}

func NewConfig(filepath string) *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.SetConfigFile(filepath)

	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	setDefaults(v)

	return &Config{
		DataFile:       v.GetString("data_file"),
		Output:         v.GetString("output"),
		Postcode:       v.GetString("postcode"),
		From:           v.GetUint("from"),
		To:             v.GetUint("to"),
		MatchedRecipes: v.GetStringSlice("matched_recipes"),
		Output_Dir:     v.GetString("output_dir"),
		Output_File:    v.GetString("output_file"),
	}

}

func setDefaults(v *viper.Viper) {
	v.SetDefault("data_file", "hf_test_calculation_fixtures.json")
	v.SetDefault("output", "stdout")
	v.SetDefault("postcode", "10120")
	v.SetDefault("from", 10)
	v.SetDefault("to", 3)
	v.SetDefault("matched_recipes", []string{"Potato", "Veggie", "Mushroom"})
	v.SetDefault("output_dir", "./")
	v.SetDefault("output_file", "stats_output.json")

}
