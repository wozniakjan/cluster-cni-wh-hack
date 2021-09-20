module github.com/wozniakjan/cluster-cni-wh-hack

go 1.17

replace (
	k8s.io/api => k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.3
	k8s.io/client-go => k8s.io/client-go v0.21.3
	k8s.io/code-generator => k8s.io/code-generator v0.21.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.3
	k8s.io/kubelet => k8s.io/kubelet v0.21.3
	k8s.io/metrics => k8s.io/metrics v0.21.3
)

require (
	github.com/gorilla/mux v1.8.0
	gomodules.xyz/jsonpatch/v2 v2.2.0
	k8c.io/kubermatic/v2 v2.18.0
	k8s.io/klog v1.0.0
)

require (
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.2.1-0.20210504230335-f78f29fc09ea // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kubermatic/machine-controller v1.35.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/open-policy-agent/frameworks/constraint v0.0.0-20210802220920-c000ec35322e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.21.3 // indirect
	k8s.io/apiextensions-apiserver v0.21.3 // indirect
	k8s.io/apimachinery v0.21.3 // indirect
	k8s.io/client-go v12.0.0+incompatible // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/utils v0.0.0-20210722164352-7f3ee0f31471 // indirect
	sigs.k8s.io/controller-runtime v0.9.6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
