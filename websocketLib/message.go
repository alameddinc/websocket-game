package websocketLib

import "encoding/json"

const (
	UP    = 1
	DOWN  = 2
	LEFT  = 3
	RIGHT = 4
)

type messageSchema struct {
	TypeId uint   `json:"type_id"`
	UserId string `json:"user_id"`
}

func (m *messageSchema) Encode() []byte {
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return jsonMessage
}
