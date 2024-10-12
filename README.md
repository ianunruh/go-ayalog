# go-ayalog

Go library for parsing logs from [Aya](https://aya-rs.dev/) eBPF programs.

This library has no external dependencies.

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
