package tests

import (
	"bytes"
	"os"
	"testing"

	sim6502 "github.com/cjbearman/sim6502/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasic6502Operation(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	// Open the source file
	file, err := os.Open("../testcodes/6502_functional_test.hex")
	require.Nil(err, "An error occurred opening the assembled code")

	// Create a new processor, raw memory impl will suffice
	proc := sim6502.NewProcessor(&sim6502.RawMemory{})

	// Set error on self jump, the test code will branch or jump to same instruction
	// in the case of an error, this will catch that
	proc.SetOption(sim6502.ErrorOnSelfJump, true)

	// Address 3469 is the self jump that signals success of the code
	// set a breakpoint here to record the success
	successHandler := &End{t: t}
	proc.SetBreakpoint(0x3469, successHandler)

	// For debugging
	// proc.SetOption(sim6502.Trace, true)

	// Alternatively to turn on debugging at a specific PC
	// proc.SetBreakpoint(<PC VAL>, &EnableDebug{t: t})

	// Run the code
	err = proc.LoadHex(file).RunFrom(0x400)
	assert.Nil(err, "Execution should not have returned an error")

	assert.True(successHandler.success, "Success handler was not called")

	if err != nil || !successHandler.success {
		t.Fail()
		var sw bytes.Buffer
		proc.DumpState(&sw)
		t.Log("State dump:\n" + sw.String())

	}

	executed := uint64(30648048)
	rep := proc.GetLastRunPerformance()
	assert.Equal(executed, rep.InstructionsExecuted, "Expected exactly 30648048 instructions to be executed")
	t.Logf("Last ran for nanos %d cycles %d effective clock: %dMhz", rep.RanForNanoseconds, rep.RanForCycles, rep.EffectiveClock/1000000)
}
