package providers

type DNSProvider interface {
	GetCurrentARecord() (string, error)
	UpdateARecord(ip string) error
}
