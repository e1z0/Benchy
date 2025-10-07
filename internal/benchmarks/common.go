/* SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Benchy
 * Copyright (C) 2025 e1z0 <e1z0@icloud.com>
 *
 * This file is part of Benchy.
 *
 * Benchy is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Benchy is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Benchy. If not, see <https://www.gnu.org/licenses/>.
 */
package benchmarks

import (
	"fmt"
	"time"
)

type Result struct {
	Name     string        `json:"name"`
	Threads  int           `json:"threads"`
	Duration time.Duration `json:"duration"`
	Ops      uint64        `json:"ops"`
	Bytes    uint64        `json:"bytes"`
	Unit     string        `json:"unit"`
	Err      string        `json:"err,omitempty"`
	Notes    string        `json:"notes,omitempty"`
}

func (r Result) ThroughputString() string {
	if r.Unit == "B/s" {
		bps := uint64(0)
		if r.Duration > 0 {
			bps = uint64(float64(r.Bytes) / r.Duration.Seconds())
		}
		return humanBytes(bps) + "/s"
	}
	if r.Unit == "GFLOP/s" {
		return fmt.Sprintf("%.2f %s", float64(r.Ops)/1e6, r.Unit)
	}
	if r.Ops > 0 {
		return fmt.Sprintf("%d %s", r.Ops/uint64(max64(1, int64(r.Duration.Seconds()))), r.Unit)
	}
	return "—"
}

func humanBytes(b uint64) string {
	suffix := []string{"B", "KB", "MB", "GB", "TB"}
	f := float64(b)
	i := 0
	for f >= 1024 && i < len(suffix)-1 {
		f /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", f, suffix[i])
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
