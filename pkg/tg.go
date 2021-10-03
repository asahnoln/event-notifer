package pkg

import (
	"net/http"
	"net/url"
)

type Tg struct {
	chatId, Endpoint string
}

func NewTg(key string, chatId string) *Tg {
	return &Tg{
		chatId,
		"https://api.telegram.org/bot" + key + "/sendMessage",
	}
}

func (t *Tg) Send(message string) error {
	resp, err := http.PostForm(
		t.Endpoint,
		url.Values{
			"chat_id": {t.chatId},
			"text":    {message},
		},
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
