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

package copy

import (
	"bytes"
	"context"
	"testing"

	"github.com/observiq/bindplane-op/client"
	"github.com/observiq/bindplane-op/common"
	"github.com/observiq/bindplane-op/internal/cli"
	"github.com/stretchr/testify/require"
)

func setupBindPlane(buffer *bytes.Buffer) *cli.BindPlane {
	bindplane := cli.NewBindPlane(common.InitConfig(""), buffer)
	bindplane.SetClient(&mockClient{})
	return bindplane
}

type mockClient struct {
	client.BindPlane
}

var gotArgs []any

func (mc *mockClient) CopyConfig(_ context.Context, configName, copyName string) error {
	gotArgs = []any{configName, copyName}
	return nil
}

func TestCopyConfigCommmand(t *testing.T) {
	out := bytes.NewBufferString("")
	bp := setupBindPlane(out)
	t.Run("errors when two arguments are not present", func(t *testing.T) {
		cmd := ConfigurationCommand(bp)
		cmd.SetArgs([]string{})
		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("calls CopyConfig with correct args", func(t *testing.T) {
		cmd := ConfigurationCommand(bp)
		cmd.SetArgs([]string{"blah", "foo"})

		cmd.Execute()

		require.Equal(t, []any{"blah", "foo"}, gotArgs)
	})
}
