#!/bin/bash

# Variables
IMAGE_TAR="oravue.tar"
IMAGE_NAME="oravue:latest"
CONTAINER_NAME="oravue-container"
PORT="8090" # Change this if needed
CONFIG_PATH="config/prod.yaml"
LOG_FILE="nohup.out"

echo "Starting the deployment process..."

# Step 1: Stop and remove any existing container
EXISTING_CONTAINER=$(sudo docker ps -aq -f name=$CONTAINER_NAME)
if [ -n "$EXISTING_CONTAINER" ]; then
    echo "Stopping and removing existing container: $CONTAINER_NAME"
    sudo docker stop $CONTAINER_NAME
    sudo docker rm $CONTAINER_NAME
else
    echo "No existing container to remove."
fi

# Step 2: Remove the existing image
EXISTING_IMAGE=$(sudo docker images -q $IMAGE_NAME)
if [ -n "$EXISTING_IMAGE" ]; then
    echo "Removing existing image: $IMAGE_NAME"
    sudo docker rmi -f $EXISTING_IMAGE
else
    echo "No existing image to remove."
fi

# Step 3: Load the new Docker image
if [ -f "$IMAGE_TAR" ]; then
    echo "Loading new Docker image from $IMAGE_TAR..."
    sudo docker load -i $IMAGE_TAR
    if [ $? -ne 0 ]; then
        echo "Error: Failed to load Docker image."
        exit 1
    fi
    echo "Docker image loaded successfully."
else
    echo "Error: Docker image tar file $IMAGE_TAR not found."
    exit 1
fi

# Step 4: Verify image was loaded successfully
IMAGE_LOADED=$(sudo docker images -q $IMAGE_NAME)
if [ -z "$IMAGE_LOADED" ]; then
    echo "Error: Docker image $IMAGE_NAME was not loaded correctly."
    exit 1
fi

# Step 5: Check if port is in use
if sudo lsof -i:$PORT -t >/dev/null; then
    echo "Port $PORT is already in use. Please free the port or choose a different one."
    exit 1
fi

# Step 6: Run the container and capture logs in nohup.out
echo "Running the Docker container..."
nohup sudo docker run -d -e CONFIG_PATH=$CONFIG_PATH --platform linux/amd64 -p $PORT:80 --name $CONTAINER_NAME $IMAGE_NAME > $LOG_FILE 2>&1 &

# Step 7: Verify the container is running
sleep 3  # Wait a few seconds for the container to start
if sudo docker ps | grep -q $CONTAINER_NAME; then
    echo "Container $CONTAINER_NAME is running successfully."
    echo "Access your application at http://<VM-PUBLIC-IP>:$PORT"
else
    echo "Error: Container $CONTAINER_NAME is not running."
    exit 1
fi

# Step 8: Remove the tar file after successful deployment
echo "Deleting $IMAGE_TAR..."
rm -f $IMAGE_TAR

echo "Deployment completed successfully!"
