# Deployment

## Initial deployment to a VPS / VM
These instructions intended for any cloud provider which offers VM / VPS instances, such as digital ocean or linode.

1. Spin up a Debian-based instance (I used Debian 12 for this example). Then, create a non-root sudo user.
```
adduser bensivo
# follow all the prompts
usermod -aG sudo bensivo
```

2. Install ufw and configure it to only allow http/s and ssh traffic
```
sudo apt install ufw
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
sudo ufw status verbose
```
3. Create 2 DNS records pointing to your new server's public IP
- saladbowl.bensivo.com
- api.saladbowl.bensivo.com

### Deploy the saladbowl service (golang app)
1. Build the go service locally, setting appropriate flags for linux
```
GOOS=linux GOARCH=amd64 go build -o saladbowl-service ./main.go
```
2. Copy the image to your VM
```
IP=x.x.x.x
ssh bensivo@$IP "mkdir -p /home/bensivo/bin"
scp ./saladbowl-service bensivo@$IP:/home/bensivo/bin/
```
3. Create a systemd service file to run the executable
```
sudo vim /etc/systemd/system/saladbowl-service.service

# File contents can be found in the same folder as this document
```

Then, enable the service so it will run automatically on a bootup. And finally, start it
```
sudo systemctl enable saladbowl-service
sudo systemctl start saladbowl-service
```

Your golang application is now deployed to the VM.

### Deploy the saladbowl webapp
1. Install nginx
```
sudo apt-get update
sudo apt-get install nginx
```

2. Build your static site with next.js, then package it into a tar file
```
npm run build
tar -zcvf saladbowl-web.tar ./out

# NOTE: WE have updated next.config.js to make it output a static site that can be served by nginx, instead of by the default next backend
```

3. Copy the tar file to the server and extract it into the 'web' folder
```
IP=saladbowl.bensivo.com
scp ./saladbowl-web.tar bensivo@$IP:/home/bensivo/
ssh bensivo@$IP "mkdir -p /home/bensivo/web && tar -zxvf /home/bensivo/saladbowl-web.tar --strip-components=2 -C /home/bensivo/web"
```

4. Update the default nginx config
```
sudo vim /etc/nginx/nginx.conf

# update the "user" section to use the bensivo user
user bensivo;
```

5. Create and enable configurations for your applications
```
sudo vim /etc/nginx/sites-available/api.saladbowl.bensivo.com
sudo vim /etc/nginx/sites-available/saladbowl.bensivo.com
# File contents can be found in the same folder as this document


sudo ln -s /etc/nginx/sites-available/api.saladbowl.bensivo.com /etc/nginx/sites-enabled/

sudo ln -s /etc/nginx/sites-available/saladbowl.bensivo.com /etc/nginx/sites-enabled/
```

6. Restart nginx
```
sudo nginx -t  # Validates your nginx configurations
sudo systemctl reload nginx
```

## Adding TLS using letsencrypt
Instructions taken from: https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-debian-11
```
sudo apt install certbot python3-certbot-nginx


sudo certbot --nginx -d saladbowl.bensivo.com -d api.saladbowl.bensivo.com
# Follow the instructions. If your DNS is configured to point to this server, it should work successfully. 
```

As a final step, we can verify that the renewal daemon is running
```
sudo systemctl status certbot.timer
sudo certbot renew --dry-run
```


## Updating deployment

After the first installation, updating either app is much quicker.

### Service
```
GOOS=linux GOARCH=amd64 go build -o saladbowl-service ./main.go

IP=saladbowl.bensivo.com
ssh bensivo@$IP "mkdir -p /home/bensivo/bin"
scp ./saladbowl-service bensivo@$IP:/home/bensivo/bin/
ssh -t bensivo@72.14.184.25 "sudo systemctl restart saladbowl-service"
```

### Webapp
```
npm run build
tar -zcvf saladbowl-web.tar ./out

IP=saladbowl.bensivo.com
scp ./saladbowl-web.tar bensivo@$IP:/home/bensivo/
ssh bensivo@$IP "mkdir -p /home/bensivo/web && tar -zxvf /home/bensivo/saladbowl-web.tar --strip-components=2 -C /home/bensivo/web"
```