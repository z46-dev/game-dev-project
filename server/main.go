package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/z46-dev/game-dev-project/server/config"
	"github.com/z46-dev/game-dev-project/server/game"
	"github.com/z46-dev/game-dev-project/server/web"
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
	"github.com/z46-dev/golog"
)

var log *golog.Logger = golog.New().Prefix("[MAIN]", golog.BoldBlue).Timestamp()

func discoverTLSKeys(dir string) (certPath, keyPath string, found bool) {
	type Candidate struct {
		cert string
		key  string
	}

	candidates := []Candidate{
		{"fullchain.pem", "privkey.pem"},
		{"cert.pem", "key.pem"},
		{"tls.crt", "tls.key"},
		{"server.crt", "server.key"},
		{"webserver.crt", "webserver.key"},
	}

	for _, c := range candidates {
		certPath = filepath.Join(dir, c.cert)
		keyPath = filepath.Join(dir, c.key)

		if fileExists(certPath) && fileExists(keyPath) {
			return certPath, keyPath, true
		}
	}

	var crtFiles, keyFiles []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", "", false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".crt") {
			crtFiles = append(crtFiles, filepath.Join(dir, name))
		}
		if strings.HasSuffix(name, ".key") {
			keyFiles = append(keyFiles, filepath.Join(dir, name))
		}
	}

	if len(crtFiles) > 0 && len(keyFiles) > 0 {
		return crtFiles[0], keyFiles[0], true
	}

	return "", "", false
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

var g *game.Game = game.NewGame()

func handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	// Set CORS headers
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	socket, err := web.Upgrade(writer, request)

	if err != nil {
		if socket != nil {
			socket.Close()
		}

		return
	}

	var ip string = request.RemoteAddr

	if strings.Contains(ip, "]:") {
		ip = strings.Split(strings.Split(ip, "]:")[0], "[")[1]
	} else if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	socket.Logger = golog.New().Prefix(fmt.Sprintf("[#%d:%s]", socket.ID, ip), golog.BoldGreen).Timestamp()
	socket.Logger.Info("Connection established")

	socket.OnClose = func() {
		socket.Logger.Info("Disconnected")
		game.RemovePlayer(g, socket.ID)
	}

	var username string = request.URL.Query().Get("name")

	if username == "" || len(username) < 3 {
		socket.Logger.Warning("Invalid username")
		var writer *protocol.Writer = new(protocol.Writer)
		writer.SetU8(protocol.PACKET_CLIENTBOUND_KICK)
		writer.SetStringUTF8("Invalid Username")
		socket.Write(writer.GetBytes())
		socket.Close()
		return
	}

	socket.Logger.Prefix(fmt.Sprintf("[#%d:%s:%s]", socket.ID, ip, username), golog.BoldGreen)

	socket.Logger.Info("Joined")

	var player *game.Player = game.NewPlayer(g, socket, username)

	go socket.InitiateUpdateLoop(func(message []byte) {
		if len(message) < 1 {
			return
		}

		var reader *protocol.Reader = protocol.NewReader(message)
		var packetType uint8 = reader.GetU8()

		switch packetType {
		case protocol.PACKET_SERVERBOUND_INPUT:
			if len(message) < 2 {
				return
			}

			var inputFlags uint8 = reader.GetU8()
			player.Body.Control.Goal = util.Vector(0, 0)

			if inputFlags&protocol.BITFLAG_INPUT_UP != 0 {
				player.Body.Control.Goal.Y -= 1
			}

			if inputFlags&protocol.BITFLAG_INPUT_DOWN != 0 {
				player.Body.Control.Goal.Y += 1
			}

			if inputFlags&protocol.BITFLAG_INPUT_LEFT != 0 {
				player.Body.Control.Goal.X -= 1
			}

			if inputFlags&protocol.BITFLAG_INPUT_RIGHT != 0 {
				player.Body.Control.Goal.X += 1
			}

			if inputFlags&protocol.BITFLAG_MOUSE_MOVE != 0 {
				if len(message) < 10 {
					return
				}

				var mouseX float32 = reader.GetF32()
				var mouseY float32 = reader.GetF32()
				player.Body.Control.PrimaryTarget = util.Vector(float64(mouseX), float64(mouseY))
			}
		}
	})
}

func main() {
	var err error
	if err = config.Init("game_server.toml"); err != nil {
		log.Panicf("Could not load configuration: %v", err)
		return
	}

	http.HandleFunc("/ws", handleWebSocket)

	g.Init()
	go g.BeginUpdateLoop(30)

	if config.Config.WebServer.TLSDir != "" {
		log.Info("Starting HTTPS server...")
		var (
			certPath, keyPath string
			found             bool
		)

		if certPath, keyPath, found = discoverTLSKeys(config.Config.WebServer.TLSDir); !found {
			log.Panicf("Could not find TLS keys in %s", config.Config.WebServer.TLSDir)
			return
		}

		http.ListenAndServeTLS(config.Config.WebServer.Address, certPath, keyPath, nil)
	} else {
		log.Info("Starting HTTP server...")
		http.ListenAndServe(config.Config.WebServer.Address, nil)
	}
}
