package response

type Response struct {
	Error string `json:"error"`
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
