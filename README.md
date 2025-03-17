# Install the Ubuntu Server

Download [Ubuntu Server](https://ubuntu.com/download/server) and flash the iso file to a USB drive using [Balena Etcher](https://etcher.balena.io/). With the USB drive in the machine, restart, go into the BIOS, and change the boot order to use the USB drive first. Using the installation TUI, choose the default options and install the Ubuntu Server onto the system. Be sure to choose to install OpenSSH when prompted.

After the system reboots, login and update the system.

```
sudo apt update && sudo apt upgrade
```

## Setup ssh on the server

Confirm the ssh service is running.

```
sudo systemctl status ssh
```

If the service is not running, start it.

```
sudo systemctl enable --now ssh
```

## Setup the server firewall

Enable the firewall.

```
sudo ufw enable
```

Allow firewall access to ssh, http, and https.

```
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow http
sudo ufw allow https
```

Confirm the firewall rules were added.

```
sudo ufw status
```

# Setup the Cloudflare Tunnel

Once ssh is setup on the server, you can connect to server from another local client to complete the tunnel setup.

```
ssh <user>@<server_ip>
```

First, follow the instructions for [adding a site to Cloudflare](https://developers.cloudflare.com/fundamentals/setup/account-setup/add-site/). Then, [delete all DNS records](https://developers.cloudflare.com/dns/manage-dns-records/how-to/create-dns-records/#delete-dns-records) from the site in the Cloudfare dashboard. You will assign new DNS records that point traffic to the tunnel.

Follow the instructions for [creating a locally managed tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/) but stop after you have authenticated cloudflared via the browser popup (step 2). Be sure to follow the commands for **Linux** when downloading and installing cloudflared. Then create a tunnel and give it a name.

```
cloudflared tunnel create <tunnel_name>
```

Confirm that the tunnel has been successfully created by running.

```
cloudflared tunnel list
```

Create a configuration file.

```
touch ~/.cloudflared/config.yml
```

```yml
# ~/.cloudflared/config.yml

tunnel: <tunnel_uuid>
credentials-file: /home/<user>/.cloudflared/<tunnel_uuid>.json

ingress:
  - hostname: example.com
    service: http://localhost:<port>
  - hostname: ssh.example.com
    service: ssh://localhost:22
  - service: http_status:404
```

Validate ingress rules.

```
cloudflared tunnel ingress validate
```

Assign a CNAME record that points traffic from your domain/subdomain to the tunnel.

```
cloudflared tunnel route dns <tunnel_uuid or tunnel_name> <hostname>
```

For example, based on the above config file's ingress rules, the command to assign CNAME records would be.

```
cloudflared tunnel route dns <tunnel_uuid> example.com
cloudflared tunnel route dns <tunnel_uuid> ssh.example.com
```

Run the tunnel.

```
cloudflared tunnel run <tunnel_uuid or tunnel_name>
```

## Run cloudflared as a service

In order to ensure the server always connects to the Cloudflare tunnel even on system reboot, a systemd service for cloudflared needs to be installed.

Stop running the tunnel from above with `Ctrl + c`, then install the cloudflared service.

```
cloudflared service install
```

Start the service.

```
sudo systemctl start cloudflared
```

Confirm the service is running.

```
sudo systemctl status cloudflared
```

# Connect to the server remotely

In order to ssh into the server from a remote machine, the remote machine must have [cloudflared installed](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/#1-download-and-install-cloudflared). Once cloudflared has been installed, add the following lines to the remote's `~/.ssh/config` file.

```
# ~/.ssh/config

Host ssh.example.com
    ProxyCommand /usr/local/bin/cloudflared access ssh --hostname %h
```

Attempt to connect to the server using ssh.

```
ssh ssh.example.com
```

## ERROR: Connection closed by UNKNOWN port 65535

If the the `config.yml` file in `~/.cloudflared/` is not identical to the one in `/etc/cloudflared/`, [attempting to ssh will throw an error](https://community.cloudflare.com/t/unable-to-ssh-using-cloudflared/357068)

First, confirm the config files are identical.

```
diff -ws ~/.cloudflared/config.yml /etc/cloudflared/config.yml
```

If the files are different, copy the `config.yml` from `~/.cloudflared/` to `/etc/cloudflared/`.

```
sudo cp ~/.cloudflared/config.yml /etc/cloudflared/config.yml
```

Restart the cloudflared service.

```
sudo systemctl restart cloudflared
```