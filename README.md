# Functional tests for sim6502
This is a test package for my sim6502 golang code as found at https://github.com/cjbearman/sim6502

The test package is entirely based on https://github.com/Klaus2m5/6502_65C02_functional_tests/tree/master, which provides an excellent resource for 6502 assembly code for testing purposes.

# Usage
Set the desired version of sim6502 in go.mod and run:

```
go test -v ./...
```

# What it does
The package from Klaus2m5 contains 4 test scripts, consisting of assembled 6502 code.

The scripts typically test various operations, and resolve either by reaching one of many jump points that jump back to themselves, thus resolving in an infinite loop.

Most of these points are failures, one (typically) is success.

The golang tests run these in sim6502, after:
* Enabling a debug option that aborts the processor if an infinite loop occurs (thus detecting failures)
* Sets a breakpoint on the loop that indicates success (thus detecting success)


The four tests:
* Test all regular 6502 opcodes
* Extending testing of BCD arithmatic
* Test extended 65C02 opcodes
* Test interrupts


The interrupt package also requires the use of a mapped memory byte. The test code sets flags in this byte when to indicate the desired states of the IRQ and NMI lines. The memory map handler translates that into changing those simulated lines, thus allowing the test to function.

