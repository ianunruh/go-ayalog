module github.com/ianunruh/go-ayalog/examples/ebpf

go 1.23.1

replace github.com/ianunruh/go-ayalog => ../../

require (
	github.com/cilium/ebpf v0.16.0
	github.com/davecgh/go-spew v1.1.1
	github.com/ianunruh/go-ayalog v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/exp v0.0.0-20230224173230-c95f2b4c22f2 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
