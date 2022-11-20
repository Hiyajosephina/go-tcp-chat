// Package api here provides methods that both client and server will use
package api

import "regexp"

// CheckHost makes sure the host matches the range for valid IP adresses
func CheckHost(host string) bool {
	match, _ := regexp.MatchString("^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9]"+
		"[0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$", host)
	return match
}

// CheckPort makes sure the host matches the range for valid ports
func CheckPort(port string) bool {
	match, _ := regexp.MatchString("^[0-9]+$", port)
	return match
}
