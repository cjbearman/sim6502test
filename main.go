package main

import (
	"errors"
	"fmt"
	"os"

	sim6502 "github.com/cjbearman/sim6502/pkg"
)

type EnableDebug struct {
	breakpoint uint16
}

func (ed *EnableDebug) GetBreakpoint() uint16 {
	return ed.breakpoint
}
func (ed *EnableDebug) HandleBreak(proc *sim6502.Processor) error {
	proc.SetOption(sim6502.Trace, true)
	return nil
}

type Success struct {
}

func (s *Success) HandleBreak(proc *sim6502.Processor) error {
	return errors.New("success reached")
}

func main() {
	file, err := os.Open("6502_functional_test.hex")
	if err != nil {
		panic(err)
	}
	proc := sim6502.NewProcessor(&sim6502.RawMemory{})
	proc.SetOption(sim6502.ErrorOnSelfJump, true)

	// proc.SetBreakpoint(&EnableDebug{breakpoint: 0x336d})
	proc.SetBreakpoint(0x3469, &Success{})

	err = proc.LoadHex(file).RunFrom(0x400)
	fmt.Printf("Aborted with error : %v\n", err)
	proc.DumpState(os.Stdout)
	fmt.Printf("Mem in 0x26: 0x%02x\n", proc.Memory().Read(0x0026))
	fmt.Printf("Mem in 0x%02x: 0x%02x\n", proc.Memory().Read(0x0026), proc.Memory().Read(uint16(proc.Memory().Read(0x0026))))
}
