- cloudflared tunnel route dns schulenburg-office-tunnel schulenburg-office.com
- cloudflared tunnel route dns schulenburg-office-tunnel api1.schulenburg-office.com
- cloudflared tunnel route dns schulenburg-office-tunnel jen1.schulenburg-office.com
- cloudflared tunnel route dns schulenburg-office-tunnel ssh1.schulenburg-office.com
- cloudflared tunnel route dns schulenburg-office-tunnel pga1.schulenburg-office.com
- 
- cloudflared tunnel route dns schulenburg-office-tunnel scriptory.club
- cloudflared tunnel route dns schulenburg-office-tunnel api.scriptory.club
- cloudflared tunnel route dns schulenburg-office-tunnel service.scriptory.club
#### und jeweils
- $ nano /etc/cloudflared/config.yml    
  `- hostname: api1.schulenburg-office.com`    
  `service: http://localhost:8080`

#### cloudflare starten    
  $ sudo cloudflared --config /etc/cloudflared/config.yml tunnel run    

### SSH & Cloudflare
- Gehe zu https://dash.cloudflare.com und wähle deine Domain aus.   
- Klicke links auf „DNS“.   
- Suche den A- oder CNAME-Eintrag deines Servers (z. B. ssh1.schulenburg-office.com).   
- Klicke auf die orange Wolke, sodass sie grau wird – das schaltet auf „DNS only“ um.
- auf dem server rechner
  - $ nano /etc/cloudflared/config.yml    
    `hostname: ssh1.schulenburg-office.com`   
    `service: ssh://localhost:22`
  - $ sudo systemctl restart cloudflared
  - *Bei Problemen mit chatgpt lösen*
  - falls beim ssh einloggen gleich rausfliegt
    - $ mv ~/.bashrc ~/.bashrc.backup   
      mv ~/.profile ~/.profile.backup
    - $ touch ~/.bashrc   
      touch ~/.profile
    - $ ssh localhost
- auf dem client rechner   
  - $ cd /home/ralf
  - $ curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o cloudflared  
  - $ chmod +x cloudflared  
  - $ sudo mv cloudflared /usr/local/bin/  
  - $ cloudflared --version
  - $ cloudflared access ssh --hostname ssh1.schulenburg-office.com

  - $nano ~/.ssh/config
  - Eintrag: (achte auf die spaces)       
    `Host schulenburg-office`   
    `  HostName ssh1.schulenburg-office.com`   
    `  ProxyCommand cloudflared access ssh --hostname %h`   
    `  User serverliebchen`   
    `  Port 22`
  - Dann einfach verbinden mit:   
    $ ssh schulenburg-office


