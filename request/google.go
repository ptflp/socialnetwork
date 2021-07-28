package request

type GoogleCallbackResponse struct {
	GoogleID  string `json:"id"`
	Name      string `json:"name"`
	GivenName string `json:"given_name"`
	Avatar    string `json:"picture"`
	Locale    string `json:"locale"`
}
