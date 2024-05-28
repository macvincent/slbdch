#!/bin/bash

# Install Go
sudo apt install -y golang-go

# Download Go tarball
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz -O go.tar.gz

# Extract Go tarball to /usr/local
sudo tar -xzvf go.tar.gz -C /usr/local

# Update PATH in .profile
echo 'export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH' >> ~/.profile

# Source .profile to update current shell session
source ~/.profile

# Navigate to the web_cache directory
cd web_cache

echo "Running Go commands..."

# Run Go commands
go mod tidy
go run main.go
