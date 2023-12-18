package san

import (
	"encoding/json"
	"time"
)

type TsInt64 int64

var shouldReadableFlag bool = true
func SetTsReadable(shouldReadable bool) {
	shouldReadableFlag = shouldReadable
}

func (t *TsInt64) MarshalJSON() ([]byte, error) {
	if shouldReadableFlag {
		return []byte(time.Unix(int64(*t), 0).Format("\"2006-01-02 15:04:05\"")), nil
	}
	var v int64 = int64(*t)
	return json.Marshal(v)
}

func (t *TsInt64) UnmarshalJSON(data []byte) error {
	if shouldReadableFlag {
		tt, err := time.Parse("\"2006-01-02 15:04:05\"", string(data))
		if err != nil {
			return err
		}
		*t = TsInt64(tt.Unix())
		return nil
	}
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err 
	}
	*t = TsInt64(v)
	return nil
}
