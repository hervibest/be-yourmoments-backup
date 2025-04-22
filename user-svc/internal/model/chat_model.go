package model

type RequestGetOrCreateRoom struct {
	SenderId   string `validate:"required"`
	ReceiverId string `json:"receiver_id" validate:"required"`
}

type GetOrCreateRoomResponse struct {
	RoomId  string `json:"room_id"`
	Created bool   `json:"created"`
}

type RequestCustomToken struct {
	UserId string `json:"user_id" validate:"required"`
}

// TODO memastikan snake_case bukan pascalCase
type RequestSendMessage struct {
	RoomId   string `json:"room_id" validate:"required"`
	SenderId string `validate:"required"`
	Message  string `json:"message" validate:"required"`
}

type CustomTokenResponse struct {
	Token string `json:"token"`
}
