package form

import (
	"net/http"
	"strconv"
	"time"
)

func PBool(r *http.Request, fv string) *bool {
	var b *bool

	switch r.Method {
		case http.MethodPost:
			if v, err := strconv.ParseBool(r.PostForm.Get(fv)); err == nil {
				b = &v
			}
		case http.MethodGet:
			if v, err := strconv.ParseBool(r.URL.Query().Get(fv)); err == nil{
				b = &v
			}
	}

	return b
}

func PInt64(r *http.Request, fv string) *int64 {
	var i *int64

	switch r.Method {
		case http.MethodPost, http.MethodPut:
			if v, err := strconv.ParseInt(r.PostForm.Get(fv), 10, 32); err == nil {
				i = &v
			}
		case http.MethodGet:
			if v, err := strconv.ParseInt(r.URL.Query().Get(fv), 10, 32); err == nil{
				i = &v
			}
	}

	return i
}

func VBool(r *http.Request, fv string) bool {
	var b bool

	switch r.Method {
		case http.MethodPost, http.MethodPatch:
			if v, err := strconv.ParseBool(r.PostForm.Get(fv)); err == nil {
				b = v
			}
		case http.MethodGet:
			if v, err := strconv.ParseBool(r.URL.Query().Get(fv)); err == nil{
				b = v
			}
	}

	return b
}

func VDate(r *http.Request, fv string) time.Time {
	var t time.Time

	if v, err := time.Parse(time.DateOnly, VText(r, fv)); err == nil {
		t = v
	}

	return t
}

func VInt(r *http.Request, fv string) int {
	var i int

	switch r.Method {
		case http.MethodPost, http.MethodPatch, http.MethodPut:
			if v, err := strconv.Atoi(r.PostForm.Get(fv)); err == nil {
				i = v
			}
		case http.MethodGet:
			if v, err := strconv.Atoi(r.URL.Query().Get(fv)); err == nil{
				i = v
			}
	}

	return i
}

func VIntArray(r *http.Request, fv string) ([]int, error) {
	var (
		sid []string
		iid []int
	)

	switch r.Method {
		case http.MethodDelete, http.MethodPatch:
			sid = r.Form[fv]
	}

	n   := len(sid)
	iid  = make([]int, n)

	for i, v := range sid {
		id, idErr := strconv.Atoi(v)
		if idErr != nil {
			return nil, idErr
		}
		iid[i] = id 
	}

	return iid, nil
}

func VText(r *http.Request, fv string) string {
	var b string

	switch r.Method {
		case http.MethodPatch, http.MethodPost, http.MethodPut:
			b = r.PostForm.Get(fv)
		case http.MethodGet:
			b = r.URL.Query().Get(fv)
	}

	return b
}

func VTextArray(r *http.Request, fv string) []string {
	var sid []string

	switch r.Method {
		case http.MethodDelete:
			sid = r.Form[fv]
	}

	return sid
}

func VTime(r *http.Request, fv string) time.Time {
	var t time.Time

	if v, err := time.Parse(time.RFC3339Nano, VText(r, fv)); err == nil {
		t = v
	}

	return t
}
