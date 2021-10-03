package pkg_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/asahnoln/event-notifier/pkg"
)

func TestTelegramUsingKey(t *testing.T) {
	sdr := pkg.NewTg("0123verySecretKey", "someId")

	assertSameString(t,
		"https://api.telegram.org/bot0123verySecretKey/sendMessage",
		sdr.Endpoint,
		"want Telegram endpoint %q, got %q",
	)
}

// TODO: Does it have a guarantee that correct Endpoint is used? Probably
func TestTelegramSendSuccess(t *testing.T) {
	var gotVals url.Values
	ts := fakeTgServer(&gotVals)
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

	sdr := pkg.NewTg("secretKey", "testId")
	sdr.Endpoint = ts.URL
	err := pkg.Send(es, sdr)

	assertNoError(t, err, "unexpected error while sending message to Telegram: %v")
	for _, want := range es {
		assertContains(t, want.What, gotVals.Get("text"))
	}
	assertSameString(t, "testId", gotVals.Get("chat_id"), "want chat id %q, got %q")
}

func fakeTgServer(vals *url.Values) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		*vals = r.PostForm
	}))

}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
