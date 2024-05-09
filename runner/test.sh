#!/bin/bash

# Set server port (optional, adjust if needed)
PORT=${PORT:-6969}

# Test function with a sample image and arguments
IMAGE="localhost:5000/js_tests"
ENV="EXERCISE=hello.js"
ARGS=""

# Define the folder to be treated as a zip
ZIP="test_data.zip"

# Start the server in the background
# ./main -port=$PORT &
SERVER_PID=$!

# Wait for server to start (adjust sleep time as needed)
# sleep 3

# Send a test request with data
curl -X POST http://localhost:$PORT/$IMAGE?env=$ENV --data-binary "@$ZIP"

# Capture server output
# OUTPUT=$(curl -s http://localhost:$PORT/)

# Stop the server
# kill $SERVER_PID
# wait $SERVER_PID

# echo "$OUTPUT"

