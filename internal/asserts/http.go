package asserts

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func StatusCode(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get expected status. want: %d. got: %d", want, got)
	}
}

func ContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("Content-Type") != want {
		t.Errorf("wrong content type. expected: %s. got %v", want, response.Result().Header)
	}
}

func JSONContentType(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	ContentType(t, response, "application/json")
}

func BodyContains(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	body := response.Body.String()
	if !strings.Contains(body, want) {
		t.Errorf("expected body to contain: %s. got %v", want, body)
	}
}
