package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// getVolumeSpec returns a volume spec
func getVolumeSpec(name string, source corev1.VolumeSource) corev1.Volume {
	return corev1.Volume{
		Name:         name,
		VolumeSource: source,
	}
}

// getVolumeSpecFromConfigMap returns a volume spec from configMap
func getVolumeSpecFromConfigMap(volumeName string, configMapName string) corev1.Volume {
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapName,
				},
			},
		},
	}
}

// getEmptyDirVolume returns and emptyDir volume spec
func getEmptyDirVolume(name string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

// getVolumeMountSpec returns a volume mount spec
func getVolumeMountSpec(volumeName string, mountPath string, readOnly bool) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
		ReadOnly:  readOnly,
	}
}

// getVolumeMountSpecWithSubPath returns a volume mount spec with subpath option
func getVolumeMountSpecWithSubPath(volumeName string, mountPath string, subPath string, readOnly bool) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
		ReadOnly:  readOnly,
		SubPath:   subPath,
	}
}

// getEnvVariable returns a Kubernetes environment variable spec
func getEnvVariable(name string, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  name,
		Value: value,
	}
}

// getEnvVariableFromFieldSelector returns a Kubernetes environment variable from a field selector
func getEnvVariableFromFieldSelector(name string, path string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: path,
			},
		},
	}
}
