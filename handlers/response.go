package handlers

type Response struct {
    Code int      `json:"code"`
    Msg  string   `json:"msg"`
    Data []string `json:"data"`
}

func Success(data []string) *Response {
    return &Response{
        Code: 200,
        Msg:  "ok",
        Data: data,
    }
}

func Error(err error) *Response {
    return &Response{
        Code: 500,
        Msg:  err.Error(),
        Data: nil,
    }
}

func ErrorM(msg string) *Response {
    return &Response{
        Code: 500,
        Msg:  msg,
        Data: nil,
    }
}
