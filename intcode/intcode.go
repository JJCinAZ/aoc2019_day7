package intcode

import (
	"fmt"
	"strconv"
	"strings"
)

type OpCode struct {
	op        int
	parmModes [3]int
}

type InputBuffer struct {
	buff     []int
	position int
}

type Program struct {
	code   []int
	output []int
	debug  bool
}

func (p *Program) GetOutput() (int, error) {
	if len(p.output) == 0 {
		return 0, fmt.Errorf("EOF")
	}
	i := p.output[0]
	p.output = p.output[1:]
	return i, nil
}

func Compile(input string) *Program {
	pgm := new(Program)
	a := strings.Split(input, ",")
	pgm.code = make([]int, len(a))
	for i := range a {
		pgm.code[i], _ = strconv.Atoi(a[i])
	}
	return pgm
}

func (p *Program) Debug(b bool) {
	p.debug = b
}

func (p *Program) Copy() *Program {
	pgm := new(Program)
	pgm.code = make([]int, len(p.code))
	copy(pgm.code, p.code)
	return pgm
}

func decodeOp(op int) OpCode {
	result := OpCode{}
	result.parmModes[2] = op / 10000
	op = op % 10000
	result.parmModes[1] = op / 1000
	op = op % 1000
	result.parmModes[0] = op / 100
	result.op = op % 100
	return result
}

func (p *Program) ExecPgm(inputBuffer InputBuffer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ip := 0
PGMLOOP:
	for {
		opcode := decodeOp(p.code[ip])
		switch opcode.op {
		case 99:
			break PGMLOOP
		case 1: // Addition
			v1, v2 := getParamsValues(opcode, p.code, ip)
			op3 := p.code[ip+3]
			p.code[op3] = v1 + v2
			ip += 4
		case 2: // Multiplication
			v1, v2 := getParamsValues(opcode, p.code, ip)
			op3 := p.code[ip+3]
			p.code[op3] = v1 * v2
			ip += 4
		case 3: // Input
			p.code[p.code[ip+1]] = inputBuffer.Get()
			if p.debug {
				fmt.Printf("INPUT:%d\n", p.code[p.code[ip+1]])
			}
			ip += 2
		case 4: // Output
			b := p.code[p.code[ip+1]]
			p.output = append(p.output, b)
			if p.debug {
				fmt.Printf("OUTPUT:%d\n", b)
			}
			ip += 2
		case 5: // Jump-if-true
			v1, v2 := getParamsValues(opcode, p.code, ip)
			if v1 != 0 {
				ip = v2
			} else {
				ip += 3
			}
		case 6: // Jump-if-false
			v1, v2 := getParamsValues(opcode, p.code, ip)
			if v1 == 0 {
				ip = v2
			} else {
				ip += 3
			}
		case 7: // Less-than
			v1, v2 := getParamsValues(opcode, p.code, ip)
			op3 := p.code[ip+3]
			if v1 < v2 {
				p.code[op3] = 1
			} else {
				p.code[op3] = 0
			}
			ip += 4
		case 8: // Equals
			v1, v2 := getParamsValues(opcode, p.code, ip)
			op3 := p.code[ip+3]
			if v1 == v2 {
				p.code[op3] = 1
			} else {
				p.code[op3] = 0
			}
			ip += 4
		default:
			panic(fmt.Errorf("illegal opcode at offset %d", ip))
		}
	}
	return nil
}

func getParamsValues(opcode OpCode, pgm []int, ip int) (int, int) {
	v1 := pgm[ip+1]
	if opcode.parmModes[0] == 0 {
		v1 = pgm[v1]
	}
	v2 := pgm[ip+2]
	if opcode.parmModes[1] == 0 {
		v2 = pgm[v2]
	}
	return v1, v2
}

func (buffer *InputBuffer) Push(x int) {
	buffer.buff = append(buffer.buff, x)
}

func (buffer *InputBuffer) Get() int {
	if buffer.position >= len(buffer.buff) {
		panic("EOF")
	}
	x := buffer.buff[buffer.position]
	buffer.position++
	return x
}
