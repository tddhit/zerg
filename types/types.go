package types

import (
	"net/http"
)

type Request struct {
	*http.Request
	RawURL string
}

type Response struct {
	*http.Response
}

type Item struct {
	//associatedWriter Writer
	Dict map[string]string
}
