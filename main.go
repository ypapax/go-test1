package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	url := "http://www.neracoos.org/erddap/tabledap/E05_aanderaa_all.json?station%2Cmooring_site_desc%2Cwater_depth%2Ctime%2Ccurrent_speed%2Ccurrent_speed_qc%2Ccurrent_direction%2Ccurrent_direction_qc%2Ccurrent_u%2Ccurrent_u_qc%2Ccurrent_v%2Ccurrent_v_qc%2Ctemperature%2Ctemperature_qc%2Cconductivity%2Cconductivity_qc%2Csalinity%2Csalinity_qc%2Csigma_t%2Csigma_t_qc%2Ctime_created%2Ctime_modified%2Clongitude%2Clatitude%2Cdepth&time%3E=2015-08-25T15%3A00%3A00Z&time%3C=2016-12-05T14%3A00%3A00Z"
	b, err := MinMaxAvg(url)
	if err != nil {
		log.Println("error: ", err)
	}
	fmt.Println(string(b))
}

func MinMaxAvg(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("invalid status code %+v for requesting url %+v", resp.Status, url)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respD respData
	if err := json.Unmarshal(b, &respD); err != nil {
		return nil, err
	}
	var speeds, temperatures, salinities valuesTime
	for _, r := range respD.Table.Rows {
		if len(r) < 18 {
			return nil, fmt.Errorf("invalid row, not enough elements in the array: %+v", r)
		}

		speed, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[4]), 64)
		if err != nil {
			return nil, err
		}
		speed_qc, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[5]), 64)
		if err != nil {
			return nil, err
		}
		date := fmt.Sprintf("%+v", r[3])
		if speed_qc != 0 {
			if len(speeds.start_date) == 0 {
				speeds.start_date = date
			}
			speeds.end_date = date
			speeds.values = append(speeds.values, speed)
		}

		temperature, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[12]), 64)
		if err != nil {
			return nil, err
		}
		temperature_qc, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[13]), 64)
		if err != nil {
			return nil, err
		}
		if temperature_qc != 0 {
			if len(temperatures.start_date) == 0 {
				temperatures.start_date = date
			}
			temperatures.end_date = date
			temperatures.values = append(temperatures.values, temperature)
		}

		salinity, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[16]), 64)
		if err != nil {
			return nil, err
		}
		salinity_qc, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[17]), 64)
		if err != nil {
			return nil, err
		}
		if salinity_qc != 0 {
			if len(salinities.start_date) == 0 {
				salinities.start_date = date
			}
			salinities.end_date = date
			salinities.values = append(salinities.values, salinity)
		}
	}
	var result = make(map[string]interface{})
	for k, v := range map[string]valuesTime{
		"current_speed": speeds,
		"temperature":   temperatures,
		"salinity":    salinities,
	} {
		min, max, avg := minMaxAvg(v.values)
		result[k] = map[string]interface{}{
			"start_date":  v.start_date,
			"end_date":    v.end_date,
			"num_records": len(v.values),
			"min_" + k:    min,
			"max_" + k:    max,
			"avg_" + k:    avg,
		}
	}
	b, err = json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type respData struct {
	Table table `json:"table"`
}

type table struct {
	Rows [][]interface{} `json:"rows"`
}

func minMaxAvg(values []float64) (min, max, avg float64) {
	var sum float64
	if len(values) == 0 {
		return
	}
	min = values[0]
	max = values[0]
	for _, v := range values {
		if v < min {
			v = min
		} else if v > max {
			v = max
		}
		sum += v
	}
	return min, max, sum / float64(len(values))
}

type valuesTime struct {
	values               []float64
	start_date, end_date string
}
