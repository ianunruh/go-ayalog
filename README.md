# go-ayalog

Go library for parsing logs from [Aya](https://aya-rs.dev/) eBPF programs.

This library has no external build dependencies.

This library has been tested with logs emitted from an eBPF program using the following library versions.

- [aya-ebpf-log v0.1.0](https://crates.io/crates/aya-log-ebpf/0.1.0)
- [aya-ebpf-log v0.1.1](https://crates.io/crates/aya-log-ebpf/0.1.1)

## Getting started

```bash
go get github.com/ianunruh/go-ayalog
```

### cilium/ebpf

This library has been tested with [cilium/epbf](https://github.com/cilium/ebpf). See `examples/ebpf` for a full example.

```go
// Declare the eBPF map for AYA_LOGS
objs := struct {
    // ... fields for eBPF program and other maps

    AyaLogs *ebpf.Map `ebpf:"AYA_LOGS"`
}{}

// ... normal eBPF setup

// Use the perf reader from cilium/ebpf
reader, err := perf.NewReader(objs.AyaLogs, os.Getpagesize())
if err != nil {
    return err
}

// Parse each message, log it as desired
for {
    record, err := reader.Read()
    if err != nil {
        if errors.Is(err, perf.ErrClosed) {
            return nil
        }
        return err
    }

    r, err := ayalog.ParseRecord(bytes.NewBuffer(record.RawSample))
    if err != nil {
        return err
    }

    spew.Dump(r)
}
```

### Using with aya-log-ebpf v0.1.0

The binary format changed between aya-log-ebpf v0.1.0 and v0.1.1. To use this library with the older binary format, configure the parser.

```go
parser := ayalog.Parser{
    LogLibraryVersion: ayalog.LogLibraryVersion0_1_0,
}
record, err := parser.Record(bytes.NewBuffer(record.RawSample))
// ...
```
