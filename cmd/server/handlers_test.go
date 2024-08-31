package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleStorage(t *testing.T) {
	storage := NewMemStorage()

	ts := httptest.NewServer(NewRouter(storage))
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Set Gauge #1",
			method: http.MethodPost,
			url:    "/update/gauge/a/12",
			want: want{
				code:        200,
				response:    "a is set to 12\n",
				contentType: "text/plain",
			},
		},
		{
			name:   "Set Gauge #2",
			method: http.MethodPost,
			url:    "/update/gauge/b/-128",
			want: want{
				code:        200,
				response:    "b is set to -128\n",
				contentType: "text/plain",
			},
		},
		{
			name:   "Upd Counter #1",
			method: http.MethodPost,
			url:    "/update/counter/e/1",
			want: want{
				code:        200,
				response:    "e is set to 1\n",
				contentType: "text/plain",
			},
		},
		{
			name:   "Upd Counter #2",
			method: http.MethodPost,
			url:    "/update/counter/e/8",
			want: want{
				code:        200,
				response:    "e is set to 9\n",
				contentType: "text/plain",
			},
		},
		{
			name:   "Upd Counter #3",
			method: http.MethodPost,
			url:    "/update/counter/f/8",
			want: want{
				code:        200,
				response:    "f is set to 8\n",
				contentType: "text/plain",
			},
		},
		{
			name:   "Set Counter Insufficient arguments #1",
			method: http.MethodPost,
			url:    "/update/counter/",
			want: want{
				code:        404,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Set Counter Insufficient arguments #2",
			method: http.MethodPost,
			url:    "/update/counter/t",
			want: want{
				code:        404,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Set Gauge Invalid argument",
			method: http.MethodPost,
			url:    "/update/gauge/a/--12",
			want: want{
				code:        400,
				response:    "Invalid argument: --12\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Set Counter Invalid argument",
			method: http.MethodPost,
			url:    "/update/counter/a/--12",
			want: want{
				code:        400,
				response:    "Invalid argument: --12\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Invalid type",
			method: http.MethodPost,
			url:    "/update/eee/a/12",
			want: want{
				code:        400,
				response:    "Invalid type: eee\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Set Counter Zero",
			method: http.MethodPost,
			url:    "/update/counter/e/-9",
			want: want{
				code:        200,
				response:    "e is set to 0\n",
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.url, nil)
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.code, resp.StatusCode)

			defer resp.Body.Close()
			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
