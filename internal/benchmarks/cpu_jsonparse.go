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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type sampleRec struct {
	ID     int               `json:"id"`
	Name   string            `json:"name"`
	Values []float64         `json:"values"`
	Tags   map[string]string `json:"tags"`
}

func genJSON() []byte {
	recs := make([]sampleRec, 500)
	r := rand.New(rand.NewSource(99))
	for i := range recs {
		vals := make([]float64, 64)
		for j := range vals {
			vals[j] = r.NormFloat64()
		}
		tags := map[string]string{"k": "v", "env": "prod", "zone": "eu"}
		recs[i] = sampleRec{ID: i, Name: fmt.Sprintf("rec-%d", i), Values: vals, Tags: tags}
	}
	b, _ := json.Marshal(recs)
	return b
}

func RunCPUJSONParse(ctx context.Context, dur time.Duration, threads int) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	payload := genJSON()
	var wg sync.WaitGroup
	var bytesOK uint64

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var local uint64
			for {
				select {
				case <-ctx.Done():
					bytesOK += local
					return
				default:
					dec := json.NewDecoder(bytes.NewReader(payload))
					var out []sampleRec
					_ = dec.Decode(&out)
					local += uint64(len(payload))
				}
			}
		}()
	}
	wg.Wait()
	return Result{Name: "JSON Parse", Threads: threads, Duration: dur, Bytes: bytesOK, Unit: "B/s"}
}
