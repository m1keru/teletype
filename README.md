## Teletype Bot
#### Telegram Bot seeding at group and trying to transcript voice messages to text.

### API Key
* You need to create Bot.
* You need to get Google Cloud API key. [Create here](https://console.cloud.google.com/apis/credentials)
* Copy config/config.yaml.tpl  to config.yaml and edit settings.
### Build
```bash 
make
```

### Install
```
* cp dist/dist.tgz to host
* tar -zxvf dist.tgz
* ./install.sh
* cp google_key.json to /etc/teletype/
* edit /etc/teletype/config.yaml
```
