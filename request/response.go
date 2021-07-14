package request

//go:generate easytags $GOFILE
type Response struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg,omitempty"`
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

type UserData struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	SecondName string `json:"second_name"`
}

type PostDataResponse struct {
	ID     int64         `json:"id"`
	Body   string        `json:"body"`
	Files  []string      `json:"files"`
	User   UserData      `json:"user"`
	Counts PostCountData `json:"counts"`
}

type PostCountData struct {
	Likes    int64 `json:"likes"`
	Comments int64 `json:"comments"`
}
