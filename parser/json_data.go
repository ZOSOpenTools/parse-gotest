package parser

import (
	"encoding/json"
)

type JsonData struct {
	Time string
	Action string
	Package string
	Test string 
	Output string 
	Elapsed float64
}

func (jd *JsonData) UnmarshalJSON(b []byte) error {
	type Alias JsonData 
	type Aux struct {
		Test *string `json:"Test"`
		Output *string `json:"Output"`
		Elapsed *float64 `json:"Elapsed"`
		*Alias
	}
	aux := &Aux{Alias: (*Alias)(jd)}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	if aux.Test == nil {
		jd.Test = ""
	} else {
		jd.Test = *aux.Test
	}

	if aux.Output == nil {
		jd.Output = ""
	} else {
		jd.Output = *aux.Output
	}

	if aux.Elapsed == nil {
		jd.Elapsed = -1
	} else {
		jd.Elapsed = *aux.Elapsed
	}
	return nil
}