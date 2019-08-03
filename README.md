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

### Provision a VM

```yaml
gcloud:
  description: "Create VM"
  groups:
  - compute
  - instances
  command: create
  arguments:
  - myinst
  flags:
    hostname: "example.com"
    labels: "FOO=BAR,STUFF=THINGS"
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