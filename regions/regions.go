package regions

import (
	"strconv"
	"strings"
)

// RegionInfo contains region details
type RegionInfo struct {
	Code    string
	Name    string
	EstPing int // Estimated ping in ms
}

// GetRegionFromIP estimates the region based on IP address
func GetRegionFromIP(addr string) RegionInfo {
	// Extract IP from addr:port format
	ip := strings.Split(addr, ":")[0]
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return RegionInfo{Code: "??", Name: "Unknown", EstPing: 200}
	}

	firstOctet, _ := strconv.Atoi(parts[0])
	secondOctet, _ := strconv.Atoi(parts[1])

	// Common CS2 server ranges
	switch {
	// North America
	case firstOctet == 64 || firstOctet == 65 || firstOctet == 66:
		return RegionInfo{Code: "NA", Name: "North America East", EstPing: 60}
	case firstOctet == 192 && secondOctet == 223:
		return RegionInfo{Code: "NA", Name: "North America West", EstPing: 90}

	// Europe
	case firstOctet == 155 || firstOctet == 162 || firstOctet == 178:
		return RegionInfo{Code: "EU", Name: "Europe West", EstPing: 120}
	case firstOctet == 185 || firstOctet == 186:
		return RegionInfo{Code: "EU", Name: "Europe East", EstPing: 140}

	// Asia
	case firstOctet == 103 || firstOctet == 111 || firstOctet == 112 || firstOctet == 113:
		return RegionInfo{Code: "AS", Name: "Asia", EstPing: 180}

	// Rough estimates based on IP ranges
	case firstOctet >= 1 && firstOctet <= 127:
		return RegionInfo{Code: "NA", Name: "North America", EstPing: 100}
	case firstOctet >= 128 && firstOctet <= 191:
		return RegionInfo{Code: "EU", Name: "Europe", EstPing: 130}
	case firstOctet >= 192 && firstOctet <= 223:
		return RegionInfo{Code: "??", Name: "Unknown", EstPing: 150}
	}

	return RegionInfo{Code: "??", Name: "Unknown", EstPing: 200}
}
