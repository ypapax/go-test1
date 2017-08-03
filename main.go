package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
)

func main(){
	url := "http://www.neracoos.org/erddap/tabledap/E05_aanderaa_all.json?station%2Cmooring_site_desc%2Cwater_depth%2Ctime%2Ccurrent_speed%2Ccurrent_speed_qc%2Ccurrent_direction%2Ccurrent_direction_qc%2Ccurrent_u%2Ccurrent_u_qc%2Ccurrent_v%2Ccurrent_v_qc%2Ctemperature%2Ctemperature_qc%2Cconductivity%2Cconductivity_qc%2Csalinity%2Csalinity_qc%2Csigma_t%2Csigma_t_qc%2Ctime_created%2Ctime_modified%2Clongitude%2Clatitude%2Cdepth&time%3E=2015-08-25T15%3A00%3A00Z&time%3C=2016-12-05T14%3A00%3A00Z"
	MinMaxAvg(url)
}

func MinMaxAvg(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Status >= 400 {
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
	var speedMax, speedMin, speedSum float64 = -1, -1, -1
	var speedCount int
	for i, r := range respD.Table.Rows {
		if len(r) < 18 {
			return nil, fmt.Errorf("invalid row, not enough elements in the array: %+v", r)
		}
		speed := r[4]
		speed_qc := r[5]
		if speed_qc != 0 {
			if speedMax <
		}

		temperature := r[12]
		temperature_qc := r[13]
		if temperature_qc != 0 {
			temperatures = append(temperatures, temperature)
		}

		salinity := r[16]
		salinity_qc := r[17]

		if salinity_qc != 0 {
			salinities = append(salinities, salinity)
		}
	}
}

type respData struct {
	Table table `json:"table"`
}

type table struct {
	Rows [][]interface{} `json:"rows"`
}