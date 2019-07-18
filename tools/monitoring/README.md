# Knative Monitoring

Knative monitoring is a service that listens to prow job pubsub messages. It scrapes through all the failure logs to catch test infrastructure failures.

## System Diagram

![alt text](systems.png)

## Setup

### Create the Cluster

```bash
gcloud container clusters create monitoring --enable-ip-alias --zone=us-central1-a
```

Note: The cluster connects to the CloudSQL instance via private IP. Thus, it is
required that the cluster is in the same zone as the CloudSQL instance.

## Build and Deploy Changes

### Update the Kubernetes components

[monitoring_service.yaml](https://github.com/knative/test-infra/blob/master/tools/monitoring/gke_deployment/monitoring_service.yaml) is the config to set up all the Kubernetes resources. Use `kubectl apply` on the monitoring_service file to make any updates.

### Update Image

1. `images/monitoring/Makefile` Commands to build and deploy the monitoring images.

1. Update to use the latest image on GKE
    ```bash
    kubectl rollout restart deployment.apps/monitoring-app
    ```

    Check the rollout status
    ```bash
    kubectl rollout status deployment.apps/monitoring-app
    ```

### Clearing the alerts

Run the `clear-alert` kubernetes job. In the monitoring tool directory, run

```bash
kubectl apply -f gke_deployment/clear_alerts_job.yaml
```

Get the pod running the job
```bash
kubectl get pods --selector=job-name=clear-alerts --output=jsonpath='{.items[*].metadata.name}'
```

View the log of the clear-alerts job
```bash
kubectl logs -f <pod-name>
```

Delete the pod after job is completed
```bash
kubectl delete job.batch/clear-alerts
```

