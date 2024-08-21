package main

import (
	"fmt"
	"shortlink/internal/common/toolkit"
)

func main() {
	ip := "127.0.0.1"
	location := toolkit.GetLocationByIP(ip)
	fmt.Printf("Location of %s is: %v", ip, location)
}
