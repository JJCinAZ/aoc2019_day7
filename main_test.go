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
		{"test1", args{pgm: `3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0`, phases: [5]int{4, 3, 2, 1, 0}}, 43210},
		{"test2", args{pgm: `3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0`, phases: [5]int{0, 1, 2, 3, 4}}, 54321},
		{"test3", args{pgm: `3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0`, phases: [5]int{1, 0, 4, 3, 2}}, 65210},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgm := intcode.Compile(tt.args.pgm)
			thrust := 0
			thrust = runAmp(pgm, tt.args.phases[0], thrust)
			thrust = runAmp(pgm, tt.args.phases[1], thrust)
			thrust = runAmp(pgm, tt.args.phases[2], thrust)
			thrust = runAmp(pgm, tt.args.phases[3], thrust)
			thrust = runAmp(pgm, tt.args.phases[4], thrust)
			if thrust != tt.want {
				t.Errorf("runAmp() = %v, want %v", thrust, tt.want)
			}
		})
	}
}
