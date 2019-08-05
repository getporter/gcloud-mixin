# Play with Google Cloud VMs

This example creates an Google Cloud VM, labels it and then deletes the test VM.

# Credentials

This is what your credentials file should look like, where the path is the path to your the service key for a service account that you have created with the Compute Instance Admin Role.

```yaml
name: gcloud
credentials:
- name: gcloud-key-file
  source:
    path: /Users/carolynvs/Downloads/porter-test-gcloud.json
```

# Try it out

## Create a VM
```console
$ porter install --cred gcloud
```

## Label a VM
```console
$ porter upgrade --cred gcloud
```

## Delete a VM
```console
$ porter uninstall --cred gcloud
```
