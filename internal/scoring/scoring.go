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
package scoring

import (
	"math"

	"github.com/e1z0/Benchy/internal/benchmarks"
)

const Baseline = 2500.0

var Reference = map[string]float64{
	"CPU SHA-256":         200000.0,    // hash/s
	"AES-CTR":             1500 << 20,  // B/s
	"Zstd Compress":       400 << 20,   // B/s
	"Gzip Compress":       250 << 20,   // B/s
	"JSON Parse":          300 << 20,   // B/s
	"MatMul":              50.0,        // GFLOP/s
	"Memory copy":         20000 << 20, // B/s
	"Gaussian Blur 1080p": 30e6,        // px/s
	"Disk seq R/W":        800 << 20,   // B/s
}

func Score(r benchmarks.Result) float64 {
	ref, ok := Reference[r.Name]
	if !ok || ref <= 0 {
		return 0
	}
	var tp float64
	switch r.Unit {
	case "B/s":
		if r.Duration.Seconds() > 0 {
			tp = float64(r.Bytes) / r.Duration.Seconds()
		}
	case "GFLOP/s":
		tp = float64(r.Ops) / 1e6
	default:
		tp = float64(r.Ops)
	}
	if tp <= 0 {
		return 0
	}
	return (tp / ref) * Baseline
}

func Aggregate(rs []benchmarks.Result) float64 {
	prod := 1.0
	n := 0
	for _, r := range rs {
		s := Score(r)
		if s > 0 {
			prod *= s
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return math.Pow(prod, 1.0/float64(n))
}
