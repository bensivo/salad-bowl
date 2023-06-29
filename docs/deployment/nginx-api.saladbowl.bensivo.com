# Save this file at: 
#       /etc/nginx/sites-available/api.saladbowl.bensivo.com
# Then run: 
#       sudo ln -s /etc/nginx/sites-available/api.saladbowl.bensivo.com /etc/nginx/sites-enabled/
server {
    listen 80;
    listen [::]:80;
    server_name api.saladbowl.bensivo.com;

    location / {
	proxy_pass http://localhost:8080;
    }

    # Special handling for our websocket connection routes
    location ~ /lobbies/.*/connect {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
    }
}