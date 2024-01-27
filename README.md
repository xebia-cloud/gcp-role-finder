gcp-role-finder - CLI and UI to explore Google Cloud Platform IAM roles
=======================================================================


## Full text role search

```shell
gcp-role-finder search +compute.instances.delete +compute.instances.list
```


## Download roles
Retrieving all the IAM roles takes quite some time. You can also download the
role definitions and search from there.

```shell
$ gcp-role-finder download 
$ gcp-role-finder --from-file search +compute.instances.delete +compute.instances.list
```

## Web interface

You can also run a web user interface. 
```shell
$ docker run -p 8080:8080 ghcr.io/xebia-cloud/gcp-role-finder:latest
```
