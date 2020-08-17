package main

import (
	sdk "agones.dev/agones/sdks/go"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func gameloop(conn net.PacketConn, stop chan struct{}, s *sdk.SDK) {
	b := make([]byte, 1024)
	for {
		sender, txt := readPacket(conn, b)
		parts := strings.Split(strings.TrimSpace(txt), " ")

		switch parts[0] {
		// shuts down the gameserver
		case "EXIT":
			// respond here, as we os.Exit() before we get to below
			respond(conn, sender, "ACK: "+txt+"\n")
			exit(s)

		// turns off the health pings
		case "UNHEALTHY":
			close(stop)

		case "GAMESERVER":
			respond(conn, sender, gameServerName(s))

		case "READY":
			ready(s)

		case "ALLOCATE":
			allocate(s)

		case "RESERVE":
			if len(parts) != 2 {
				respond(conn, sender, "ERROR: Invalid RESERVE, should have 1 argument\n")
				continue
			}
			if dur, err := time.ParseDuration(parts[1]); err != nil {
				respond(conn, sender, fmt.Sprintf("ERROR: %s\n", err))
				continue
			} else {
				reserve(s, dur)
			}

		case "WATCH":
			watchGameServerEvents(s)

		case "LABEL":
			switch len(parts) {
			case 1:
				// legacy format
				setLabel(s, "timestamp", strconv.FormatInt(time.Now().Unix(), 10))
			case 3:
				setLabel(s, parts[1], parts[2])
			default:
				respond(conn, sender, "ERROR: Invalid LABEL command, must use zero or 2 arguments")
				continue
			}

		case "CRASH":
			log.Print("Crashing.")
			os.Exit(1)

		case "ANNOTATION":
			switch len(parts) {
			case 1:
				// legacy format
				setAnnotation(s, "timestamp", time.Now().UTC().String())
			case 3:
				setAnnotation(s, parts[1], parts[2])
			default:
				respond(conn, sender, "ERROR: Invalid ANNOTATION command, must use zero or 2 arguments\n")
				continue
			}

		case "PLAYER_CONNECT":
			if len(parts) < 2 {
				respond(conn, sender, "ERROR: Invalid PLAYER_CONNECT, should have 1 arguments\n")
				continue
			}
			playerConnect(s, parts[1])

		case "PLAYER_DISCONNECT":
			if len(parts) < 2 {
				respond(conn, sender, "ERROR: Invalid PLAYER_CONNECT, should have 1 arguments\n")
				continue
			}
			playerDisconnect(s, parts[1])

		case "PLAYER_CONNECTED":
			if len(parts) < 2 {
				respond(conn, sender, "ERROR: Invalid PLAYER_CONNECTED, should have 1 arguments\n")
				continue
			}
			respond(conn, sender, playerIsConnected(s, parts[1]))
			continue

		case "GET_PLAYERS":
			respond(conn, sender, getConnectedPlayers(s))
			continue

		case "PLAYER_COUNT":
			respond(conn, sender, getPlayerCount(s))
			continue
		}

		respond(conn, sender, "ACK: "+txt+"\n")
	}
}
