package easycsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// WriteCSVData ...
func WriteCSVData(t interface{}, pathfile string) {
	data := manipulateReflection(t)

	file, err := os.Create(pathfile)
	if err != nil {
		fmt.Println("Fail Create File")
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			fmt.Println("Fail Write CSV")
			panic(err)
		}
	}

}

func manipulateReflection(t interface{}) (data [][]string) {
	var tags []string
	data = append(data, tags)

	if reflect.TypeOf(t).Kind() == reflect.Slice {
		s := reflect.ValueOf(t)
		loopTags := 0

		for x := 0; x < s.Len(); x++ {
			v := reflect.Indirect(s.Index(x))
			vt := v.Type()

			var dataNested []string
			for i, n := 0, v.NumField(); i < n; i++ {
				ft := vt.Field(i)
				if ft.Tag.Get("csv") != "" && loopTags == 0 {
					tags = append(tags, ft.Tag.Get("csv"))
				}

				f := v.Field(i)
				n := checkTypeAndReturnString(f.Interface())
				if n != "" {
					dataNested = append(dataNested, n)
				}
			}
			data = append(data, dataNested)

			loopTags++
		}
	}
	data[0] = tags

	return
}

func checkTypeAndReturnString(t interface{}) (n string) {
	switch v := t.(type) {
	case int:
		n = strconv.Itoa(v)
	case int64:
		n = strconv.FormatInt(v, 10)
	case string:
		n = v
	case float64:
		n = strconv.FormatFloat(v, 'f', 5, 64)
	case time.Time:
		n = v.String()
	default:
		fmt.Println("LOG ERROR DATA == Parsing WriteVal")
		fmt.Println(v)
		fmt.Println(reflect.TypeOf(v))
	}

	return
}
