#!/usr/bin/env bash

sudo yum -y install git

curl -L -o- https://dl.google.com/go/go1.11.linux-amd64.tar.gz | sudo tar -C /usr/local -xz
export PATH=$PATH:/usr/local/go/bin

cd go/src/topaz.io/batch
go get -v .
go build
cd ~

cat > batch.service <<EOL
[Unit]
Description=batch
[Service]
ExecStart=/home/ec2-user/go/src/topaz.io/batch
Restart=always
User=ec2-user
[Install]
WantedBy=multi-user.target
EOL

sudo iptables -A PREROUTING -t nat -p tcp --dport 80 -j REDIRECT --to-ports 8081
sudo mv topaz.service /lib/systemd/system/batch.service
sudo systemctl daemon-reload
sudo systemctl enable batch
sudo systemctl start batch
sudo systemctl status batch
