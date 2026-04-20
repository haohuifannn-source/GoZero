package base62

import (
	"math"
	"strings"
)

// 62进制转换

// 0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ

// 0-9:0-9
// a-z:10-35
// A-Z: 36-62

// 转换公式  ： a(十进制) / 62 = c ......d, c + index(d), 其中c要一直除，得到1/0为止

// 19进制数    转换    62进制数
//    0                  0
//    1                  1
//    10                 a
//    11                 b
//    61                 Z
//    62                 10
//    63                 11
//    6347               ?

// 代码实现转换

//const base62Str = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`

// 为了避免被人恶意请求，我们可以将上面的字符串打乱
// const base62Str = `KLM0149abcdefghijklmno78pqrst56uvwxyzABCDEFGHIJNOPQRS23TUVWXYZ`
// 但是就是测试用例部分得自己去重写算一遍

var (
	base62Str string
)

// MustInit要使用base62这包必须要调用该函数进行初始化
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string!")
	}
	base62Str = bs
}

// Int2String十进制数转换62进制
func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}

	bl := []byte{}
	for seq > 0 {
		mod := seq % 62
		div := seq / 62
		bl = append(bl, []byte(base62Str)[mod])
		seq = div
	}

	// 最后把得到的数据反转
	return string(reverse(bl))

}

// String2Int进制数转换10进制
func String2Int(s string) (seq uint64) {
	bl := []byte(s)
	bl = reverse(bl)
	for idx, b := range bl {
		base := math.Pow(62, float64(idx))
		seq += uint64(strings.Index(base62Str, string(b))) * uint64(base)
	}
	return seq

}

// reverse将字符串转置
func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
