package pkg_test

import (
	"strings"
	"testing"

	"github.com/asahnoln/event-notifier/pkg"
)

var source = strings.NewReader(`{
	"ivan@gmail.com": "Ivan Kolosov",
	"kamila@gmail.com": "Kamila Kenesbay"
}`)

func TestMailsToNames(t *testing.T) {
	source.Seek(0, 0)
	mails := []string{
		"ivan@gmail.com",
		"kamila@gmail.com",
	}
	names, err := pkg.MailsToNames(mails, source)

	assertNoError(t, err, "unexpected error while converting mails to names: %v")
	assertSameStruct(t, []string{
		"Ivan Kolosov",
		"Kamila Kenesbay",
	}, names)
}

func TestMailsToNamesWithMissingMail(t *testing.T) {
	source.Seek(0, 0)
	mails := []string{
		"ivan@gmail.com",
		"ilya@gmail.com",
	}
	_, err := pkg.MailsToNames(mails, source)

	want := &pkg.ErrorNameMissing{"ilya@gmail.com"}
	assertSameString(t, want.Error(), err.Error(), "want error %q, got %q")
}
