package namecheap

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"woodpecker/src/config"
	"woodpecker/src/providers"
)

type Namecheap struct {
	config *config.Config
}

func New(config *config.Config) providers.DNSProvider {
	return &Namecheap{config: config}
}

func (n *Namecheap) GetCurrentARecord() (string, error) {
	return "", fmt.Errorf("GetCurrentARecord not supported in namecheap via API access")
}

func (n *Namecheap) UpdateARecord(ip string) error {
	url := fmt.Sprintf(
		"%s?host=%s&domain=%s&password=%s",
		n.config.NamecheapEditURL,
		n.config.NamecheapSubdomain,
		n.config.NamecheapDomain,
		n.config.NamecheapPassword,
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to send update request to Namecheap: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update DNS record, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response from Namecheap: %v", err)
	}

	if strings.Contains(string(body), "<ErrCount>1</ErrCount>") {
		return fmt.Errorf("failed to update Namecheap DNS record: %s", string(body))
	}

	fmt.Println("Namecheap DNS A record updated successfully")
	return nil
}
