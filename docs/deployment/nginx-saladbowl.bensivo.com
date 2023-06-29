# Save this file at: 
#       /etc/nginx/sites-available/saladbowl.bensivo.com
# Then run: 
#       sudo ln -s /etc/nginx/sites-available/saladbowl.bensivo.com /etc/nginx/sites-enabled/
server {
    listen 80;
    listen [::]:80;
    server_name saladbowl.bensivo.com;

    root /home/bensivo/web;

    location / {
    }
}