# Troubleshooting

## Node scale up failed: Pod is at risk of not being scheduled

If this occurs on GKE Autopilot, see
[this GCP doc](https://cloud.google.com/kubernetes-engine/docs/troubleshooting/troubleshooting-autopilot-clusters#scale-up-failed-serial-port-logging)
which explains that you need to enable serial port logging either in your
organization or project, which can be done by running:

```
gcloud compute project-info add-metadata \
    --metadata serial-port-logging-enable=true
```
