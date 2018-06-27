package analyzer_test

import (
	"fmt"
	"strings"

	"github.com/alex-slynko/haornot/analyzer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func HaveMatchingElement(expected interface{}) types.GomegaMatcher {
	return &matchingElementMatcher{
		expected: expected,
	}
}

type matchingElementMatcher struct {
	expected interface{}
}

func (matcher *matchingElementMatcher) Match(actual interface{}) (success bool, err error) {
	array := actual.([]string)

	if len(array) == 0 {
		return false, nil
	}

	substring := matcher.expected.(string)
	for _, element := range array {
		if strings.Contains(element, substring) {
			return true, nil
		}
	}

	return false, nil
}

func (matcher *matchingElementMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto contain the element that matches\n\t%s", actual, matcher.expected)
}

func (matcher *matchingElementMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to contain the element that matches\n\t%s", actual, matcher.expected)
}

var _ = Describe("Analyze", func() {
	Context("when spec is invalid", func() {
		It("returns error", func() {
			template := []byte(`apiVersion: v1
kind: Deployment
metadata:
  name: nginx
`)
			_, err := analyzer.Analyze(template)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when spec is not for deployment", func() {
		It("returns error", func() {
			template := []byte(`apiVersion: v1
kind: Service
metadata:
  labels:
    name: nginx
  name: nginx
spec:
  ports:
    - port: 80
  selector:
    app: nginx
  type: NodePort
`)
			_, err := analyzer.Analyze(template)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("number of replicas", func() {
		It("is successful for 3 and more replicas when old api format is used", func() {
			template := []byte(`apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
`)
			output, err := analyzer.Analyze(template)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).NotTo(HaveMatchingElement("replicas"))
		})

		It("is successful for 3 and more replicas", func() {
			template := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
`)
			output, err := analyzer.Analyze(template)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).NotTo(HaveMatchingElement("replicas"))
		})

		It("returns message when there are not enough replicas", func() {
			template := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
`)
			output, err := analyzer.Analyze(template)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(HaveMatchingElement("replicas"))
		})

		It("returns message when replicas are not specified", func() {
			template := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
`)
			output, err := analyzer.Analyze(template)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(HaveMatchingElement("replicas"))
		})
	})

	Context("readiness probe", func() {
		It("is successful for 3 and more replicas", func() {
			template := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 1
          timeoutSeconds: 1
`)
			output, err := analyzer.Analyze(template)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).NotTo(HaveMatchingElement("readiness"))
		})

		It("returns message when replicas are not specified", func() {
			template := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
`)
			output, err := analyzer.Analyze(template)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(HaveMatchingElement("readiness"))
		})
	})

})
