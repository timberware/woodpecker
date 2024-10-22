package main

import (
	"fmt"
	"os"

	"woodpecker/config"
	"woodpecker/providers/porkbun"
	"woodpecker/services"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("error loading environment file:", err)
		os.Exit(1)
	}

	ip, err := services.GetPublicIP(loadConfig)
	if err != nil {
		fmt.Println("failed to retrieve public IP:", err)
		os.Exit(1)
	}
	fmt.Println("current public IP address:", ip)

	dnsProvider := porkbun.New(loadConfig)

	dnsIP, err := dnsProvider.GetCurrentARecord()
	if err != nil {
		fmt.Println("failed to retrieve DNS A record:", err)
		os.Exit(1)
	}
	fmt.Println("current DNS A record IP address:", dnsIP)

	if dnsIP != ip {
		fmt.Printf("DNS record is outdated- current: %s, expected: %s\n", dnsIP, ip)
		err := dnsProvider.UpdateARecord(ip)
		if err != nil {
			fmt.Println("failed to update DNS A record:", err)
			os.Exit(1)
		}

		fmt.Println("DNS A record updated successfully.")
	}
}
