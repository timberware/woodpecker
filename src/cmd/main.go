package main

import (
	"fmt"
	"os"
	"woodpecker/src/providers/namecheap"

	"woodpecker/src/config"
	"woodpecker/src/providers/porkbun"
	"woodpecker/src/services"
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

	fmt.Println("checking Porkbun DNS records...")
	dnsProvider := porkbun.New(loadConfig)
	dnsIP, err := dnsProvider.GetCurrentARecord()
	if err != nil {
		fmt.Println("failed to retrieve Porkbun DNS A record:", err)
		os.Exit(1)
	}
	fmt.Println("current Porkbun DNS A record IP address:", dnsIP)

	if dnsIP != ip {
		fmt.Printf("Porkbun DNS record is outdated- current: %s, expected: %s\n", dnsIP, ip)
		err := dnsProvider.UpdateARecord(ip)
		if err != nil {
			fmt.Println("failed to update DNS A record:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Porkbun DNS A record already up to date")
	}

	fmt.Println("updating Namecheap DNS records...")
	dnsProviderNamecheap := namecheap.New(loadConfig)

	err = dnsProviderNamecheap.UpdateARecord(ip)
	if err != nil {
		fmt.Printf("failed to update Namecheap DNS A record: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("DNS updates completed for both providers.")
	os.Exit(0)
}
