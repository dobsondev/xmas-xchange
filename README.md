# Christmas Gift Exchange Script

## Docker

Below you will find a summary of helpful Docker commands that might need to be used for this project.

### Build the Docker Image

```
docker build -t xmas-xchange .
```

### Run the Docker Container

```
docker run --rm xmas-xchange
```

### Dry Run the Docker Container

```
docker run --rm xmas-xchange --dry-mode
```

### Build the Docker Image with Tag and Push to DockerHub

```
docker build -t xmas-xchange:v0.1 .
docker push dobsondev/xmas-xchange:v0.1
```