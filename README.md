# Recipe stats count application

## Setup

### Prerequisite
1. go1.20
2. Docker for Desktop 
3. `hf_test_calculation_fixtures.json` file in present in root of the project, you can change file name in config.yml if you want. If input file name change then change in Dockerfile also.
### Clone Repository
Clone this repository somewhere on your computer and `cd aniljaiswalcs-HFtest-backend-go`
```
git clone -b anil-branch git@github.com:hellofreshdevtests/aniljaiswalcs-HFtest-backend-go.git
```



## Application setup
1. update config.yml file, if you need to update filename, postcode, time, receipe, output file dir and file name.
2. Run from console or Docker setup. 
3. Provide info for calculations
4. Based on the input provide for Output in config.yml, json output will be shown on console for "stdout" or redirect to output file "stat_output.json" for "file" input.

## Usage
### Run application from console:
1. To build the source code:
```    
make cli_build
```
2. To run the cli applicattion:
```
make cli_run
```

3. To clean the cli applicattion:
```
make clean
```

### Run application from docker:

1. create image :
```
make build
```

2. Run image:
```
make run
```

3. Delete container and image:
```
make docker_clean
```

4. Run test:
```
make test
```

## Notes

#### Regarding Json file reading

For loading the json file, viper package is used. Which read the filename from config.yml and pass it to Read().

Read() method in `handler/dataReader.go` is used to parse given Json file and return each Json object in reader channel which is then received in `main.go` and passed into CalculateStats() method in `repository/aggregate.go` for stats calculations.

Read() method doesn't load whole Json file in memory, it instead reads each row gradually using json.NewDecoder() and pass it to CalculateStats() through channel.

#### Regarding a functional requirement #4 and 5

Below confinguration for postcode, from and to read from config.yml and pass it to main.go for further processing.
```
postcode: "10224"
from: 1
to: 7
```

Regaring the Functional req #5, everything is configurable in config.yml.

#### Weekday check

_Important notes_

1. Property value `"delivery"` always has the following format: "{weekday} {h}AM - {h}PM", i.e. "Monday 9AM - 5PM"

Here clearly mention about {weekday} i.e. Monday to Friday only but hf_test_calculation_fixtures.json file contains Saturday data also which is quite big number.

```
grep -nr "Saturday" hf_test_calculation_fixtures.json | wc -l
1667026
```
I am skipping data for Saturday.
