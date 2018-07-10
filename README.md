# go-simple-web-server

# Install docker on aws linux
./install_aws_linux.sh
Logout and login

# Deploy
docker pull wayming/golang-simple-web-server

# Run
docker run -dit --restart unless-stopped --rm -p 8080:8080 --name simple-web-server wayming/golang-simple-web-server 

# Build
docker build -t golang-simple-web-server .
docker tag golang-simple-web-server wayming/golang-simple-web-server
docker push wayming/golang-simple-web-server
