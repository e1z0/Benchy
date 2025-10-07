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
	"context"
	"crypto/sha256"
	"runtime"
	"sync"
	"time"
)

type shaState struct{ buf []byte }

func RunCPUSHA256(ctx context.Context, dur time.Duration, threads int) Result {
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	var wg sync.WaitGroup
	var ops uint64
	states := make([]shaState, threads)
	for i := range states {
		states[i] = shaState{buf: make([]byte, 8*1024)}
	}

	done := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		st := &states[i]
		go func(seed byte) {
			defer wg.Done()
			var local uint64
			b := st.buf
			for {
				select {
				case <-ctx.Done():
					ops += local
					return
				default:
					for j := range b {
						b[j] ^= seed + byte(j)
					}
					_ = sha256.Sum256(b)
					local++
				}
			}
		}(byte(i))
	}
	go func() { wg.Wait(); close(done) }()
	<-done
	return Result{Name: "CPU SHA-256", Threads: threads, Duration: dur, Ops: ops, Unit: "hash/s"}
}
