package pkg_test

import (
	"reflect"
	"strings"
	"testing"
)

func assertNoError(t testing.TB, err error, message string) {
	t.Helper()

	if err != nil {
		t.Fatalf(message, err)
	}
}

func assertError(t testing.TB, err error, message string) {
	t.Helper()

	if err == nil {
		t.Fatal(message)
	}
}

func assertContains(t testing.TB, want, got string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Errorf("want substring %q in string %q, don't have it", want, got)
	}
}

func assertSameStruct(t testing.TB, want, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want event structure %+v, got %+v", want, got)
	}
}

func assertSameLength(t testing.TB, want, got int) {
	t.Helper()

	if want != got {
		t.Fatalf("want events length %d, got %d", want, got)
	}
}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}
