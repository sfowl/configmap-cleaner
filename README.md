# ConfigMap Cleaner Webhook

A Kubernetes mutating webhook that cleans ConfigMaps by removing private keys.

**USE AT YOUR OWN RISK**


## Deploy

Log in to OpenShift cluster, then:

```bash
$ make deploy IMG=quay.io/sfowler/configmap-cleaner
```

Afterwards all ConfigMap creation/update requests will be checked for private keys, and if found will be removed.

Pre-existing ConfigMap objects that contain private keys will need to be removed manually.


## Monitor

After deployment, monitor the webhook status with

```bash
$ oc logs -n configmap-cleaner -f deploy/configmap-cleaner-manager
```
