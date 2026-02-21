# Helm Deployment Documentation

## Task 1: Helm Chart Setup and Installation

### Overview

Created Helm charts for both applications (Moscow Time App and Wordle Game) with proper values configuration, image repository settings, and service configurations.

### Helm Charts Created

- **moscow-time-app**: Helm chart for Python Moscow Time application
  - Container Port: 5000
  - Replicas: 3
  - Image: escalopax/moscow-time-app:latest
  - Service Type: NodePort

- **wordle-game**: Helm chart for Go Wordle Game application
  - Container Port: 8080
  - Replicas: 3
  - Image: escalopax/wordle-game:latest
  - Service Type: NodePort

### Application Access

Services are exposed via NodePort and can be accessed using:
```bash
minikube service moscow-time-moscow-time-app
minikube service wordle-wordle-game
```

## Task 2: Helm Chart Hooks

### Hook Implementation

Implemented pre-install and post-install hooks for both applications with hook delete policy.

#### Hook Configuration

- **Pre-install Hook**: Runs before chart installation with weight -5
- **Post-install Hook**: Runs after chart installation with weight 5
- **Delete Policy**: `before-hook-creation` - ensures hooks are deleted before the next hook execution

### Hook Output

#### kubectl get pods,svc Output

```
NAME                                                READY   STATUS      RESTARTS   AGE
pod/moscow-time-app-9c658464-67bt8                  1/1     Running     0          74m
pod/moscow-time-app-9c658464-hxqlh                  1/1     Running     0          74m
pod/moscow-time-app-9c658464-qlvrw                  1/1     Running     0          74m
pod/moscow-time-moscow-time-app-7d768d9bb9-5d6db   1/1     Running     0          15m
pod/moscow-time-moscow-time-app-7d768d9bb9-8wdhq   1/1     Running     0          15m
pod/moscow-time-moscow-time-app-7d768d9bb9-j6bqm   1/1     Running     0          15m
pod/moscow-time-moscow-time-app-post-install-hook  0/1     Completed   0          15m
pod/moscow-time-moscow-time-app-pre-install-hook   0/1     Completed   0          15m
pod/wordle-game-app-79fff8c6fc-56z59               1/1     Running     0          74m
pod/wordle-game-app-79fff8c6fc-sq2dj               1/1     Running     0          74m
pod/wordle-game-app-79fff8c6fc-xvz4d               1/1     Running     0          74m
pod/wordle-wordle-game-7dc4689f45-9vkjg            1/1     Running     0          14m
pod/wordle-wordle-game-7dc4689f45-bs8hg            1/1     Running     0          14m
pod/wordle-wordle-game-7dc4689f45-kjjjx            1/1     Running     0          14m
pod/wordle-wordle-game-post-install-hook           0/1     Completed   0          14m
pod/wordle-wordle-game-pre-install-hook            0/1     Completed   0          14m

NAME                                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
service/kubernetes                    ClusterIP   10.96.0.1        <none>        443/TCP        78m
service/moscow-time-app               NodePort    10.109.192.164   <none>        80:30001/TCP   74m
service/moscow-time-moscow-time-app   NodePort    10.97.28.123     <none>        80:30439/TCP   15m
service/wordle-game-app               NodePort    10.100.138.147   <none>        80:30002/TCP   74m
service/wordle-wordle-game            NodePort    10.111.50.60     <none>        80:32232/TCP   14m
```

#### Pre-install Hook Pod Description

