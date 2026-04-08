package types

// Generic API response wrapper
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Auth
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// User
type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Gender   int32  `json:"gender"`
}

// Friend
type SendFriendRequestReq struct {
	ToUserID int64  `json:"to_user_id"`
	Remark   string `json:"remark"`
}

type RespondFriendRequestReq struct {
	RequestID int64  `json:"request_id"`
	Action    string `json:"action"` // "accept" or "reject"
}
