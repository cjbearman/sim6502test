module github.com/cjbearman/sim6502test

require github.com/cjbearman/sim6502 v1.0.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cjbearman/sim6502 v1.0.0 => ../sim6502

go 1.20
