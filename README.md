# go-test1

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