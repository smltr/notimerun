#!/bin/bash

# Create vendor directory if it doesn't exist
mkdir -p static/vendor

# Function to download a file if it doesn't exist
download_if_not_exists() {
    local url=$1
    local output=$2

    if [ ! -f "$output" ]; then
        echo "Downloading $output..."
        curl -L "$url" -o "$output" || {
            echo "Failed to download $url"
            exit 1
        }
    else
        echo "$output already exists, skipping download"
    fi
}

# Download dependencies
download_if_not_exists "https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js" "static/vendor/htmx.min.js"
download_if_not_exists "https://unpkg.com/alpinejs@3.13.5/dist/cdn.min.js" "static/vendor/alpine.min.js"

# Make sure the files are readable
chmod 644 static/vendor/*.js

echo "Vendor files setup completed successfully!"
