package connect

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	convey.Convey("基础用例", t, func() {
		url := "https://www.liwenzhou.com/posts/Go/golang-menu"
		got := Get(url)
		// 断言
		convey.So(got, convey.ShouldEqual, true)
	})
	convey.Convey("url请求不通用例", t, func() {
		url := "posts/Go/url-test-5/"
		got := Get(url)
		// 断言
		//convey.So(got, convey.ShouldEqual, true)
		convey.ShouldBeFalse(got)
	})
}
