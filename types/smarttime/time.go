package smarttime

import(
	"time"
	"bytes"
	"encoding/json"
	"strconv"
	"errors"
	// "reflect"
)

const timeTemplate = "2006-01-02 15:04:05"

type Time time.Time

func Parse(data interface{}) (Time, error) {
	// t := reflect.TypeOf(data)
	switch i := data.(type) {
	case time.Time:
		return Time(i), nil
	case int:
		timestamp := int64(i)
		sec := timestamp / 1000
		nano := (timestamp % 1000) * 1e6	
		t := time.Unix(sec, nano)
		return Time(t), nil
	case int32:
		timestamp := int64(i)
		sec := timestamp / 1000
		nano := (timestamp % 1000) * 1e6	
		t := time.Unix(sec, nano)
		return Time(t), nil
	case int64:
		timestamp := int64(i)
		sec := timestamp / 1000
		nano := (timestamp % 1000) * 1e6	
		t := time.Unix(sec, nano)
		return Time(t), nil
	case string:
		var t time.Time
		timestamp, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			t, err = time.ParseInLocation(timeTemplate, i, time.Local)
			if err != nil {
				return Time(t), err
			}
			return Time(t), nil
		} else {
			sec := timestamp / 1000
			nano := (timestamp % 1000) * 1e6	
			t = time.Unix(sec, nano)
			return Time(t), nil
		}
	default:
		return Time(time.Unix(0, 0)), errors.New("time parse error")
	}
}

func (self *Time) UnmarshalJSON(data []byte) error {
	var timeStr string
	decoder := json.NewDecoder(bytes.NewReader(data))
	// decoder.UseNumber()
	if err := decoder.Decode(&timeStr); err != nil {
		return err
	}

	var err error
	*self, err = Parse(timeStr)
	return err
}

func (self Time) MarshalJSON() ([]byte, error) {
	t := time.Time(self)
	timestamp := t.UnixNano() / 1e6
	if timestamp <= 0 {
		timestamp = 0
	}

	return json.Marshal(timestamp)
}
