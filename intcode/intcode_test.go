package intcode

import (
	"reflect"
	"testing"
)

func Test_decodeOp(t *testing.T) {
	type args struct {
		op int
	}
	tests := []struct {
		name string
		args args
		want OpCode
	}{
		{"test1", args{1002}, OpCode{op: 2, parmModes: [3]int{0, 1, 0}}},
		{"test2", args{11199}, OpCode{op: 99, parmModes: [3]int{1, 1, 1}}},
		{"test3", args{42}, OpCode{op: 42, parmModes: [3]int{}}},
		{"test4", args{10011}, OpCode{op: 11, parmModes: [3]int{0, 0, 1}}},
		{"test5", args{102}, OpCode{op: 2, parmModes: [3]int{1, 0, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeOp(tt.args.op); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeOp() = %v, want %v", got, tt.want)
			}
		})
	}
}
