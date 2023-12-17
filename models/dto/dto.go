package dto

import "encoding/json"


type Note struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Note      string `json:"note,omitempty"`
}

func NewNote() *Note {
	return &Note{ ID: -1 }
}

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}


func (r *Response) Wrap(result string, data json.RawMessage, error string) {
	r.Result = result
	r.Error = error
	r.Data = data
}
