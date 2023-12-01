# Install the Ubuntu Server

Download [Ubuntu Server](https://ubuntu.com/download/server) and flash the `.iso` file to a USB drive using [Balena Etcher](https://etcher.balena.io/). With the USB drive on the system, go into the BIOS and changed the boot order to use the USB drive first. Using the installation TUI, choose the default options and install the Ubuntu Server onto the system.
<!-- I had to hold F2 on restart to enter into the BIOS on my system. --> 

NOTE: be sure to choose to install OpenSSH when prompted.

After the system reboots, login and update the system:

```
sudo apt update && sudo apt upgrade
```

# Setup `ssh` on the server

Confirm the `ssh` service is running:

```
sudo systemctl status ssh
```

If the service is not running, start it:

```
sudo systemctl enable --now ssh
```

## Setup the server firewall

Enable the firewall:

```
sudo ufw enable
```

Allow firewall access to `ssh`, `http`, and `https`:

```
sudo ufw allow ssh ;
sudo ufw allow http ;
sudo ufw allow https
```

Confirm the firewall rules were added:

```
sudo ufw status
```

Output:

```
To                         Action      From
--                         ------      ----
22/tcp                     ALLOW       Anywhere                  
80/tcp                     ALLOW       Anywhere                  
443                        ALLOW       Anywhere                  
22/tcp (v6)                ALLOW       Anywhere (v6)             
80/tcp (v6)                ALLOW       Anywhere (v6)             
443 (v6)                   ALLOW       Anywhere (v6)
```

# Setup the Cloudflare Tunnel

Follow the instructions for [creating a locally managed tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/). After you have authenticated `cloudflared` via the browser popup (step 2), begin creating the tunnel.

NOTE: be sure to follow the commands for **Linux** when downloading and installing `cloudflared`.



Create a tunnel and give it a name: 

```
cloudflared tunnel create <NAME>
```

Confirm that the tunnel has been successfully created by running:

```
cloudflared tunnel list
```

Create a configuration file:

```
touch ~/.cloudflared/config.yml
```

```yml
# config.yml

tunnel: <Tunnel-UUID>
credentials-file: /home/sam/.cloudflared/<Tunnel-UUID>.json

ingress:
  - hostname: example.com
    service: http://localhost:8000
  - hostname: ssh.example.com
    service: ssh://localhost:22
  - service: http_status:404
```

Validate ingress rules:

```
cloudflared tunnel ingress validate
```

Assign a `CNAME` record that points traffic to your tunnel domain/subdomain:

```
cloudflared tunnel route dns <UUID or NAME> <hostname>
```

For example, based on the above config file's ingress rules, the command to assign `CNAME` records would be:

```
cloudflared tunnel route dns <Tunnel-UUID> example.com ;
cloudflared tunnel route dns <Tunnel-UUID> ssh.example.com
```

Run the tunnel:

```
cloudflared tunnel run <UUID or NAME>
```

# Run `cloudflared` as a service

In order to make the ensure the server always connects to the Cloudflare tunnel even on system reboot, `cloudflared` needs to be as a service. 

Stop running the tunnel from above with `Ctrl + c`

Install `cloudflared` service:

```
cloudflared service install
```

Start the service:

```
systemctl start cloudflared
```

Confirm the service is running:

```
systemctl status cloudflared
```

If the the `config.yml` file in `~/.cloudflared/` is not identical to the one in `/etc/cloudflared/`, connecting to the server with `ssh` will throw:

```
kex_exchange_identification: Connection closed by remote host
Connection closed by UNKNOWN port 65535
```

Confirm the config files are identical:

```
diff -ws ~/.cloudflared/config.yml /etc/cloudflared/config.yml
```

If the files are different, copy the `config.yml` from `~/.cloudflared/` to `/etc/cloudflared/`:

```
sudo cp ~/.cloudflared/config.yml /etc/cloudflared/config.yml
```