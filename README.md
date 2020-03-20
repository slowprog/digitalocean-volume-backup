## DigitalOcean Volumes Backup

A small container (with GO-script) for automatic backup of DigitalOcean volumes using the [Block Storage API](https://developers.digitalocean.com/documentation/v2/#block-storage).

At each launch, the script creates a snapshot of the volumes and the extra ones are deleted. But you can’t run the script more than once every 10 minutes because you can create a snapshot for volum only once every 10 minutes. This is an API limit. And remember that by default you cannot create more than 25 snapshots per volume.

The following environment variables can be used to configure:

Environment Variable | Description | Default | Required
---|---|---|---
ACCESS_TOKEN | Your DigitalOcean access token. | | +
VOLUMES_BACKUP | The list of volumes for backup (separated by commas). If empty, then all volumes will be backed up. |  |
SNAPSHOTS_PREFIX | Prefix for the name of the snapshots. | backup |
SNAPSHOTS_MAX | The maximum number of stored snapshots for one volume. | 5 |

## Usage

```shell
$ docker run --rm \
    -e ACCESS_TOKEN=$DIGITALOCEAN_ACCESS_TOKEN \
    slowprog/digitalocean-volumes-backup
```

or

```shell
$ docker run --rm \
    -e ACCESS_TOKEN=$DIGITALOCEAN_ACCESS_TOKEN \
    -e SNAPSHOTS_PREFIX=auto \
    -e VOLUMES_BACKUP=postgresql_data,queue_data,monogodb_data \
    -e SNAPSHOTS_MAX=10 \
    slowprog/digitalocean-volumes-backup
```