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

type VM struct {
	ID            int
	Pgm           *Program
	Input         chan int
	Output        chan int
	StartingPhase int
}

type Program struct {
	code  []int
	debug bool
}

func NewVM(id int, pgm *Program, startingPhase int, in, out chan int) *VM {
	vm := new(VM)
	vm.ID = id
	vm.Pgm = pgm.Copy()
	vm.Input = in
	vm.Output = out
	vm.StartingPhase = startingPhase
	return vm
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

func (vm *VM) ExecPgm() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ip := 0
	p := vm.Pgm
PGMLOOP:
	for {
		opcode := decodeOp(p.code[ip])
		switch opcode.op {
		case 99:
			if p.debug {
				fmt.Printf("%02d: HALT\n", vm.ID)
			}
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
			var b int
			if vm.StartingPhase >= 0 {
				b = vm.StartingPhase
				vm.StartingPhase = -1
			} else {
				b = <-vm.Input
			}
			p.code[p.code[ip+1]] = b
			if p.debug {
				fmt.Printf("%02d: INPUT:%d\n", vm.ID, b)
			}
			ip += 2
		case 4: // Output
			b := p.code[p.code[ip+1]]
			vm.Output <- b
			if p.debug {
				fmt.Printf("%02d: OUTPUT:%d\n", vm.ID, b)
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
