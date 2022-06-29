// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package get

import (
	"github.com/spf13/cobra"

	"github.com/observiq/bindplane/internal/cli"
)

// Command returns the BindPlane get cobra command.
func Command(bindplane *cli.BindPlane) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more resources",
	}

	cmd.AddCommand(
		AgentsCommand(bindplane),
		ConfigurationsCommand(bindplane),
		DestinationsCommand(bindplane),
		DestinationTypesCommand(bindplane),
		SourcesCommand(bindplane),
		SourceTypesCommand(bindplane),
	)

	return cmd
}
