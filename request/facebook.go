package request

//go:generate easytags $GOFILE
type FacebookCallbackRequest struct {
	FacebookID string `json:"id"`
	Name       string `json:"name"`
}
