package kubeadmconfig

import (
	"bytes"
	"text/template"
)

// KubeadmInit defines "kubeadm init" configuration.
// https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/
type KubeadmInit struct {
	MasterNodePrivateDNS string `json:"master-node-private-dns"`
}

// Script returns the service file setup script.
func (ka *KubeadmInit) Script() (s string, err error) {
	return createScriptInit(scriptInit{
		MasterNodePrivateDNS: ka.MasterNodePrivateDNS,
	})
}

func createScriptInit(si scriptInit) (string, error) {
	tpl := template.Must(template.New("scriptInitTmpl").Parse(scriptInitTmpl))
	buf := bytes.NewBuffer(nil)
	if err := tpl.Execute(buf, si); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type scriptInit struct {
	MasterNodePrivateDNS string
}

// make sure to run as root, otherwise "[ERROR IsPrivilegedUser]: user is not running as root".
const scriptInitTmpl = `#!/usr/bin/env bash

sudo touch /home/ec2-user/kubeadm-init.log
sudo mkdir -p /home/ec2-user/.kube
sudo mkdir -p /etc/kubernetes/pki/


cat > /tmp/cluster.yaml <<EOF
apiVersion: kubeadm.k8s.io/v1beta1
kind: InitConfiguration
nodeRegistration:
  name: {{ .MasterNodePrivateDNS }}
  kubeletExtraArgs:
    cloud-provider: aws
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
apiServer:
  extraArgs:
    cloud-provider: aws
controllerManager:
  extraArgs:
    cloud-provider: aws
EOF
cat /tmp/cluster.yaml


sudo kubeadm init --config /tmp/cluster.yaml 1>>/home/ec2-user/kubeadm-init.log 2>&1
`
