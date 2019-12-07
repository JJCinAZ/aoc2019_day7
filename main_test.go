package main

import (
	"testing"

	"cloud.google.com/aoc2019/day7/intcode"
)

func Test_runAmp(t *testing.T) {
	type args struct {
		pgm    string
		phases [5]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"test1", args{pgm: `3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5`, phases: [5]int{9, 8, 7, 6, 5}}, 139629729},
		{"test2", args{pgm: `3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10`,
			phases: [5]int{9, 7, 8, 5, 6}}, 18216},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgm := intcode.Compile(tt.args.pgm)
			thrust := runAmps(pgm, tt.args.phases)
			if thrust != tt.want {
				t.Errorf("runAmp() = %v, want %v", thrust, tt.want)
			}
		})
	}
}
