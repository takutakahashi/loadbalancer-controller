### Backend

Backend resource has responsibilities for building loadbalancer using cloud provider, gathering access info to reach an app.

Once Backend resource created, controller start below reconcile sequence.

- get credentials info from configmap if needed.
- build loadbalancer using cloud provider (ex: AWS)
- gather loadbalancer access info and register own status

Backend implementations can be added as a CRD.