```
Name:             moscow-time-moscow-time-app-pre-install-hook
Namespace:        default
Priority:         0
Service Account:  default
Node:             minikube/192.168.49.2
Start Time:       Sat, 21 Feb 2026 15:47:17 +0300
Labels:           app.kubernetes.io/instance=moscow-time
                  app.kubernetes.io/managed-by=Helm
                  app.kubernetes.io/name=moscow-time-app
                  app.kubernetes.io/version=1.16.0
                  helm.sh/chart=moscow-time-app-0.1.0
Annotations:      helm.sh/hook: pre-install
                  helm.sh/hook-delete-policy: before-hook-creation
                  helm.sh/hook-weight: -5
Status:           Succeeded
IP:               10.244.0.25
IPs:
  IP:  10.244.0.25
Containers:
  pre-install:
    Container ID:  docker://eeea7ddb196dc7d17642b76aaa2451cd274499332d3e878d3ad9a261979d469c
    Image:         busybox
    Image ID:      docker-pullable://busybox@sha256:b3255e7dfbcd10cb367af0d40974d511aeb66dfac98cf30e97e87e4207dd76f
    Port:          <none>
    Host Port:     <none>
    Command:
      sh
      -c
      echo "Running pre-install hook" && sleep 20 && echo "Pre-install hook completed"
    State:          Terminated
      Reason:       Completed
      Exit Code:    0
      Started:      Sat, 21 Feb 2026 15:47:19 +0300
      Finished:     Sat, 21 Feb 2026 15:47:39 +0300
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-dsj92 (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   False
  Initialized                 True
  Ready                       False
  ContainersReady             False
  PodScheduled                True
Volumes:
  kube-api-access-dsj92:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    Optional:                false
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  15m   default-scheduler  Successfully assigned default/moscow-time-moscow-time-app-pre-install-hook to minikube
  Normal  Pulling    15m   kubelet            spec.containers{pre-install}: Pulling image "busybox"
  Normal  Pulled     15m   kubelet            spec.containers{pre-install}: Successfully pulled image "busybox" in 1.399s (1.399s including waiting). Image size: 4105254 bytes.
  Normal  Created    15m   kubelet            spec.containers{pre-install}: Container created
  Normal  Started    15m   kubelet            spec.containers{pre-install}: Container started
```

#### Post-install Hook Pod Description

```
Name:             moscow-time-moscow-time-app-post-install-hook
Namespace:        default
Priority:         0
Service Account:  default
Node:             minikube/192.168.49.2
Start Time:       Sat, 21 Feb 2026 15:47:41 +0300
Labels:           app.kubernetes.io/instance=moscow-time
                  app.kubernetes.io/managed-by=Helm
                  app.kubernetes.io/name=moscow-time-app
                  app.kubernetes.io/version=1.16.0
                  helm.sh/chart=moscow-time-app-0.1.0
Annotations:      helm.sh/hook: post-install
                  helm.sh/hook-delete-policy: before-hook-creation
                  helm.sh/hook-weight: 5
Status:           Succeeded
IP:               10.244.0.28
IPs:
  IP:  10.244.0.28
Containers:
  post-install:
    Container ID:  docker://1c928617cf46df4ef7e9223c5a2785575fcb3fffbe25e6a4639e8f36a098ae58
    Image:         busybox
    Image ID:      docker-pullable://busybox@sha256:b3255e7dfbcd10cb367af0d40974d511aeb66dfac98cf30e97e87e4207dd76f
    Port:          <none>
    Host Port:     <none>
    Command:
      sh
      -c
      echo "Running post-install hook" && sleep 20 && echo "Post-install hook completed"
    State:          Terminated
      Reason:       Completed
      Exit Code:    0
      Started:      Sat, 21 Feb 2026 15:47:42 +0300
      Finished:     Sat, 21 Feb 2026 15:48:03 +0300
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-tqr7w (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   False
  Initialized                 True
  Ready                       False
  ContainersReady             False
  PodScheduled                True
Volumes:
  kube-api-access-tqr7w:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    Optional:                false
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  15m   default-scheduler  Successfully assigned default/moscow-time-moscow-time-app-post-install-hook to minikube
  Normal  Pulling    15m   kubelet            spec.containers{post-install}: Pulling image "busybox"
  Normal  Pulled     15m   kubelet            spec.containers{post-install}: Successfully pulled image "busybox" in 1.179s (1.179s including waiting). Image size: 4105254 bytes.
  Normal  Created    15m   kubelet            spec.containers{post-install}: Container created
  Normal  Started    15m   kubelet            spec.containers{post-install}: Container started
```
