# Sonar

Sonar deploys a configurable (if you wish) debugging container to a Kubernetes cluster.

Why *Sonar*? Well it allows you to see deep into your cluster. That, and most nautically-themed names are already taken.

```bash
$ sonar create --image glitchcrab/ubuntu-debug:latest --networkpolicy --podsecuritypolicy \
   --pod-command sleep --pod-args 1h --name glitchcrab-debug --namespace sonar
INFO[0000] serviceaccount "sonar/sonar-glitchcrab-debug" created
INFO[0000] podsecuritypolicy "sonar-glitchcrab-debug" created
INFO[0000] clusterrole "sonar-glitchcrab-debug" created
INFO[0000] clusterrolebinding "sonar-glitchcrab-debug" created
INFO[0000] networkpolicy "sonar-glitchcrab-debug" created
INFO[0000] deployment "sonar/sonar-glitchcrab-debug" created

$ kubectl get po -n sonar
NAME                                      READY   STATUS    RESTARTS   AGE
sonar-glitchcrab-debug-575db85b54-gss4v   1/1     Running   0          43s

$ kubectl exec -it sonar-glitchcrab-debug-575db85b54-gss4v -- bash

notroot@sonar-glitchcrab-debug-575db85b54-gss4v:/$ whoami
notroot

notroot@sonar-glitchcrab-debug-575db85b54-gss4v:/$ hostname
sonar-glitchcrab-debug-575db85b54-gss4v

$ sonar delete --name glitchcrab-debug --namespace sonar --force
INFO[0000] force was set, not asking for confirmation before deleting resources
INFO[0000] deleting deployment
INFO[0000] deleting podsecuritypolicy
INFO[0000] deleting clusterrolebinding
INFO[0000] deleting clusterrole
INFO[0000] deleting networkpolicy
INFO[0000] deleting serviceaccount
INFO[0000] resources deleted: deployment, podsecuritypolicy, clusterrolebinding, clusterrole, networkpolicy, serviceaccount
```

## Configuration

### Global flags

| flag               | default                       | description                                             |
|--------------------|-------------------------------|---------------------------------------------------------|
| `--kubeconfig`     | `/home/$user/.kube/config`    | Absolute path to the kubeconfig file to use.            |
| `--context`        | current context in kubeconfig | Name of the context to use.                             |
| `--name`/`-N`      | `debug`                       | Name given to all resources. Max 50 chars. (see note 1) |
| `--namespace`/`-n` | `default`                     | Namespace to deploy resources to.                       |

#### Notes

1. All names are automatically prepended with `sonar-` for visibility. `--name debug` will result in resources named `sonar-debug`.

### Create

| flag                  | default          | description                                                       |
|-----------------------|------------------|-------------------------------------------------------------------|
| `--image`/`-i`        | `busybox:latest` | Name of the image to use. (see note 1)                            |
| `--networkpolicy`     | `false`          | Creates a NetworkPolicy allowing all ingress & egress.            |
| `--node-exec`         | `null`           | Creates the pod in the host's IPC/net/PID namespaces (see note 2) |
| `--node-name`         | `null`           | Attempt to schedule the pod on the named node.                    |
| `--privileged`        | `false`          | Allow the pod to run as a privileged pod. (see note 3)            |
| `--pod-args`          | `24h`            | Args to pass to the command.                                      |
| `--pod-cmd`           | `sleep`          | Command to use as the entrypoint.                                 |
| `--podsecuritypolicy` | `false`          | Create a PodSecurityPolicy. (see note 4)                          |
| `--pod-userid`        | `1000`           | User ID to run the container as.                                  |

#### Notes

1. If no tag is provided then `latest` is automatically used.
2. A node name to schedule onto must also be provided. Note that the following flags will be ignored: `networkpolicy`, `podsecuritypolicy`, `privileged`.
3. Must be provided at the same time as `--podsecuritypolicy` to have any effect.
4. The PSP will inherit the value set via --pod-userid and configure the minimum value of the RunAs range accordingly.

#### Examples

- `sonar create`
  - accept all defaults. Creates a deployment in namespace `default` called `sonar-debug`.  The pod image will be `busybox:latest` with `sleep 24h` as the initial command.

- `sonar create --image glitchcrab/ubuntu-debug:v1.0 --pod-cmd sleep --pod-args 1h`
  - uses the provided image, command and args.

- `sonar create --podsecuritypolicy --pod-userid 0 --privileged`
  - creates a deployment which runs as root. Also creates a PodSecurityPolicy (and associated RBAC) which allows the pod to run as root/privileged.

- `sonar create --networkpolicy`
  - also creates a NetworkPolicy which allows all ingress and traffic to the Sonar pod.

- `sonar create --context foo-context --namespace bar`
  - create a deployment using context `foo-context` in namespace `bar`.

- `sonar create --node-exec true --node-name worker2 --pod-userid 0`
  - create a pod with root access to the node named `worker2`.

### Delete

| flag      | default | description                                                       |
|-----------|---------|-------------------------------------------------------------------|
| `--force` | `false` | Skips all interaction and deletes all resources created by Sonar. |

#### Examples

- `sonar delete`
  - deletes all resources which match the defaults. This will result in deleting all resources in namespace `default` which are named `sonar-debug`.

- `sonar delete --name test --namespace kube-system`
  - deletes all resources in namespace `kube-system` named `sonar-test`.

## Installing

**Release artifacts**:
- Grab the binary for your distribution from the [releases](https://github.com/glitchcrab/sonar/releases).

**Supported distributions**:
- [Arch Linux](https://github.com/glitchcrab/sonar/tree/main/distribution/arch-linux)

**Build from source**:

```
go build -o sonar
```
