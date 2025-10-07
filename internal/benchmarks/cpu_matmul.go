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
	"fmt"
	"runtime"
	"sync"
	"time"
)

// RunMatMul performs naive C = A*B on n x n matrices of float64.
// Reports operations as floating point ops per second (approx 2*n^3) and returns GFLOP/s.
func RunMatMul(ctx context.Context, dur time.Duration, threads, n int) Result {
	if n <= 0 {
		n = 256
	}
	if threads <= 0 {
		threads = runtime.NumCPU()
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	type job struct{ repeat int }
	jobs := make(chan job, threads)
	var wg sync.WaitGroup
	var flops uint64

	worker := func() {
		defer wg.Done()
		A := make([]float64, n*n)
		B := make([]float64, n*n)
		C := make([]float64, n*n)
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-jobs:
				if !ok {
					return
				}
				// i-k-j loop order
				for i := 0; i < n; i++ {
					for k := 0; k < n; k++ {
						aik := A[i*n+k]
						row := i * n
						col := k * n
						for j := 0; j < n; j++ {
							C[row+j] += aik * B[col+j]
						}
					}
				}
				// count 2*n^3 flops per multiply-add
				flops += uint64(2) * uint64(n) * uint64(n) * uint64(n)
			}
		}
	}
	for w := 0; w < threads; w++ {
		wg.Add(1)
		go worker()
	}

	start := time.Now()
	for time.Since(start) < dur {
		select {
		case <-ctx.Done():
			break
		default:
			jobs <- job{1}
		}
	}
	close(jobs)
	wg.Wait()

	gflops := float64(flops) / 1e9 / dur.Seconds()
	return Result{Name: "MatMul", Threads: threads, Duration: dur, Ops: uint64(gflops * 1e6), Unit: "GFLOP/s", Notes: fmt.Sprintf("n=%d", n)}
}
