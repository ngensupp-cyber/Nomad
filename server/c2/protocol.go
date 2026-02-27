package c2

import (
	"encoding/json"
)

type PacketType string

const (
	TypeBeacon   PacketType = "BEACON"
	TypeCommand  PacketType = "COMMAND"
	TypeResponse PacketType = "RESPONSE"
)

type Packet struct {
	Type    PacketType      `json:"type"`
	AgentID string          `json:"agent_id"`
	Payload json.RawMessage `json:"payload"`
}

type BeaconPayload struct {
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
}

type CommandPayload struct {
	ID   int    `json:"id"`
	Line string `json:"line"`
}

type ResponsePayload struct {
	CommandID int    `json:"command_id"`
	Output    string `json:"output"`
}
