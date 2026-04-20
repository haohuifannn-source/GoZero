package urltool_test

import (
	"shortener/pkg/urltool"
	"testing"
)

func TestGetBasePath(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		targeturl string
		want      string
		wantErr   bool
	}{
		{name: "基本实例", targeturl: "https://www.liwenzhou.com/posts/Go/golang-menu", want: "golang-menu", wantErr: false},
		{name: "相对路径url", targeturl: "/xxx/123", want: "", wantErr: true},
		{name: "空字符串", targeturl: "", want: "", wantErr: true},
		{name: "带query的url", targeturl: "https://www.liwenzhou.com/posts/Go/golang-menu?a=1&b=2", want: "golang-menu", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := urltool.GetBasePath(tt.targeturl)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBasePath() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBasePath() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("GetBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
