package createconfig

var defaultConfig = `---
# deployment-specific settings
image: "busybox:latest"
networkpolicy: true
non-root: true
pod-args: "24h"
pod-command: "sleep"
pod-groupid: 1000
pod-userid: 1000
privilege-escalation: false
privileged: false

# global settings
# name: ""
`
