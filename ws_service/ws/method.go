package ws

const (
	//Heart
	HEART_BEAT = "HEART_BEAT"

	//Notify
	WS_SERVER_NOTIFY_TICK_INFO = "WS_SERVER_NOTIFY_TICK_INFO"

	//Connect
	WS_CONNECT = "WS_CONNECT"

	//common
	WS_RESPONSE_SUCCESS = "WS_RESPONSE_SUCCESS"
	WS_RESPONSE_ERROR   = "WS_RESPONSE_ERROR"

	//disconnect
	WS_DISCONNECT = "WS_DISCONNECT"
)

const (
	WS_CODE_HEART_BEAT      = 10
	WS_CODE_HEART_BEAT_BACK = 10
	WS_CODE_SERVER          = 0
	WS_CODE_SEND_SUCCESS    = 200
	WS_CODE_SEND_ERROR      = 400
)