/*
Copyright The containerd Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/containerd/containerd/plugin"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	for desc, test := range map[string]struct {
		config      *PluginConfig
		expectedErr string
		expected    *PluginConfig
	}{
		"deprecated untrusted_workload_runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					UntrustedWorkloadRuntime: Runtime{
						Type: "untrusted",
					},
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					UntrustedWorkloadRuntime: Runtime{
						Type: "untrusted",
					},
					Runtimes: map[string]Runtime{
						RuntimeUntrusted: {
							Type: "untrusted",
						},
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
			},
		},
		"both untrusted_workload_runtime and runtime[untrusted]": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					UntrustedWorkloadRuntime: Runtime{
						Type: "untrusted-1",
					},
					Runtimes: map[string]Runtime{
						RuntimeUntrusted: {
							Type: "untrusted-2",
						},
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("conflicting definitions: configuration includes both `untrusted_workload_runtime` and `runtimes[%q]`", RuntimeUntrusted),
		},
		"deprecated default_runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntime: Runtime{
						Type: "default",
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntime: Runtime{
						Type: "default",
					},
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
			},
		},
		"no default_runtime_name": {
			config:      &PluginConfig{},
			expectedErr: "`default_runtime_name` is empty",
		},
		"no runtime[default_runtime_name]": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
				},
			},
			expectedErr: "no corresponding runtime configured in `runtimes` for `default_runtime_name`",
		},
		"deprecated systemd_cgroup for v1 runtime": {
			config: &PluginConfig{
				SystemdCgroup: true,
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				SystemdCgroup: true,
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"deprecated systemd_cgroup for v2 runtime": {
			config: &PluginConfig{
				SystemdCgroup: true,
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`systemd_cgroup` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"no_pivot for v1 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					NoPivot:            true,
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					NoPivot:            true,
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"no_pivot for v2 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					NoPivot:            true,
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`no_pivot` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"deprecated runtime_engine for v1 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Engine: "runc",
							Type:   plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Engine: "runc",
							Type:   plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"deprecated runtime_engine for v2 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Engine: "runc",
							Type:   plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`runtime_engine` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"deprecated runtime_root for v1 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Root: "/run/containerd/runc",
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Root: "/run/containerd/runc",
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
			},
		},
		"deprecated runtime_root for v2 runtime": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Root: "/run/containerd/runc",
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
			},
			expectedErr: fmt.Sprintf("`runtime_root` only works for runtime %s", plugin.RuntimeLinuxV1),
		},
		"deprecated auths": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
				Registry: Registry{
					Auths: map[string]AuthConfig{
						"https://gcr.io": {Username: "test"},
					},
				},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeRuncV1,
						},
					},
				},
				Registry: Registry{
					Configs: map[string]RegistryConfig{
						"https://gcr.io": {
							Auth: &AuthConfig{
								Username: "test",
							},
						},
					},
					Auths: map[string]AuthConfig{
						"https://gcr.io": {Username: "test"},
					},
				},
			},
		},
		"invalid stream_idle_timeout": {
			config: &PluginConfig{
				StreamIdleTimeout: "invalid",
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: "default",
						},
					},
				},
			},
			expectedErr: "invalid stream idle timeout",
		},
		"valid id mapping": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
				NodeWideUIDMapping: LinuxIDMapping{0, 800000, 65536},
				NodeWideGIDMapping: LinuxIDMapping{0, 800000, 65536},
			},
			expected: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
				NodeWideUIDMapping: LinuxIDMapping{0, 800000, 65536},
				NodeWideGIDMapping: LinuxIDMapping{0, 800000, 65536},
			},
		},
		"invalid id mapping": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
				NodeWideUIDMapping: LinuxIDMapping{1, 100000, 65536},
				NodeWideGIDMapping: LinuxIDMapping{1, 100000, 65536},
			},
			expectedErr: "missing root id in container",
		},
		"different uid and gid mappings": {
			config: &PluginConfig{
				ContainerdConfig: ContainerdConfig{
					DefaultRuntimeName: RuntimeDefault,
					Runtimes: map[string]Runtime{
						RuntimeDefault: {
							Type: plugin.RuntimeLinuxV1,
						},
					},
				},
				NodeWideUIDMapping: LinuxIDMapping{0, 100000, 65536},
				NodeWideGIDMapping: LinuxIDMapping{0, 200000, 65536},
			},
			expectedErr: "different mappings for uid and gid not yet supported",
		},
	} {
		t.Run(desc, func(t *testing.T) {
			err := ValidatePluginConfig(context.Background(), test.config)
			if test.expectedErr != "" {
				assert.Contains(t, err.Error(), test.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, test.config)
			}
		})
	}
}
