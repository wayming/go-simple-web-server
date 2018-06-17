# go-simple-web-server

# Deploy
docker pull wayming/golang-simple-web-server

# Run
docker run -it --rm -p 8080:8080 --name simple-web-server golang-simple-web-server


# Build
docker build -t golang-simple-web-server .
docker tag golang-simple-web-server wayming/golang-simple-web-server
docker push wayming/golang-simple-web-server