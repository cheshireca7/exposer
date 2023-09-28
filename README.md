<h1 align="center">Exp👁️ser</h1>
<h4 align="center">exposer is a go tool supported by <a href="https://github.com/projectdiscovery/uncover">uncover</a> to perform query monitoring to different search engines, and storing results in Elasticsearch.</h4>

<p align="center">
  <a href="#dependencies">Dependencies</a> •
  <a href="#installation-instructions">Installation</a> •
  <a href="#installation-with-docker">Installation with Docker</a> •
  <a href="#usage">Usage</a> •
  <a href="#running-exposer">Running Exposer</a>
</p>

---

# Dependencies
exposer requires a running Elasticsearch cluster to work properly. Information required to establish communication should be specified in `config.yaml`, or `.env` if running exposer via **docker-compose**.

## Provider Configuration

Uncover requires API keys to the different search engiens to be used. Exposer will not run until one API key is specified at least.
The provider configuration file should be located at `$HOME/.config/uncover/provider-config.yaml`

# Installation Instructions
exposer requires **go1.21** to install successfully. Run the following command to get the repo -

```sh
go install -v github.com/cheshireca7/exposer@latest
```

### Before runnning
1. Edit `$HOME/.config/uncover/provider-config.yaml` with API keys for search engines.
2. Edit `$HOME/.config/exposer/config.yaml` file with the data regarding Elasticsearch communication. As an example:

```yaml
URL: localhost
PORT: 9200
USERNAME: elastic
PASSWORD: elastic
```

3. Set the Elasticsearch CA to be at `$HOME/.config/exposer/http_ca.crt`

## Docker
exposer has its own image that could be downloaded from Docker Hub

```sh
docker pull cheshireca7/exposer
```

### Before runnning
1. Edit `$HOME/.config/uncover/provider-config.yaml` with API keys for search engines.
```sh
docker run -it exposer vim ~/.config/uncover/provider-config.yaml
```

2. Edit `$HOME/.config/exposer/config.yaml` file with the data regarding Elasticsearch communication.
```sh
docker run -it exposer vim ~/.config/exposer/config.yaml
```

3. Get the certificate from the elasticsearch container and upload it to the exposer container
```sh
docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt exposer:/root/.config/exposer/http_ca.crt
```

## Docker compose
By running docker-compose, it will load a clear elasticsearch container, as well as exposer at once. 

1. Credentials for Elasticsearch communication should be set at `docker/.env` file, then run `docker-compose up -d`
2. Edit `$HOME/.config/uncover/provider-config.yaml` with API keys for search engines.

# Usage

```sh
exposer -h
```

### Docker 

```sh
docker run -it exposer exposer -h
```

# Running Exposer
Default run just require a query

```console
exposer -q 'ssl:hackerone.com'

                                                                                      
                                                                                      
 ,adPPYba,  8b,     ,d8  8b,dPPYba,    ,adPPYba,   ,adPPYba,   ,adPPYba,  8b,dPPYba,  
a8P_____88   `Y8, ,8P'   88P'    "8a  a8"     "8a  I8[    ""  a8P_____88  88P'   "Y8  
8PP"""""""     )888(     88       d8  8b       d8   `"Y8ba,   8PP"""""""  88          
"8b,   ,aa   ,d8" "8b,   88b,   ,a8"  "8a,   ,a8"  aa    ]8I  "8b,   ,aa  88          
 `"Ybbd8"'  8P'     `Y8  88`YbbdP"'    `"YbbdP"'   `"YbbdP"'   `"Ybbd8"'  88          
                         88                                                           
                         88                                                           

-- Monitor your favorite services exposed to the Internet 👀


[INF] Creating new index: 2023-09-26-13-56-10_uncover_results
[INF] Monitoring query: 'ssl:hackerone.com'
[INF] Number of entries stored: 4

```
# TODO

[] Interactive console to get more information about stored results
[] Improve installation
