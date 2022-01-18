package v1alpha1

import (
	v1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubernetes Helpers", func() {
	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Helpers functions", func() {
		It("should return a correct volume", func() {
			vol := getVolumeSpec("vol1", v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "cfg1",
					},
				},
			})

			Expect(vol).ToNot(BeNil())
			Expect(vol.Name).To(Equal("vol1"))
			Expect(vol.ConfigMap.Name).To(Equal("cfg1"))
		})

		It("should return a correct volume from a configmap", func() {
			vol := getVolumeSpecFromConfigMap("vol1", "cfg1")

			Expect(vol).ToNot(BeNil())
			Expect(vol.Name).To(Equal("vol1"))
			Expect(vol.ConfigMap.Name).To(Equal("cfg1"))
		})

		It("should return an emptyDir volume", func() {
			vol := getEmptyDirVolume("vol1")

			Expect(vol).ToNot(BeNil())
			Expect(vol.Name).To(Equal("vol1"))
			Expect(vol.EmptyDir).ToNot(BeNil())
		})

		It("should return a valid volume mount spec", func() {
			mount := getVolumeMountSpec("vol1", "/tmp", true)

			Expect(mount).ToNot(BeNil())
			Expect(mount.Name).To(Equal("vol1"))
			Expect(mount.MountPath).To(Equal("/tmp"))
			Expect(mount.ReadOnly).To(BeTrue())
		})

		It("should return a valid volume mount spec with subpath", func() {
			mount := getVolumeMountSpecWithSubPath("vol1", "/tmp/file.tf", "file.tf", true)

			Expect(mount).ToNot(BeNil())
			Expect(mount.Name).To(Equal("vol1"))
			Expect(mount.MountPath).To(Equal("/tmp/file.tf"))
			Expect(mount.SubPath).To(Equal("file.tf"))
			Expect(mount.ReadOnly).To(BeTrue())
		})

		It("should return an env var", func() {
			env := getEnvVariable("key", "value")

			Expect(env).ToNot(BeNil())
			Expect(env.Name).To(Equal("key"))
			Expect(env.Value).To(Equal("value"))
		})

		It("should return an env var", func() {
			env := getEnvVariableFromFieldSelector("key", "metadata.name")

			Expect(env).ToNot(BeNil())
			Expect(env.Name).To(Equal("key"))
			Expect(env.ValueFrom.FieldRef).ToNot(BeNil())
			Expect(env.ValueFrom.FieldRef.FieldPath).To(Equal("metadata.name"))
		})
	})
})
