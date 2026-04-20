package base62_test

import (
	"shortener/pkg/base62"
	"testing"
)

func TestInt2String(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		seq  uint64
		want string
	}{
		// TODO: Add test cases.
		{name: "0", seq: 0, want: "0"},
		{name: "1", seq: 1, want: "1"},
		{name: "62", seq: 62, want: "10"},
		{name: "6347", seq: 6347, want: "1En"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base62.Int2String(tt.seq)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("Int2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString2Int(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want uint64
	}{
		// TODO: Add test cases.
		{name: "0", s: "0", want: 0},
		{name: "10", s: "10", want: 62},
		{name: "1En", s: "1En", want: 6347},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base62.String2Int(tt.s)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("String2Int() = %v, want %v", got, tt.want)
			}
		})
	}
}
