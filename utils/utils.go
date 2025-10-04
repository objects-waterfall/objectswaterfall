package utils

import (
	"errors"
	"math/rand"
	"time"
)

func RandomString(min, max int) (string, error) {
	if min > max {
		return "", errors.New("min length cannot be greater than max length")
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(max-min+1) + min

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result), nil
}

func RandomBool() bool {
	return rand.Intn(2) == 1
}

func RandomDate(start, end time.Time) time.Time {
	duration := end.Sub(start)
	randomDuration := time.Duration(rand.Int63n(duration.Nanoseconds()))
	return start.Add(randomDuration)
}

func FillWithDummyData(obj *interface{}) error {
	switch val := (*obj).(type) {
	case bool:
		*obj = RandomBool()
	case string:
		str, err := RandomString(10, 50)
		if err != nil {
			return err
		}
		*obj = str
	case float64:
		switch val {
		case float64(int(val)):
			*obj = rand.Intn(50)
		case float64(int32(val)):
			*obj = rand.Int31n(50)
		case float64(int64(val)):
			*obj = rand.Int63n(50)
		case float64(float32(val)):
			*obj = rand.Float32()
		default:
			*obj = rand.Float64()
		}
	case time.Time:
		*obj = RandomDate(
			time.Now().Add(time.Duration(-rand.Float32())),
			time.Now().Add(time.Duration(rand.Float32())))
	case map[string]interface{}:
		for k, v := range val {
			ptr := v
			FillWithDummyData(&ptr)
			val[k] = ptr
		}
	case []interface{}:
		for k := range val {
			FillWithDummyData(&val[k])
		}
	}

	return nil
}
