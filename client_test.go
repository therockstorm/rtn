package rtn

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

type getFileRes struct {
	res string
	err error
}

func TestNewClient(t *testing.T) {
	cj := cookieJarMock{}
	if newClient(cj.new) == nil {
		t.Errorf("Nil client.")
	}
}

func TestNewClientPanic(t *testing.T) {
	cj := cookieJarMock{errors.New("error")}
	defer func() {
		if r := recover(); r != nil {
			s := r.(string)
			if s != "CookieJar required." {
				t.Errorf("Unexpected panic=%v", s)
			}
		} else {
			t.Errorf("Expected panic but wasn't one.")
		}
	}()

	newClient(cj.new)
}

func TestGetFile(t *testing.T) {
	e := errors.New("error")
	const ok = http.StatusOK
	cases := map[string]struct {
		postErr  error
		getErr   error
		readErr  error
		status   int
		expected getFileRes
	}{
		"post error":             {e, nil, nil, ok, getFileRes{"", e}},
		"get error":              {nil, e, nil, ok, getFileRes{"", e}},
		"read error":             {nil, nil, e, ok, getFileRes{"", e}},
		"unexpected status code": {nil, nil, nil, http.StatusInternalServerError, getFileRes{"", errors.New("unexpected statusCode=500, body='rtnInfo'")}},
		"success":                {nil, nil, nil, ok, getFileRes{"rtnInfo", nil}},
	}

	for k, c := range cases {
		r := readerMock{c.readErr}
		cl := client{&httpMock{call{c.status, c.getErr}, call{c.status, c.postErr}}, r.read}
		res, err := cl.getFile()
		if res != c.expected.res || (err != c.expected.err && err.Error() != c.expected.err.Error()) {
			t.Errorf("getFile() == \"%v\", %v, expected \"%v\", %v. %v", res, err, c.expected.res, c.expected.err, k)
		}
	}
}

type readerMock struct {
	err error
}

func (m *readerMock) read(r io.Reader) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return ioutil.ReadAll(r)
}

func nop() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer([]byte("rtnInfo")))
}

type call struct {
	statusCode int
	err        error
}

type httpMock struct {
	get  call
	post call
}

func (m *httpMock) Get(url string) (resp *http.Response, err error) {
	return &http.Response{Body: nop(), StatusCode: m.get.statusCode}, m.get.err
}

func (m *httpMock) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return &http.Response{Body: nop(), StatusCode: m.post.statusCode}, m.post.err
}

type cookieJarMock struct {
	err error
}

func (m *cookieJarMock) new(o *cookiejar.Options) (*cookiejar.Jar, error) {
	if m.err != nil {
		return nil, m.err
	}
	return cookiejar.New(o)
}
