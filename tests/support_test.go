package tests

import (
	"testing"

	sim6502 "github.com/cjbearman/sim6502/pkg"
)

// EnableDebug can be used to enable tracing during execution
// Use it as a break point handler at the PC at which you want tracing
// to be enabled
type EnableDebug struct {
	t *testing.T
}

func (ed *EnableDebug) HandleBreak(proc *sim6502.Processor) error {
	ed.t.Log("Enabling debug")
	proc.SetOption(sim6502.Trace, true)
	return nil
}

// Success is a breakpoint handler that should be set at the code location (PC)
// who's execution indicates success of the test
// It will stop the processor and record the success
type Success struct {
	t       *testing.T
	success bool
}

func (s *Success) HandleBreak(proc *sim6502.Processor) error {
	proc.Stop()
	s.success = true
	s.t.Log("Success handler called")
	return nil
}
