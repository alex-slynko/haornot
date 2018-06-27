package analyzer

import (
	"encoding/json"
	"fmt"

	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

const notEnoughReplicas = "At least 2 replicas required for deployment"
const readinessProbeMissing = "Pod %s does not have readiness probe"

func Analyze(yaml []byte) ([]string, error) {
	deployment, err := parseDeployment(yaml)

	if err != nil {
		return []string{}, err
	}

	if deployment.Spec.Replicas == nil {
		return []string{notEnoughReplicas}, nil
	}
	r := *deployment.Spec.Replicas
	if r < 2 {
		return []string{notEnoughReplicas}, nil
	}

	errors := []string{}
	for _, c := range deployment.Spec.Template.Spec.Containers {
		if c.ReadinessProbe == nil {
			errors = append(errors, fmt.Sprintf(readinessProbeMissing, c.Name))
		}
	}

	return errors, nil
}

func parseDeployment(yaml []byte) (*v1.Deployment, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode

	gvk := &schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	obj, meta, err := decode([]byte(yaml), gvk, nil)
	if err != nil {
		return &v1.Deployment{}, err
	}

	if meta.Kind == "Deployment" {
		b, err := json.Marshal(obj)
		if err != nil {
			return &v1.Deployment{}, err
		}
		deployment := &v1.Deployment{}
		err = json.Unmarshal(b, deployment)
		return deployment, err
	}
	return &v1.Deployment{}, fmt.Errorf("only deployment can be analyzed")

}
