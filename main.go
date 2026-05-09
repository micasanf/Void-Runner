package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	dir, _ := os.Getwd()
	fs := http.FileServer(http.Dir(dir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(dir, "void_runner.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	addrs := localIPs()
	fmt.Printf("===== VOID RUNNER server =====\n")
	fmt.Printf("  Local:   http://localhost:%s\n", port)
	for _, a := range addrs {
		fmt.Printf("  Network: http://%s:%s\n", a, port)
	}
	fmt.Printf("==============================\n")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func localIPs() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips
}
