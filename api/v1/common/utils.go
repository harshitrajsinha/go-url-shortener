package common

import (
	"net"
	"net/http"
	"strings"
)

func convertIPv6LoopbackToIPv4(ipStr string) string {
	trimmed := strings.Trim(ipStr, "[]")
	ip := net.ParseIP(trimmed)
	if ip == nil {
		return ipStr
	}

	if ip.Equal(net.IPv6loopback) {
		return "127.0.0.1"
	}

	return ipStr
}

// Function to check if address is IPv6
func getIPv4Addr(clientIpAddr string) string{
	var networkID string
	if strings.Contains(clientIpAddr, "]"){
		indexVal := strings.Index(clientIpAddr, "]:")
		if indexVal != -1{
			clientIpAddr = clientIpAddr[:indexVal+1]
		}
		networkID = convertIPv6LoopbackToIPv4(clientIpAddr)
	}else{
		networkID = strings.Split(clientIpAddr, ":")[0]
	// 	octet := strings.Split(networkID, ".")
	// 	_ = octet[len(octet)-1]
	}
	return networkID
}

// Function to fetch client's ip address
func GetClientIpAddr(r *http.Request) string {
	var clientIpAddr string
	clientIpAddr = r.Header.Get("X-Forwarded-For")
	if clientIpAddr != "" {
		// The client IP is the first one in the list
		ips := strings.Split(clientIpAddr, ",")
		if len(ips) > 0 {
			clientIpAddr = strings.TrimSpace(ips[0]) // Return the first IP
		}
	} 
	if clientIpAddr == "" {
	// Fallback to X-Real-IP
		clientIpAddr = r.Header.Get("X-Real-IP");
	}
	if clientIpAddr == "" {
		// Fallback to remoteAddr
		clientIpAddr = r.RemoteAddr
	}
	if clientIpAddr != ""{
		// check if address is IPv6
		clientIpAddr = getIPv4Addr(clientIpAddr)
		return clientIpAddr
	}
	return ""
}
