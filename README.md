# Task
This test will measure your ability to work with large amounts of data obtained via API.
The following URL can be used to obtain surface current data for a point in the Atlantic Ocean from 2015-08-15 - 2016-12-05. Write the necessary code to:

1) Fetch data from the API,
2) Calculate the min, max and average current speed, salinity and water temperature for all records.

Your code should output the result in json format such as 

results = {
	current_speed: {
		start_date: "2015-08-15",
		end_date: "2016-12-05",
		num_records: 10000,
		min_current_speed: 0.0,
		max_current_speed: 32.81,
		avg_current_speed: 22.45,
	},
	...
}

**PLEASE NOTE: All _qc values != 0 indicate a faulty sensor reading and should not be included in your calculations.**


http://www.neracoos.org/erddap/tabledap/E05_aanderaa_all.json?station%2Cmooring_site_desc%2Cwater_depth%2Ctime%2Ccurrent_speed%2Ccurrent_speed_qc%2Ccurrent_direction%2Ccurrent_direction_qc%2Ccurrent_u%2Ccurrent_u_qc%2Ccurrent_v%2Ccurrent_v_qc%2Ctemperature%2Ctemperature_qc%2Cconductivity%2Cconductivity_qc%2Csalinity%2Csalinity_qc%2Csigma_t%2Csigma_t_qc%2Ctime_created%2Ctime_modified%2Clongitude%2Clatitude%2Cdepth&time%3E=2015-08-25T15%3A00%3A00Z&time%3C=2016-12-05T14%3A00%3A00Z

# Solution

```
$ go get github.com/ypapax/go-test1
$ cd $GOPATH/src/github.com/ypapax/go-test1
$ go run main.go | jq .
```
```json
{
  "current_speed": {
    "avg_current_speed": 21.631215155608178,
    "end_date": "2016-12-05T14:00:00Z",
    "max_current_speed": 81.2441,
    "min_current_speed": 0,
    "num_records": 22203,
    "start_date": "2015-08-25T15:00:00Z"
  },
  "salinity": {
    "avg_salinity": 31.834862969460694,
    "end_date": "2015-09-15T23:40:00Z",
    "max_salinity": 32.16124,
    "min_salinity": 31.116806,
    "num_records": 1539,
    "start_date": "2015-08-25T15:00:00Z"
  },
  "temperature": {
    "avg_temperature": 12.697609162955994,
    "end_date": "2016-12-05T14:00:00Z",
    "max_temperature": 19.601141,
    "min_temperature": 6.325187,
    "num_records": 22203,
    "start_date": "2015-08-25T15:00:00Z"
  }
}
```
