package pkg

import (
	"encoding/json"
	"io"
)

type ErrorNameMissing struct {
	Mail string
}

func (e *ErrorNameMissing) Error() string {
	return "Missing name for mail " + e.Mail
}

func MailsToNames(ms []string, r io.Reader) ([]string, error) {
	data := make(map[string]string, len(ms))
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(ms))
	for i, m := range ms {
		var name string
		var ok bool
		if name, ok = data[m]; !ok {
			return names, &ErrorNameMissing{m}
		}

		names[i] = name
	}

	return names, nil
}
