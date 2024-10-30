package porkbun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"woodpecker/internal/config"
	"woodpecker/internal/providers"
)

type Porkbun struct {
	config *config.Config
}

type DNSRecord struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Content string  `json:"content"`
	TTL     string  `json:"ttl"`
	Prio    string  `json:"prio"`
	Notes   *string `json:"notes,omitempty"`
}

type DNSResponse struct {
	Status     string      `json:"status"`
	Cloudflare string      `json:"cloudflare,omitempty"`
	Records    []DNSRecord `json:"records"`
}

func New(config *config.Config) providers.DNSProvider {
	return &Porkbun{config: config}
}

func (p *Porkbun) GetCurrentARecord() (string, error) {
	url := fmt.Sprintf("%s%s/A/%s", p.config.PorkbunRetrieveByNameTypeURL, p.config.PorkbunDomain, p.config.PorkbunSubdomain)

	body := map[string]string{
		"secretapikey": p.config.PorkbunSecretKey,
		"apikey":       p.config.PorkbunAPIKey,
	}

	bodyData, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyData))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve DNS records, status code: %d", resp.StatusCode)
	}

	var dnsResponse DNSResponse
	err = json.NewDecoder(resp.Body).Decode(&dnsResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse DNS response: %v", err)
	}

	if len(dnsResponse.Records) == 0 {
		return "", fmt.Errorf("no A records found for %s.%s", p.config.PorkbunSubdomain, p.config.PorkbunDomain)
	}

	return dnsResponse.Records[0].Content, nil
}

func (p *Porkbun) UpdateARecord(ip string) error {
	url := fmt.Sprintf("%s%s/A/%s", p.config.PorkbunEditByNameTypeURL, p.config.PorkbunDomain, p.config.PorkbunSubdomain)

	body := map[string]string{
		"secretapikey": p.config.PorkbunSecretKey,
		"apikey":       p.config.PorkbunAPIKey,
		"content":      ip,
	}

	bodyData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to parse request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyData))
	if err != nil {
		return fmt.Errorf("failed to send update request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update DNS record, status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to parse DNS update response: %v", err)
	}

	if response["status"] != "SUCCESS" {
		return fmt.Errorf("failed to update DNS record, response: %v", response)
	}

	return nil
}
