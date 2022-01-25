#!/bin/bash
set -euxo pipefail

DLE_TEST_MOUNT_DIR="/var/lib/test/dblab_synthetic"
DLE_TEST_POOL_NAME="test_dblab_pool_synthetic"
ZFS_FILE="$(pwd)/zfs_file_synth"
POSTGRES_VERSION="${POSTGRES_VERSION:-13}"

# If previous run was interrupted without cleanup,
# test_dblab_pool and $ZFS_FILE are still here. Cleanup.
#sudo zpool destroy test_dblab_pool_synthetic || true
#sudo rm -f "${ZFS_FILE}"

if [[ ! -d ${ZFS_FILE} ]]; then
  truncate --size 1GB "${ZFS_FILE}"

  sudo zpool create -f \
    -O compression=on \
    -O atime=off \
    -O recordsize=128k \
    -O logbias=throughput \
    -m ${DLE_TEST_MOUNT_DIR} \
    test_dblab_pool_synthetic \
    "${ZFS_FILE}"
fi

# Create dataset
sudo zfs create test_dblab_pool_synthetic/${POSTGRES_VERSION}/${DLE_TEST_POOL_NAME}

sudo zfs list
