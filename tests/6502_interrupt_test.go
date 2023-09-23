package tests

import (
	"bytes"
	"os"
	"testing"

	sim6502 "github.com/cjbearman/sim6502/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This tests IRQ and NMI interrupts

// The test code sets bits 1 at $bffc when it wants the IRQ to go low (asserted)
// and bit 2 when it wants the NMI to go low (assert)
// We handle this by using a memory mapped memory implementation and mapping
// $bffc to a handler

// .. the handler
type interruptFeedbackHandler struct {
	p   *sim6502.Processor
	t   *testing.T
	val uint8
}

// Just mapping $bffc
func (i *interruptFeedbackHandler) AddressRange() []sim6502.MappedMemoryAddressRange {
	return []sim6502.MappedMemoryAddressRange{{Start: 0xbffc, End: 0xbffc}}
}

// Write to $bffc, this is the script telling us to change IRQ/NMI state:
func (i *interruptFeedbackHandler) Write(addr uint16, val uint8) {
	// To help with debugging...
	// i.t.Logf("Feedback register written as $%02x", val)
	if val&0x01 == 0x01 {
		i.p.IRQ(true)
	} else {
		i.p.IRQ(false)
	}

	if val&0x02 == 0x02 {
		i.p.NMI(true)
		// Turn on trace on an NMI, to help with debugging...
		//i.p.SetOption(sim6502.Trace, true)
	} else {
		i.p.NMI(false)
	}

	i.val = val
}

// We'll return the written value if it's ever read
func (i *interruptFeedbackHandler) Read(addr uint16) uint8 {
	return i.val
}

func TestInterrupt6502Operation(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	// Open the source file
	file, err := os.Open("../testcodes/6502_interrupt_test.hex")
	require.Nil(err, "An error occurred opening the assembled code")

	// Create a new processor, we need mappable memory for this
	mem := &sim6502.MappableMemory{}
	proc := sim6502.NewProcessor(mem)
	// If debugging, slow down the clock rate
	// proc.SetClock(100000)

	ifh := &interruptFeedbackHandler{p: proc, t: t}
	mem.Map(ifh)

	// Set error on self jump, the test code will branch or jump to same instruction
	// in the case of an error, this will catch that
	proc.SetOption(sim6502.ErrorOnSelfJump, true)

	// Address 3469 is the self jump that signals success of the code
	// set a breakpoint here to record the success
	successHandler := &End{t: t}
	proc.SetBreakpoint(0x06f5, successHandler)

	// For debugging
	// proc.SetOption(sim6502.Trace, true)
	// proc.SetOption(sim6502.TraceInterrupts, true)

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

	rep := proc.GetLastRunPerformance()
	t.Logf("Last ran for nanos %d cycles %d effective clock: %dMhz", rep.RanForNanoseconds, rep.RanForCycles, rep.EffectiveClock/1000000)
}
