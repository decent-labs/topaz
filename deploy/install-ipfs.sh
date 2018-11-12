#!/usr/bin/env bash

curl https://dist.ipfs.io/go-ipfs/v0.4.17/go-ipfs_v0.4.17_linux-amd64.tar.gz -o go-ipfs.tar.gz
tar xvfz go-ipfs.tar.gz
cd go-ipfs
sudo ./install.sh
cd ..
rm -rf go-ipfs go-ipfs.tar.gz

ipfs init
ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/9001

cat > ipfs.service <<EOL
[Unit]
Description=ipfs
[Service]
ExecStart=/usr/local/bin/ipfs daemon --enable-gc
Restart=always
User=ec2-user
[Install]
WantedBy=multi-user.target
EOL

sudo mv ipfs.service /lib/systemd/system/ipfs.service
sudo systemctl daemon-reload
sudo systemctl enable ipfs
sudo systemctl start ipfs
sudo systemctl status ipfs
