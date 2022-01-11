package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// returns a volume spec
func getVolumeSpec(name string, source corev1.VolumeSource) corev1.Volume {
	return corev1.Volume{
		Name:         name,
		VolumeSource: source,
	}
}

// returns and emptyDir volume spec
func getEmptyDirVolume(name string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

// returns a volume mount spec
func getVolumeMountSpec(volumeName string, mountPath string, readOnly bool) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
		ReadOnly:  readOnly,
	}
}

func getEnvVariable(name string, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  name,
		Value: value,
	}
}

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
