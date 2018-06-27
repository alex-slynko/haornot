package analyzer_test

import (
	"github.com/alex-slynko/haornot/analyzer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
			Expect(output).To(BeEmpty())
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
			Expect(output).To(BeEmpty())
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
			Expect(output).NotTo(BeEmpty())
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
			Expect(output).NotTo(BeEmpty())
		})
	})

})