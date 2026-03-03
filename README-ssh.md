### Auf dem Host
- *.env Datei auf prod kopieren
  $ scp .env-prod schulenburg-office:/tmp/gowrite-api-go.env
  $ sudo mv /tmp/gowrite-api-go.env /etc/gowrite-api-go.env
  $ sudo chown root:root /etc/gowrite-api-go.env
  $ sudo chmod 600 /etc/gowrite-api-go.env

- *Führe auf dem Host (Server) aus:*  
  $ **sudo systemctl status ssh**

- *Oder um nur zu prüfen, ob Port 22 lauscht:*  
  $ **sudo ss -tlnp | grep :22**

- *Wenn du nur testen willst, ob du dich lokal verbinden kannst:*  
  $ **ssh localhost**

- *Installieren*   
  $ **sudo apt install openssh-server**

- *Starten*   
  $ **sudo systemctl enable --now ssh**

### SSH & Cloudflare
- siehe [README-cloudflare.md](README-cloudflare.md)   
  unter gleicher überschrift : SSH & Cloudflare

