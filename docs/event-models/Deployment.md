# Deployment

## Deploy to a VPS / VM
These instructions intended for any cloud provider which offers VM / VPS instances, such as digital ocean or linode.

1. Spin up a Debian-based instance (used Debian 12 for this example). Then, create a non-root sudo user.
```
adduser bensivo
# follow all the prompts
usermod -aG sudo bensivo
```

### Deploy the saladbowl service
1. Build the go service locally, setting appropriate flags to build for linux
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
sudo vim  /etc/systemd/system/saladbowl-service.service
```

And write these contents to it:
```
[Unit]
Description=Saladbowl service
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=bensivo
ExecStart=/home/bensivo/bin/saladbowl-service

[Install]
WantedBy=multi-user.target
```

Then, enable the service so it will run automatically on a bootup. And finally, start it
```
sudo systemctl enable saladbowl-service
sudo systemctl start saladbowl-service
```

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
```

3. Copy the tar file to the server and extract it into the 'web' folder
```
IP=72.14.184.25
scp ./saladbowl-web.tar bensivo@$IP:/home/bensivo/
ssh bensivo@$IP "mkdir -p /home/bensivo/web && tar -zxvf /home/bensivo/saladbowl-web.tar --strip-components=2 -C /home/bensivo/web"
```

4. Update the nginx config
```
sudo vim /etc/nginx/nginx.conf


# update the "user" section to use the bensivo user
user bensivo;

...

# Add this block within the "http" section, at the end
server {
    listen 80;
    listen [::]:80;
    server_name 72.14.184.25;

    root /home/bensivo/web;

    location / {
    }
}o

```
```
sudo systemctl reload nginx
```


## Updating deployment

After the first installation, updating either app is much quicker.

### Service
```
GOOS=linux GOARCH=amd64 go build -o saladbowl-service ./main.go

IP=x.x.x.x
ssh bensivo@$IP "mkdir -p /home/bensivo/bin"
scp ./saladbowl-service bensivo@$IP:/home/bensivo/bin/
ssh -t bensivo@72.14.184.25 "sudo systemctl restart saladbowl-service"
```

### Webapp
```
npm run build
tar -zcvf saladbowl-web.tar ./out

IP=x.x.x.x
scp ./saladbowl-web.tar bensivo@$IP:/home/bensivo/
ssh bensivo@$IP "mkdir -p /home/bensivo/web && tar -zxvf /home/bensivo/saladbowl-web.tar --strip-components=2 -C /home/bensivo/web"
```