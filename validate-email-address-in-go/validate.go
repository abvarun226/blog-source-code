package main

import (
	"net"
	"regexp"
	"strings"
)

func validEmail(id string) bool {
	if len(id) > maxEmailLength {
		return false
	}

	if !emailRegex.MatchString(id) {
		return false
	}

	domain := strings.Split(id, "@")[1]
	mx, errLookup := net.LookupMX(domain)
	if errLookup != nil || len(mx) == 0 {
		return false
	}

	return true
}

const maxEmailLength = 254

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
