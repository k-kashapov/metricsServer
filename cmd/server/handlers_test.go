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
	handler := HandleStorage(storage)

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
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Set Gauge #2",
			method: http.MethodPost,
			url:    "/update/gauge/b/-128",
			want: want{
				code:        200,
				response:    "b is set to -128\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Upd Counter #1",
			method: http.MethodPost,
			url:    "/update/counter/e/1",
			want: want{
				code:        200,
				response:    "e is set to 1\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Upd Counter #2",
			method: http.MethodPost,
			url:    "/update/counter/e/8",
			want: want{
				code:        200,
				response:    "e is set to 9\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Upd Counter #3",
			method: http.MethodPost,
			url:    "/update/counter/f/8",
			want: want{
				code:        200,
				response:    "f is set to 8\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Set Counter Insufficient arguments #1",
			method: http.MethodPost,
			url:    "/update/counter/",
			want: want{
				code:        404,
				response:    "Insufficient arguments\n",
				contentType: "",
			},
		},
		{
			name:   "Set Counter Insufficient arguments #2",
			method: http.MethodPost,
			url:    "/update/counter/t",
			want: want{
				code:        404,
				response:    "Insufficient arguments\n",
				contentType: "",
			},
		},
		{
			name:   "Set Gauge Invalid argument",
			method: http.MethodPost,
			url:    "/update/gauge/a/--12",
			want: want{
				code:        400,
				response:    "Invalid argument: --12\n",
				contentType: "",
			},
		},
		{
			name:   "Set Counter Invalid argument",
			method: http.MethodPost,
			url:    "/update/counter/a/--12",
			want: want{
				code:        400,
				response:    "Invalid argument: --12\n",
				contentType: "",
			},
		},
		{
			name:   "Invalid type",
			method: http.MethodPost,
			url:    "/update/eee/a/12",
			want: want{
				code:        400,
				response:    "Invalid type: eee\n",
				contentType: "",
			},
		},
		{
			name:   "Not Post",
			method: http.MethodGet,
			url:    "/update/gauge/a/12",
			want: want{
				code:        400,
				response:    "Request type is not POST\n",
				contentType: "",
			},
		},
		{
			name:   "Set Counter Zero",
			method: http.MethodPost,
			url:    "/update/counter/e/-9",
			want: want{
				code:        200,
				response:    "e is set to 0\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()

			handler(w, req)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
