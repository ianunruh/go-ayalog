package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/davecgh/go-spew/spew"

	"github.com/ianunruh/go-ayalog"
)

func main() {
	// Parse command-line arguments
	ifaceName := flag.String("iface", "eth0", "Network interface to attach the program to")
	objPath := flag.String("obj-path", "../../../xdp-hello/target/bpfel-unknown-none/debug/xdp-hello", "Path to built eBPF object file")
	blockIP := flag.String("block-ip", "1.1.1.1", "IPv4 addr to block")
	flag.Parse()

	// Initialize logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load the eBPF object file
	bpfBytes, err := os.ReadFile(*objPath)
	if err != nil {
		log.Fatalf("Failed to read eBPF object file: %v", err)
	}

	// Load the eBPF collection spec
	spec, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(bpfBytes))
	if err != nil {
		log.Fatalf("Failed to load eBPF collection spec: %v", err)
	}

	// override license to GPL
	spec.Programs["xdp_hello"].License = "GPL"

	// Load the eBPF objects into the kernel
	objs := struct {
		XdpHello  *ebpf.Program `ebpf:"xdp_hello"`
		Blocklist *ebpf.Map     `ebpf:"BLOCKLIST"`
		AyaLogs   *ebpf.Map     `ebpf:"AYA_LOGS"`
	}{}
	if err := spec.LoadAndAssign(&objs, nil); err != nil {
		log.Fatalf("Failed to load and assign eBPF objects: %v", err)
	}
	defer objs.XdpHello.Close()
	defer objs.Blocklist.Close()
	defer objs.AyaLogs.Close()

	// Attach the XDP program to the network interface
	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		log.Fatalf("Could not find interface %s: %v", *ifaceName, err)
	}

	lnk, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpHello,
		Interface: iface.Index,
		Flags:     link.XDPGenericMode,
	})
	if err != nil {
		log.Fatalf("Failed to attach XDP program: %v", err)
	}
	defer lnk.Close()

	// Insert IP address into the BLOCKLIST map
	ipAddr := net.ParseIP(*blockIP).To4()
	if ipAddr == nil {
		log.Fatalf("Invalid IP address: %s", *blockIP)
	}
	ipKey := binary.BigEndian.Uint32(ipAddr)
	var zeroValue uint32 = 0

	if err := objs.Blocklist.Put(ipKey, zeroValue); err != nil {
		log.Fatalf("Failed to insert into BLOCKLIST map: %v", err)
	}

	reader, err := perf.NewReader(objs.AyaLogs, os.Getpagesize())
	if err != nil {
		log.Fatalf("Failed to open aya log reader: %v", err)
	}

	var wg sync.WaitGroup

	// start logger
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			record, err := reader.Read()
			if err != nil {
				if errors.Is(err, perf.ErrClosed) {
					return
				}
				log.Fatalf("Failed to read log record: %v", err)
			}

			r, err := ayalog.ParseRecord(bytes.NewBuffer(record.RawSample))
			if err != nil {
				log.Fatalf("Failed to parse log record: %v", err)
			}

			spew.Dump(r)
		}
	}()

	fmt.Println("XDP program attached and BLOCKLIST updated.")
	fmt.Println("Waiting for Ctrl-C to exit...")

	// Set up signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	reader.Close()
	wg.Wait()

	// Cleanup is handled by deferred calls
}
