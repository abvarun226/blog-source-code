package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	userIP := "195.176.3.19"
	fmt.Printf("Does `%s` belong to the Tor network? %t\n", userIP, validateTor(userIP))
}

func validateTor(clientIP string) bool {
	torExt := fmt.Sprintf(torDNSEL, reverseOctets(clientIP))

	addrs, _ := net.LookupHost(torExt)
	for _, addr := range addrs {
		if addr == torValidationAddr {
			return true
		}
	}

	return false
}

func reverseOctets(ip string) string {
	oct := strings.Split(ip, ".")
	for i := len(oct)/2 - 1; i >= 0; i-- {
		opp := len(oct) - 1 - i
		oct[i], oct[opp] = oct[opp], oct[i]
	}
	return strings.Join(oct, ".")
}

const (
	torValidationAddr = "127.0.0.2"
	torDNSEL          = "%s.dnsel.torproject.org"
)
