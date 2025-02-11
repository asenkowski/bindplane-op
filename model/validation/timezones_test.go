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

package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsTimezone(t *testing.T) {
	testCases := []struct {
		value       string
		expectValid bool
	}{
		{
			"UTC",
			true,
		},
		{
			"America/NewJersey",
			false,
		},
	}

	for _, test := range testCases {
		require.Equal(t, test.expectValid, IsTimezone(test.value))
	}
}
