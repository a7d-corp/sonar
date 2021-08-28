# Sonar

Sonar deploys a debugging container to a Kubernetes cluster.

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

```
--kube-config (default: '/home/$user/.kube/config')

Absolute path to the kubeconfig file to use.

--context (default: current context in kube config)

Name of the kube context to use.

--name (default: 'debug')

Name given to all the created resources. This will be automatically
prepended with 'sonar-', so a provided name of 'test' will result
in a deployment named 'sonar-debug'. Provided name can be a max of
50 characters.

--namespace (default: 'default')

Namespace to deploy resources to.
```

### Create

#### Options

```
--image (default: 'busybox:latest')

Name of the image to use. Image names may be provided with or without a
tag; if no tag is detected then 'latest' is automatically used.

--pod-cmd (default: 'sleep')

Command to use as the entrypoint.

--pod-args (default: '24h')

Args to pass to the command.

--pod-userid (default: 1000)

User ID to run the container as (set in deployment's SecurityContext).

--podsecuritypolicy (default: false)

Create a PodSecurityPolicy and the associated ClusterRole and Binding.
The PSP will inherit the value set via --pod-userid and configure the
minimum value of the RunAs range accordingly.

--privileged (default: false)

Allow the pod to run as a privileged pod; must be provided at the same
time as --podsecuritypolicy to have any effect.

--networkpolicy (default: false)

Apply a NetworkPolicy which allows all ingress and egress traffic.
```

#### Examples

```
sonar create
```
 - accept all defaults. Creates a deployment in namespace `default` called `sonar-debug`.  The pod image will be `busybox:latest` with `sleep 24h` as the initial command.

```
sonar create --image glitchcrab/ubuntu-debug:v1.0 --pod-cmd sleep \
    --pod-args 1h
```
 - uses the provided image, command and args.

```
sonar create --podsecuritypolicy --pod-userid 0 --privileged
```

- creates a deployment which runs as root. Also creates a PodSecurityPolicy (and associated RBAC) which allows the pod to run as root/privileged.

```
sonar create --networkpolicy
```

 - also creates a NetworkPolicy which allows all ingress and traffic to the Sonar pod.

```
sonar create --context foo-context --namespace bar
```

- create a deployment using context `foo-context` in namespace `bar`.

### Delete

#### Options

```
--force (default: false)

Skips all interaction and deletes all resources created by Sonar.
```

#### Examples

```
sonar delete
```
 - deletes all resources which match the defaults. This will result in deleting all resources in namespace 'default' which are named 'sonar-debug'.

```
sonar delete --name test --namespace kube-system
```
 - deletes all resources in namespace 'kube-system' named 'sonar-test'.

## Installing

Either:
 - Grab the binary for your distribution from the [releases](https://github.com/glitchcrab/sonar/releases).

or:
 - Build it from source:

```
go build -o sonar
```
