# gcloud CLI Mixin for Porter

This is a mixin for Porter that provides the gcloud CLI.

## Mixin Syntax

See the [gcloud CLI Command Reference](https://cloud.google.com/sdk/gcloud/reference/) for the supported commands

```yaml
gcloud:
  description: "Description of the command"
  groups: GROUP
  command: COMMAND
  arguments:
  - arg1
  - arg2
  flags:
    FLAGNAME: FLAGVALUE
    REPEATED_FLAG:
    - FLAGVALUE1
    - FLAGVALUE2
```

```yaml
gcloud:
  description: "Description of the command"
  groups:
  - GROUP 1
  - GROUP 2
  command: COMMAND
  arguments:
  - arg1
  - arg2
  flags:
    FLAGNAME: FLAGVALUE
    REPEATED_FLAG:
    - FLAGVALUE1
    - FLAGVALUE2
```

## Examples

### Authenticate

```yaml
gcloud:
  description: "Authenticate"
  groups:
    - auth
  command: activate-service-account
  flags:
    key-file: gcloud.json
```

### Provision a VM

```yaml
gcloud:
  description: "Create VM"
  groups:
    - compute
    - instances
  command: create
  arguments:
    - porter-test
  flags:
    project: porterci
    zone: us-central1-a
    machine-type: f1-micro
    image: debian-9-stretch-v20190729
    image-project: debian-cloud
    boot-disk-size: 10GB
    boot-disk-type: pd-standard
    boot-disk-device-name: porter-test
```

### Configure SSH Keys

```yaml
gcloud:
  description: "Configure SSH"
  groups: compute
  command: config-ssh
  flags:
    ssh-config-file: ./gce-ssh-config
    ssh-key-file: ./gce-ssh-key
```