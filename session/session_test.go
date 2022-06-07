package sesman

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_getUser(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUser(tt.args.w, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_alreadyLoggedIn(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alreadyLoggedIn(tt.args.w, tt.args.r); got != tt.want {
				t.Errorf("alreadyLoggedIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
