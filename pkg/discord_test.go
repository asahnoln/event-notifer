package pkg_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/asahnoln/event-notifier/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscordUsingKey(t *testing.T) {
	sdr := pkg.NewDiscord("https://discord.com/secretwebhook")

	assert.Equal(t,
		"https://discord.com/secretwebhook",
		sdr.Endpoint,
		"want proper discord endpoint",
	)
}

func TestDiscordSendSuccess(t *testing.T) {
	var gotVals url.Values
	ds := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		gotVals = r.PostForm
	}))
	es := []pkg.Event{
		{
			What:  "Good thing",
			Where: "Paradise",
			Who:   []string{"Angel", "God"},
			Start: "01.02.1999 09:00",
			End:   "01.02.1999 10:00",
		},
		{
			What:  "Bad thing",
			Where: "Hell",
			Who:   []string{"Demon", "Satan"},
			Start: "05.12.1999 07:00",
			End:   "05.12.1999 08:00",
		},
	}

	sdr := pkg.NewDiscord(ds.URL)
	err := pkg.Send(es, sdr, pkg.Tomorrow)

	require.NoError(t, err, "unexpected error while sending message to Discord")
	require.True(t, gotVals.Has("content"), "want content value set")
	for _, want := range es {
		assert.Contains(t, gotVals.Get("content"), want.What, "content should contain event info")
	}
}
