package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const qcSuffix = "_qc"
const dateField = "time"

func main() {
	url := "http://www.neracoos.org/erddap/tabledap/E05_aanderaa_all.json?station%2Cmooring_site_desc%2Cwater_depth%2Ctime%2Ccurrent_speed%2Ccurrent_speed_qc%2Ccurrent_direction%2Ccurrent_direction_qc%2Ccurrent_u%2Ccurrent_u_qc%2Ccurrent_v%2Ccurrent_v_qc%2Ctemperature%2Ctemperature_qc%2Cconductivity%2Cconductivity_qc%2Csalinity%2Csalinity_qc%2Csigma_t%2Csigma_t_qc%2Ctime_created%2Ctime_modified%2Clongitude%2Clatitude%2Cdepth&time%3E=2015-08-25T15%3A00%3A00Z&time%3C=2016-12-05T14%3A00%3A00Z"
	b, err := minMaxAvg(url, "current_speed", "temperature", "salinity")
	if err != nil {
		log.Println("error: ", err)
	}
	fmt.Println(string(b))
}

func minMaxAvg(url string, fields ...string) ([]byte, error) {
	if len(fields) == 0 {
		return nil, errors.New("at least one field is required")
	}
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
	columnsMap := make(map[string]int)
	for i, f := range respD.Table.ColumnNames {
		columnsMap[f] = i
	}
	var values []*valuesTime
	for _, f := range fields {
		var vt = valuesTime{column: f}
		var ok bool
		vt.index, ok = columnsMap[f]
		if !ok {
			return nil, fmt.Errorf("column %+v is not found in response", f)
		}
		vt.qcIndex, ok = columnsMap[f+qcSuffix]
		if !ok {
			return nil, fmt.Errorf("column %+v is not found in response", f+qcSuffix)
		}
		vt.dateIndex, ok = columnsMap[dateField]
		if !ok {
			return nil, fmt.Errorf("column %+v is not found in response", dateField)
		}
		values = append(values, &vt)
	}
	for _, r := range respD.Table.Rows {
		for _, o := range values {
			if len(r) <= o.index || len(r) <= o.qcIndex || len(r) <= o.dateIndex {
				return nil, fmt.Errorf("invalid row, not enough elements in the array: %+v", r)
			}
			qc, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[o.qcIndex]), 64)
			if err != nil {
				return nil, err
			}
			if qc != 0 {
				continue
			}
			v, err := strconv.ParseFloat(fmt.Sprintf("%+v", r[o.index]), 64)
			if err != nil {
				return nil, err
			}
			date := fmt.Sprintf("%+v", r[o.dateIndex])
			if o.count == 0 {
				o.min = v
				o.max = v
			} else {
				if o.min > v {
					o.min = v
				}
				if o.max < v {
					o.max = v
				}
			}
			if len(o.start_date) == 0 {
				o.start_date = date
			}
			o.end_date = date // assuming data in rows is sorted: newest (with latest time field) data is in the end of the array
			o.count++
			o.sum += v
		}
	}
	var result = make(map[string]interface{})
	for _, o := range values {
		var avg float64
		if o.count != 0 {
			avg = o.sum / float64(o.count)
		}
		result[o.column] = map[string]interface{}{
			"start_date":      o.start_date,
			"end_date":        o.end_date,
			"num_records":     o.count,
			"min_" + o.column: o.min,
			"max_" + o.column: o.max,
			"avg_" + o.column: avg,
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
	ColumnNames []string        `json:"columnNames"`
	Rows        [][]interface{} `json:"rows"`
}

type valuesTime struct {
	start_date, end_date      string
	index, qcIndex, dateIndex int
	column                    string
	count                     int
	min, max, sum             float64
}
