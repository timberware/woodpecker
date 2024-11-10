package main

import (
	"fmt"
	"os"

	"woodpecker/internal/config"
	"woodpecker/internal/providers/namecheap"
	"woodpecker/internal/providers/porkbun"
	"woodpecker/internal/services"
	"woodpecker/internal/utils"

	"github.com/robfig/cron/v3"
)

func main() {
	utils.InitLogger()

	configDir, err := utils.GetAppPath()
	if err != nil {
		utils.Log.Error().Err(err).Msg("an error occurred")
		os.Exit(1)
	}

	loadConfig, err := config.LoadConfig()
	if err != nil {
		utils.Log.Error().Err(err).Msg("an error occurred")
		os.Exit(1)
	}

	err = updateDNS(loadConfig, configDir)
	if err != nil {
		utils.Log.Error().Err(err).Msg("an error occurred")
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
			utils.Log.Error().Err(err).Msg("error during DNS update: " + err.Error())
		}
	})

	if err != nil {
		utils.Log.Fatal().Err(err).Msg("failed to schedule DNS update")
		os.Exit(1)
	}

	c.Start()
}

func updateDNS(config *config.Config, configPath string) error {
	ip, err := services.GetPublicIP(config)
	if err != nil {
		return err
	}

	storedIP, err := utils.ReadIPFromFile(configPath)
	if err != nil {
		return err
	}

	if storedIP == ip {
		utils.Log.Info().Msgf("public and stored IP address equal (%s), sleeping for %d minute(s)", ip, config.CheckInterval)
		return nil
	}

	utils.Log.Info().Msg("public and stored IP address not equal, proceeding to update DNS records...")

	if config.PorkbunAPIKey != "" && config.PorkbunSecretKey != "" {
		err = updatePorkbunDNS(config, ip)
		if err != nil {
			return err
		}
	}

	if config.NamecheapPassword != "" {
		err = updateNamecheapDNS(config, ip)
		if err != nil {
			return err
		}
	}

	err = utils.WriteIPToFile(ip, configPath)
	if err != nil {
		return err
	}

	utils.Log.Info().Str("level", "update").Msgf("update complete, sleeping for %d minute(s)", config.CheckInterval)
	return nil
}

func updatePorkbunDNS(config *config.Config, ip string) error {
	utils.Log.Info().Msg("checking Porkbun DNS records...")
	porkbunProvider := porkbun.New(config)

	dnsIP, err := porkbunProvider.GetCurrentARecord()
	if err != nil {
		utils.Log.Error().Err(err).Msg("failed to retrieve Porkbun DNS A record")
		return err
	}
	utils.Log.Info().Msg("current Porkbun DNS A record IP address")

	if dnsIP != ip {
		utils.Log.Info().Msg("Porkbun DNS record is outdated")
		err := porkbunProvider.UpdateARecord(ip)
		if err != nil {
			return fmt.Errorf("failed to update Porkbun DNS A record: %w", err)
		}

		utils.Log.Info().Str("level", "update").Msgf("PorkBun (%s.%s) | DNS A record updated successfully.", config.PorkbunSubdomain, config.PorkbunDomain)
	}

	return nil
}

func updateNamecheapDNS(config *config.Config, ip string) error {
	utils.Log.Info().Str("level", "update").Msg("updating Namecheap DNS records...")
	namecheapProvider := namecheap.New(config)

	err := namecheapProvider.UpdateARecord(ip)
	if err != nil {
		return fmt.Errorf("failed to update Namecheap DNS A record: %w", err)
	}

	utils.Log.Info().Str("level", "update").Msgf("Namecheap (%s.%s) | DNS A record updated successfully.", config.NamecheapSubdomain, config.NamecheapDomain)

	return nil
}
