package san

import (
	"encoding/json"
	"time"
)

type Version struct {
	Time int64
	Hash string
	Msg  string
}

func (v *Version) MarshalJSON() ([]byte, error) {
	var vv = struct {
		Time string
		Hash string
		Msg  string
	}{
		Time: time.Unix(v.Time, 0).Format("2006-01-02 15:04:05"),
		Hash: v.Hash,
		Msg:  v.Msg,
	}
	return json.Marshal(vv)
}

func (v *Version) UnmarshalJSON(data []byte) error {
	var vv = struct {
		Time string
		Hash string
		Msg  string
	}{}

	if err := json.Unmarshal(data, &vv); err != nil {
		return err
	}
	v.Hash = vv.Hash
	v.Msg = vv.Msg
	t, err := time.Parse("2006-01-02 15:04:05", vv.Time)
	if err != nil {
		return err
	}
	v.Time = t.Unix()
	return nil
}
