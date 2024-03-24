# Outline Web

Minimalistic web interface for Outline

![screenshot](screenshot.png)

# Usage with docker


### Build
```
docker build -t outline_web .
```


### Run
```
docker run --rm \
  --env OUTLINE_API_URL=%OUTLINE_API_URLS% \
  [--env ADMIN_PASSWORD=%ADMIN_PASSWORD%] \
  [--env PORT=8080] \
  --env ADDR=0.0.0.0 \
  -p 8080 \
  outline_web
```
