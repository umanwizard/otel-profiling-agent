//go:build amd64 && !dummy

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package support

import (
	_ "embed"
)

//go:embed ebpf/tracer.ebpf.amd64
var tracerData []byte
