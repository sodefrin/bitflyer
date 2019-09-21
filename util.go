package bitflyer

import "time"

func parseTimeString(str string) (time.Time, error) {
	tmp, err := time.Parse(`2006-01-02T15:04:05`, string(str[0:19]))
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return tmp.In(jst), err
}
