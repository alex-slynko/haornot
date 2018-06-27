package analyzer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alex-slynko/haornot/types"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

const notEnoughReplicasMessage = "At least 2 replicas required for deployment"
const readinessProbeMissingMessage = "Pod %s does not have readiness probe"
const imageVersionMessage = "Image %s for pod %s does not have version. It will always use latest"

var ErrNotADeployment = fmt.Errorf("Not a deployment")

func Analyze(yaml []byte) (*types.Message, error) {
	deployment, err := parseDeployment(yaml)

	if err != nil {
		return nil, err
	}

	msg := types.Message{Name: deployment.Name}

	if deployment.Spec.Replicas == nil {
		msg.Errors = []string{notEnoughReplicasMessage}
		return &msg, nil
	}
	r := *deployment.Spec.Replicas
	if r < 2 {
		msg.Errors = []string{notEnoughReplicasMessage}
		return &msg, nil
	}

	errors := []string{}
	for _, c := range deployment.Spec.Template.Spec.Containers {
		if c.ReadinessProbe == nil {
			errors = append(errors, fmt.Sprintf(readinessProbeMissingMessage, c.Name))
		}

		if !strings.Contains(c.Image, ":") {
			errors = append(errors, fmt.Sprintf(imageVersionMessage, c.Image, c.Name))
		}
	}
	msg.Errors = errors

	return &msg, nil
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
	return &v1.Deployment{}, ErrNotADeployment

}
