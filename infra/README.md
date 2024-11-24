# Infra files

## Docker (USE THIS!)

`docker compose up -d` and `docker compose down`

## Systemd service (DEPRECATED)

`platform.service` -> `/etc/systemd/system/platform.service`

Make sure to run `systemctl enable platform.service` and `systemctl start platform.service`

## Nginx configuration

`cyberchallenge.diag.uniroma1.it` -> `/etc/nginx/sites-available/cyberchallenge.diag.uniroma1.it`

Make sure to symlink to it from `sites-enabled`

## Certbot

`sudo snap install --classic certbot`

Initial setup:
`sudo certbot --nginx -d cyberchallenge.diag.uniroma1.it`

It should renew automatically.
