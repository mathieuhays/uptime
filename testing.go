package uptime

import (
	"net/http/httptest"
	"testing"
)

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get expected status. got %v, want %v", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("Content-Type") != want {
		t.Errorf("wrong content type. expected %s. got %v", want, response.Result().Header)
	}
}

func assertJSONContentType(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	assertContentType(t, response, "application/json")
}
