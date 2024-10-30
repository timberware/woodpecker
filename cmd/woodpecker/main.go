package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"woodpecker/internal/config"
	"woodpecker/internal/providers/namecheap"
	"woodpecker/internal/providers/porkbun"
	"woodpecker/internal/services"
	"woodpecker/internal/utils"
)

func main() {
	configDir, err := utils.GetAppPath()
	if err != nil {
		fmt.Println("failed to get config path:", err)
		os.Exit(1)
	}

	loadConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("error loading environment file:", err)
		os.Exit(1)
	}

	err = updateDNS(loadConfig, configDir)
	if err != nil {
		fmt.Printf("error during initial DNS update: %v\n", err)
		os.Exit(1)
	}

	setupCron(loadConfig, configDir)

	select {}
}

func setupCron(config *config.Config, configPath string) {
	c := cron.New()
	checkInterval := fmt.Sprintf("@every %dm", config.CheckInterval)

	_, err := c.AddFunc(checkInterval, func() {
		err := updateDNS(config, configPath)
		if err != nil {
			fmt.Printf("error during DNS update: %v\n", err)
		}
	})

	if err != nil {
		fmt.Println("failed to schedule DNS update:", err)
		os.Exit(1)
	}

	c.Start()
}

func updateDNS(config *config.Config, configPath string) error {
	ip, err := services.GetPublicIP(config)
	if err != nil {
		return fmt.Errorf("failed to retrieve public IP: %v", err)
	}
	fmt.Println("current public IP address:", ip)

	storedIP, err := utils.ReadIPFromFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read stored IP: %v", err)
	}
	fmt.Println("current stored public IP address:", storedIP)

	if storedIP == ip {
		fmt.Println("IP address unchanged, skipping DNS updates")
		return nil
	}

	fmt.Println("IP address has changed, proceeding to update DNS records...")

	if config.PorkbunAPIKey != "" && config.PorkbunSecretKey != "" {
		err = updatePorkbunDNS(config, ip)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("skipping Porkbun DNS update as required config not provided")
	}

	if config.NamecheapPassword != "" {
		err = updateNamecheapDNS(config, ip)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("skipping Namecheap DNS update as required config not provided")
	}

	err = utils.WriteIPToFile(ip, configPath)
	if err != nil {
		return fmt.Errorf("failed to store new updated public IP address: %v", err)
	}

	fmt.Println("DNS updates completed for both providers.")
	return nil
}

func updatePorkbunDNS(config *config.Config, ip string) error {
	fmt.Println("checking Porkbun DNS records...")
	porkbunProvider := porkbun.New(config)

	dnsIP, err := porkbunProvider.GetCurrentARecord()
	if err != nil {
		return fmt.Errorf("failed to retrieve Porkbun DNS A record: %v", err)
	}
	fmt.Println("current Porkbun DNS A record IP address:", dnsIP)

	if dnsIP != ip {
		fmt.Printf("Porkbun DNS record is outdated- current: %s, expected: %s\n", dnsIP, ip)
		err := porkbunProvider.UpdateARecord(ip)
		if err != nil {
			return fmt.Errorf("failed to update Porkbun DNS A record: %v", err)
		}
		fmt.Println("Porkbun DNS A record updated successfully.")
	} else {
		fmt.Println("Porkbun DNS A record already up to date.")
	}

	return nil
}

func updateNamecheapDNS(config *config.Config, ip string) error {
	fmt.Println("updating Namecheap DNS records...")
	namecheapProvider := namecheap.New(config)

	err := namecheapProvider.UpdateARecord(ip)
	if err != nil {
		return fmt.Errorf("failed to update Namecheap DNS A record: %v", err)
	}
	fmt.Println("Namecheap DNS A record updated successfully.")

	return nil
}
