# [davisssamuel.net](https://davisssamuel.net)


# Setup the Ubuntu Server

Flash an [Ubuntu Server](https://ubuntu.com/download/server) .iso file to a USB flash drive using [Balena Etcher](https://etcher.balena.io/)

With the flash drive in the desktop, go into the BIOS and changed the boot order to use the USB flash drive first. 
<!-- I had to hold F2 on restart to enter into the BIOS on my system. -->

Using the installation TUI, choose the default options and install the Ubuntu Server onto the desktop.

___
# Setup `ssh` on the server

Confirm the `ssh` service is running:

```
sudo systemctl status ssh
```

If the service is not running, start it:

```
sudo systemctl enable --now ssh
```

---
# Setup the server firewall

Enable the firewall with

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

___
# Setup the Cloudflare Tunnel

Follow the instructions for [creating a locally managed tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/).
NOTE: be sure to follow the commands for Linux when downloading and installing `cloudflared`

After you have authenticated `cloudflared` via the browser popup (step 2), begin creating the tunnel. 

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

For example, based on the above config file's ingress rules:

```
cloudflared tunnel route dns <Tunnel-UUID> example.com ;
cloudflared tunnel route dns <Tunnel-UUID> ssh.example.com
```

Run the tunnel:

```
cloudflared tunnel run <UUID or NAME>
```

___
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

Ensure the `config.yml` file in `~/.cloudflared/` is identical to the one in `/etc/cloudflared/`:

```
diff -ws ~/.cloudflared/config.yml /etc/cloudflared/config.yml
```

If the files are different, copy the `config.yml` from `~/.cloudflared/` to `/etc/cloudflared/`:

```
sudo cp ~/.cloudflared/config.yml /etc/cloudflared/config.yml
```