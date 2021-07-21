# sonar

Sonar is a tool for deploying a standalone debugging container to a Kubernetes cluster.

## Configuration

Sonar accepts configuration either via flags or a file. Flags have a higher priority.

Example config file:

```yaml
image: glitchcrab/ubuntu-debug:latest
name: debug
namespace: default
networkpolicy: true
podsecuritypolicy: true
privileged: false
```
