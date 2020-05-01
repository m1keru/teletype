#!/bin/bash

cp teletype /usr/local/bin/
mkdir -p /etc/teletype
cp config.yaml.tpl /etc/teletype/
cp teletype.service /etc/systemd/system
systemctl enable teletype.service
