# Exp👁️ser
exposer is a go tool supported by <a href="https://github.com/projectdiscovery/uncover">uncover</a> to perform query monitoring to different search engines, and storing results in Elasticsearch.

# Dependencies
exposer requires a running Elasticsearch cluster to work properly. Information required to establish communication should be specified in `config.yaml`, or `.env` if running exposer via **docker-compose**.

## Provider Configuration

Uncover requires API keys to the different search engiens to be used. Exposer will not run until one API key is specified at least.
The provider configuration file should be located at `$HOME/.config/uncover/provider-config.yaml`

# Install
exposer requires **go1.21** to install successfully. Run the following command to get the repo -

```sh
go install -v github.com/cheshireca7/exposer@latest
```

## Usage

```sh
exposer -h
```

# Install with docker
exposer has its own image that could be downloaded from Docker Hub

```sh
docker pull cheshireca7/exposer
```
Keys should be specified within the exposer container

```sh
docker run -it exposer exposer vim ~/.config/uncover/provider-config.yaml
```
## Usage

```sh
docker run -it exposer exposer -h
```

## Running Exposer
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
