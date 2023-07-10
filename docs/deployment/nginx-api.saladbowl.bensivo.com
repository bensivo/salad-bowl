# Save this file at: 
#       /etc/nginx/sites-available/api.saladbowl.bensivo.com
# Then run: 
#       sudo ln -s /etc/nginx/sites-available/api.saladbowl.bensivo.com /etc/nginx/sites-enabled/
server {
    listen 80;
    listen [::]:80;
    server_name api.saladbowl.bensivo.com;

    location / {
	proxy_pass https://api.saladbowl.bensivo.com;
    }

    # Special handling for our websocket connection routes
    location ~ /games/.*/connect {
        proxy_pass https://api.saladbowl.bensivo.com;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}