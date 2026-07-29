package types

// SystemInfo identifies the device. Read-only.
type SystemInfo struct {
	Model    string `json:"model"`
	Firmware string `json:"firmware_version"`
	MAC      string `json:"mac,omitempty"`
	Uptime   int64  `json:"uptime,omitempty"`
}
