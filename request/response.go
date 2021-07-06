package request

type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type AuthTokenResponse struct {
	Success bool          `json:"success"`
	Msg     string        `json:"msg"`
	Data    AuthTokenData `json:"data"`
}

type AuthTokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
