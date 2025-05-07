/*
Copyright 2024-2025 the Unikorn Authors.

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

package get

import (
	"github.com/spf13/cobra"

	"github.com/unikorn-cloud/kubectl-unikorn/pkg/factory"
	"github.com/unikorn-cloud/kubectl-unikorn/pkg/get/clustermanager"
	"github.com/unikorn-cloud/kubectl-unikorn/pkg/get/kubernetescluster"
	"github.com/unikorn-cloud/kubectl-unikorn/pkg/get/user"
)

func Command(factory *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource information",
	}

	cmd.AddCommand(
		user.Command(factory),
		kubernetescluster.Command(factory),
		clustermanager.Command(factory),
	)

	return cmd
}
