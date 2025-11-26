package easyExpression

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FormulaAction struct {
}

/*-----------------Math---------------------------*/

func (f FormulaAction) Sum(values ...any) interface{} {
	result := float64(0)
	for _, v := range values {
		//if v is string
		if _, ok := v.(string); ok {
			value := strings.Replace(v.(string), " ", "", -1)
			temp, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic("function sum error: " + v.(string) + "not a number")
			}
			result = result + temp
		} else if _, ok := v.(float64); ok {
			result = result + v.(float64)
		}
	}
	return result
}
func (f FormulaAction) Avg(values ...any) interface{} {
	result := float64(0)
	for _, v := range values {
		//if v is string
		if _, ok := v.(string); ok {
			value := strings.Replace(v.(string), " ", "", -1)
			temp, err := strconv.ParseFloat(value, 64)
			if err != nil {
				panic("function sum error: " + v.(string) + "not a number")
			}
			result = result + temp
		} else if _, ok := v.(float64); ok {
			result = result + v.(float64)
		}
	}
	return result / float64(len(values))
}
func (f FormulaAction) Round(values ...any) interface{} {
	if len(values) == 0 {
		panic("function Round called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() != reflect.Slice || len(values[0].([]interface{})) != 3 {
		panic("function Round called with invalid arguments")
	}
	array := values[0].([]interface{})
	v := InterfaceToFloat64(array[0])
	accuracy := InterfaceToFloat64(array[1])
	mode := InterfaceToFloat64(array[2])
	var delta = 5 / math.Pow(10, accuracy+1)
	switch mode {
	case -1:
		return f.CustomerRound(v-delta, accuracy)
	case 0:
		return f.CustomerRound(v, accuracy)
	case 1:
		return f.CustomerRound(v+delta, accuracy)
	}
	panic("round mode error")
}

/*-----------------Math---------------------------*/

/*-----------------String---------------------------*/

func (f FormulaAction) Contains(values ...any) interface{} {
	if len(values) == 0 {
		panic("function equals called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() == reflect.Slice && len(values[0].([]interface{})) == 2 {
		array := values[0].([]interface{})
		key := fmt.Sprintf("%v", array[0])
		str := fmt.Sprintf("%v", array[1])
		if strings.Contains(key, str) {
			return float64(1)
		}
		return float64(0)
	}
	return float64(0)
}
func (f FormulaAction) Excluding(values ...any) interface{} {
	key := values[0].(string)
	str := values[1].(string)
	if len(key) == 0 {
		return float64(0)
	}
	if strings.Contains(str, key) {
		return float64(0)
	}
	return float64(1)
}
func (f FormulaAction) Equals(values ...any) interface{} {
	if len(values) == 0 {
		panic("function equals called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() == reflect.Slice && len(values[0].([]interface{})) == 2 {
		array := values[0].([]interface{})
		key := fmt.Sprintf("%v", array[0])
		str := fmt.Sprintf("%v", array[1])
		if key == str {
			return float64(1)
		}
		return float64(0)
	}
	panic("function equals called with invalid arguments")
}
func (f FormulaAction) StartWith(values ...any) interface{} {
	key := values[0].(string)
	str := values[1].(string)
	if len(key) == 0 {
		return float64(1)
	}
	if strings.HasPrefix(str, key) {
		return float64(1)
	}
	return float64(0)
}
func (f FormulaAction) EndWith(values ...any) interface{} {
	key := values[0].(string)
	str := values[1].(string)
	if len(key) == 0 {
		return float64(1)
	}
	if strings.HasSuffix(str, key) {
		return float64(1)
	}
	return float64(0)
}
func (f FormulaAction) Different(values ...any) interface{} {
	if len(values) == 0 {
		panic("function Different called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() == reflect.Slice && len(values[0].([]interface{})) == 2 {
		array := values[0].([]interface{})
		key := fmt.Sprintf("%v", array[0])
		str := fmt.Sprintf("%v", array[1])
		if key == str {
			return float64(0)
		}
		return float64(1)
	}
	panic("function Different called with invalid arguments")
}

/*-----------------String---------------------------*/

/*-----------------Time---------------------------*/

func (f FormulaAction) EDate(values ...any) interface{} {
	if len(values) == 0 {
		panic("function EDate called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() != reflect.Slice || len(values[0].([]interface{})) != 3 {
		panic("function EDate called with invalid arguments")
	}
	array := values[0].([]interface{})
	if date, ok := array[0].(time.Time); ok {
		value := InterfaceToInt(array[1])
		format := array[2].(string)
		switch format {
		case "Y", "y":
			date = date.AddDate(value, 0, 0)
			return date
		case "M":
			date = date.AddDate(0, value, 0)
			return date
		case "D", "d":
			date = date.AddDate(0, 0, value)
			return date
		case "H", "h":
			date = date.Add(time.Hour * time.Duration(value))
			return date
		case "m":
			date = date.Add(time.Minute * time.Duration(value))
			return date
		case "S", "s":
			date = date.Add(time.Second * time.Duration(value))
			return date
		case "F", "f":
			date = date.Add(time.Millisecond * time.Duration(value))
			return date
		}
	}
	panic("date parse error")
}
func (f FormulaAction) EODate(values ...any) interface{} {
	if len(values) == 0 {
		panic("function EODate called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() != reflect.Slice || len(values[0].([]interface{})) != 3 {
		panic("function EODate called with invalid arguments")
	}
	array := values[0].([]interface{})
	if date, ok := array[0].(time.Time); ok {
		value := InterfaceToInt(array[1])
		format := array[2].(string)
		newDate := date.AddDate(0, value, 0)
		switch format {
		case "S", "s":
			return time.Date(newDate.Year(), newDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		case "E", "e":
			newDate = time.Date(newDate.Year(), newDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			newDate = newDate.AddDate(0, 1, 0)
			newDate = newDate.AddDate(0, 0, -1)
			return newDate
		}
	}
	panic("EODate execute error")
}
func (f FormulaAction) NowTime(values ...any) interface{} {
	return time.Now()
}
func (f FormulaAction) TimeToString(values ...any) interface{} {
	if len(values) == 0 {
		panic("function TimeToString called with no arguments")
	}
	if reflect.TypeOf(values[0]).Kind() != reflect.Slice || len(values[0].([]interface{})) != 2 {
		panic("function TimeToString called with invalid arguments")
	}
	array := values[0].([]interface{})
	if date, ok := array[0].(time.Time); ok {
		value := array[1].(string)
		formatting := "2006-01-02 15:04:05"
		if len(value) > 1 {
			formatting = value
			formatting = strings.Replace(formatting, "yyyy", "2006", -1)
			formatting = strings.Replace(formatting, "YYYY", "2006", -1)
			formatting = strings.Replace(formatting, "MM", "01", -1)
			formatting = strings.Replace(formatting, "dd", "02", -1)
			formatting = strings.Replace(formatting, "DD", "02", -1)
			formatting = strings.Replace(formatting, "HH", "15", -1)
			formatting = strings.Replace(formatting, "hh", "15", -1)
			formatting = strings.Replace(formatting, "mm", "04", -1)
			formatting = strings.Replace(formatting, "ss", "05", -1)
			formatting = strings.Replace(formatting, "SS", "05", -1)
		}
		return date.Format(formatting)
	}

	panic("TimeToString execute error")
}

func (f FormulaAction) Days(values ...any) interface{} {
	if duration, ok := values[0].([]interface{})[0].(time.Duration); ok {
		return duration.Hours() / 24
	}
	panic("Days execute error")
}
func (f FormulaAction) Hours(values ...any) interface{} {
	if duration, ok := values[0].([]interface{})[0].(time.Duration); ok {
		return duration.Hours()
	}
	panic("Hours execute error")
}
func (f FormulaAction) Minutes(values ...any) interface{} {
	if duration, ok := values[0].([]interface{})[0].(time.Duration); ok {
		return duration.Minutes()
	}
	panic("Minutes execute error")
}
func (f FormulaAction) Seconds(values ...any) interface{} {
	if duration, ok := values[0].([]interface{})[0].(time.Duration); ok {
		return duration.Seconds()
	}
	panic("Seconds execute error")
}
func (f FormulaAction) MillSeconds(values ...any) interface{} {
	if duration, ok := values[0].([]interface{})[0].(time.Duration); ok {
		return duration.Milliseconds()
	}
	panic("MillSeconds execute error")
}

/*-----------------Time---------------------------*/

/*-----------------object---------------------------*/

func (f FormulaAction) IsNull(values ...any) interface{} {
	if reflect.ValueOf(values[0]).IsNil() {
		return float64(1)
	} else {
		return float64(0)
	}
}

/*-----------------object---------------------------*/
