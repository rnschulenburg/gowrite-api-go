- ### .env Datei auf prod kopieren    
  scp .env-prod schulenburg-office:/tmp/gowrite-api-go.env    
  sudo mv /tmp/gowrite-api-go.env /etc/gowrite-api-go.env    
  sudo chown root:root /etc/gowrite-api-go.env    
  sudo chmod 600 /etc/gowrite-api-go.env    
 
- ### service erstellen    
  sudo systemctl daemon-reload    
  sudo systemctl enable gowrite-api-go    
  sudo systemctl start gowrite-api-go    