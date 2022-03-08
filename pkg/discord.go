package pkg

import (
	"net/http"
	"net/url"
)

type Discord struct {
	Endpoint string
}

func NewDiscord(target string) *Discord {
	return &Discord{Endpoint: target}
}

func (d *Discord) Send(message string) error {
	vals := url.Values{}
	vals.Set("content", message)
	_, err := http.PostForm(d.Endpoint, vals)
	return err
}
