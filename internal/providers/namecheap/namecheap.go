package namecheap

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"woodpecker/internal/config"
	"woodpecker/internal/providers"
)

type Namecheap struct {
	config *config.Config
}

func New(config *config.Config) providers.DNSProvider {
	return &Namecheap{config: config}
}

func (n *Namecheap) GetCurrentARecord() (string, error) {
	return "", fmt.Errorf("GetCurrentARecord not supported in namecheap without API access requirements: https://www.namecheap.com/support/knowledgebase/article.aspx/9739/63/api-faq/#c")
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
		return fmt.Errorf("failed to send update request to Namecheap: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update DNS record, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response from Namecheap: %w", err)
	}

	if strings.Contains(string(body), "<ErrCount>1</ErrCount>") {
		return fmt.Errorf("failed to update Namecheap DNS record: %s", string(body))
	}

	return nil
}
