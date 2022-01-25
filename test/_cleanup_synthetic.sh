#!/bin/bash
set -euxo pipefail

# Stop and remove test Docker containers
sudo docker ps -aq --filter label="test_dblab_pool" | xargs --no-run-if-empty sudo docker rm -f
sudo docker ps -aq --filter label="dblab_test" | xargs --no-run-if-empty sudo docker rm -f

# Remove unused Docker images
sudo docker images --filter=reference='registry.gitlab.com/postgres-ai/database-lab/dblab-server:*' -q | xargs --no-run-if-empty sudo docker rmi || true

# To start from the very beginning: destroy ZFS storage pool
#sudo zpool destroy test_dblab_pool_synthetic || true

# Remove CLI configuration
dblab config remove test || true
