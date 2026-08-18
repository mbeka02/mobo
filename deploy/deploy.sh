#!/bin/bash

# Exit on error
set -e

# Check if commit hash is passed as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <commit-hash>"
  exit 1
fi

COMMIT_HASH=$1
RELEASES_DIR="/home/ubuntu/releases"
DEPLOY_API_BIN="/home/ubuntu/production/mobo"
DEPLOY_WORKER_BIN="/home/ubuntu/production/email-worker"
SERVICE_NAME="mobo"
WORKER_SERVICE="mobo-email-worker"
API_BINARY_NAME="mobo-${COMMIT_HASH}"
WORKER_BINARY_NAME="email-worker-${COMMIT_HASH}"
declare -a PORTS=("3000" "3001" "3002")

# Check if the binaries exist
if [ ! -f "${RELEASES_DIR}/${API_BINARY_NAME}" ]; then
  echo "Binary ${API_BINARY_NAME} not found in ${RELEASES_DIR}"
  exit 1
fi

if [ ! -f "${RELEASES_DIR}/${WORKER_BINARY_NAME}" ]; then
  echo "Binary ${WORKER_BINARY_NAME} not found in ${RELEASES_DIR}"
  exit 1
fi

# Keep a reference to the previous API binary from the symlink
if [ -L "${DEPLOY_API_BIN}" ]; then
  PREVIOUS_API=$(readlink -f $DEPLOY_API_BIN)
  echo "Current API binary is ${PREVIOUS_API}, saved for rollback."
else
  echo "No symbolic link found for API, no previous binary to backup."
  PREVIOUS_API=""
fi

# Keep a reference to the previous email worker binary from the symlink
if [ -L "${DEPLOY_WORKER_BIN}" ]; then
  PREVIOUS_WORKER=$(readlink -f $DEPLOY_WORKER_BIN)
  echo "Current email worker binary is ${PREVIOUS_WORKER}, saved for rollback."
else
  echo "No symbolic link found for email worker, no previous binary to backup."
  PREVIOUS_WORKER=""
fi

rollback_api() {
  if [ -n "$PREVIOUS_API" ]; then
    echo "Rolling back API to previous binary: ${PREVIOUS_API}"
    ln -sfn "${PREVIOUS_API}" "${DEPLOY_API_BIN}"
  else
    echo "No previous API binary to roll back to."
  fi

  # wait to restart the services
  sleep 10

  # Restart all API services with the previous binary
  for port in "${PORTS[@]}"; do
    SERVICE="${SERVICE_NAME}@${port}.service"
    echo "Restarting $SERVICE..."
    sudo systemctl restart $SERVICE
  done

  echo "API rollback completed."
}

rollback_worker() {
  if [ -n "$PREVIOUS_WORKER" ]; then
    echo "Rolling back email worker to previous binary: ${PREVIOUS_WORKER}"
    ln -sfn "${PREVIOUS_WORKER}" "${DEPLOY_WORKER_BIN}"
  else
    echo "No previous email worker binary to roll back to."
  fi

  sleep 5
  echo "Restarting ${WORKER_SERVICE}.service..."
  sudo systemctl restart ${WORKER_SERVICE}.service
  echo "Email worker rollback completed."
}

# --- Deploy API ---

echo "Promoting ${API_BINARY_NAME} to ${DEPLOY_API_BIN}..."
ln -sf "${RELEASES_DIR}/${API_BINARY_NAME}" "${DEPLOY_API_BIN}"

WAIT_TIME=5
restart_service() {
  local port=$1
  local SERVICE="${SERVICE_NAME}@${port}.service"
  echo "Restarting ${SERVICE}..."

  # Restart the service
  if ! sudo systemctl restart "$SERVICE"; then
    echo "Error: Failed to restart ${SERVICE}. Rolling back API deployment."

    # Call the rollback function
    rollback_api
    exit 1
  fi

  # Wait a few seconds to allow the service to fully start
  echo "Waiting for ${SERVICE} to fully start..."
  sleep $WAIT_TIME

  # Check the status of the service
  if ! systemctl is-active --quiet "${SERVICE}"; then
    echo "Error: ${SERVICE} failed to start correctly. Rolling back API deployment."

    # Call the rollback function
    rollback_api
    exit 1
  fi

  echo "${SERVICE}.service restarted successfully."
}

for port in "${PORTS[@]}"; do
  restart_service $port
done

echo "API deployment completed successfully."

# --- Deploy Email Worker ---

echo "Promoting ${WORKER_BINARY_NAME} to ${DEPLOY_WORKER_BIN}..."
ln -sf "${RELEASES_DIR}/${WORKER_BINARY_NAME}" "${DEPLOY_WORKER_BIN}"

echo "Restarting ${WORKER_SERVICE}.service..."
if ! sudo systemctl restart "${WORKER_SERVICE}.service"; then
  echo "Error: Failed to restart ${WORKER_SERVICE}. Rolling back email worker."
  rollback_worker
  exit 1
fi

sleep $WAIT_TIME

if ! systemctl is-active --quiet "${WORKER_SERVICE}.service"; then
  echo "Error: ${WORKER_SERVICE} failed to start correctly. Rolling back email worker."
  rollback_worker
  exit 1
fi

echo "Email worker deployment completed successfully."
echo "Full deployment completed successfully."
