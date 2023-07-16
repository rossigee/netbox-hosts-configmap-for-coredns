# Netbox hosts ConfigMap for CoreDNS

This microservice is designed to integrate with a Kubernetes cluster, providing a webhook-triggered mechanism to retrieve a filtered list of IP addresses from Netbox. The IP address data, including hostname and description, is used to generate a typical hosts file. This file is then written to an entry called 'hosts' in a Kubernetes ConfigMap and a success log is generated.

## Requirements

- Go 1.20 or later
- Docker
- Docker Compose
- Kubernetes Cluster
- Netbox instance

## Environment Variables

The application requires several environment variables to run:

- `NETBOX_API_URL`: The URL for your Netbox instance's API.
- `NETBOX_API_TOKEN`: Your Netbox API token.
- `K8S_NAMESPACE`: The Kubernetes namespace to use for the ConfigMap. Defaults to 'kube-system' if not provided.
- `K8S_CONFIGMAP`: The name of the ConfigMap to update. Defaults to 'netbox-hosts' if not provided.

## Building the Application

To build the application, you can use Docker Compose:

```shell
docker-compose build
```

## Testing

To run unit tests:

```shell
go test -v ./...
```

## Running the Application

To run the application:

```shell
docker-compose up
```

## Pushing the Docker Image

To push the Docker image to Dockerhub:

```shell
docker login
docker-compose push
```

# Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

# License

MIT

# The non-automated part of the README...

At home I use Netbox as the Single Source Of Truth for gathering and organising my internal network IP addresses and hostnames. I want to have CoreDNS resolve these hostnames in my cluster with as little manual work as possible. I could use the 'netbox' CoreDNS plugin, which looks cool. However, I don't want to have to build a custom CoreDNS container for my clusters if I can help it, and I want things to continue working even if Netbox is unreachable.

So I figured I could run a small microservice that re-generates a 'hosts' file whenever the IP address data in Netbox changes. It would store it in a ConfigMap which CoreDNS could refer to via the 'hosts' plugin. This should be enough to get me going.

Maybe later I can switch to using the `zone` plugin instead to serve up some useful Netbox data as SRV/TXT fields.

## How did I write this?

To save me writing it and debugging it all on my own from scratch like I normally would, and because it's the end of a long weekend and I'm tired, I enlisted the help of GPT-4, starting with this...

"i need a simple golang microservice for a kubernetes cluster that takes a webhook as a trigger to connect to a netbox instance, retrieve a filtered list of ip addresses with hostname and description, and use them to generate a typical hosts file, which it will write to an entry called 'hosts' in a ConfigMap and log success"

It produced a very tidy initial piece of code that got us started. It was originally based on the netbox library, but we had problems with that so went for straight HTTP requests instead. The Dockerfile, docker-compose, README, github CI config etc. were largely written by GPT-4, and only lightly editted by myself.
