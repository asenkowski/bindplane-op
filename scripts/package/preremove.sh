#!/bin/bash
# Copyright  observIQ, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script must work on the following:
#  debian 10
#  debian 11
#  ubuntu 18.04
#  ubuntu 20.04
#  ubuntu 22.04
#  centos / rhel 7
#  centos / rhel 8
#  SLES 12
#  SLES 15

set -e

# Stop and disable service
if systemctl disable --now bindplane ; then
    echo "Service stopped: bindplane"
    echo "Service disabled: bindplane"
else
    exit 0
fi
