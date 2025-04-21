package utils

type MetaData struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

type Paginator struct {
	ItemCount   int  `json:"itemCount"`
	Limit       int  `json:"limit"`
	PageCount   int  `json:"pageCount"`
	Page        int  `json:"page"`
	HasPrevPage bool `json:"hasPrevPage"`
	HasNextPage bool `json:"hasNextPage"`
	PrevPage    *int `json:"prevPage,omitempty"`
	NextPage    *int `json:"nextPage,omitempty"`
}

type ResponseFormat[T any] struct {
	Meta      MetaData   `json:"meta"`
	Data      *T         `json:"data"`
	Paginator *Paginator `json:"paginator,omitempty"`
}

func ResponseSuccess[T any](code int, message string, data *T, paginator *Paginator) ResponseFormat[T] {
	return ResponseFormat[T]{
		Meta: MetaData{
			Code:    code,
			Success: true,
			Message: message,
		},
		Data:      data,
		Paginator: paginator,
	}
}

func ResponseError(code int, message interface{}) ResponseFormat[any] {
	return ResponseFormat[any]{
		Meta: MetaData{
			Code:    code,
			Success: false,
			Message: message,
		},
		Data: nil,
	}
}
