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
package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/e1z0/Benchy/internal/benchmarks"

	"github.com/mappu/miqt/qt"
	"github.com/mappu/miqt/qt/mainthread"
)

type testFn func(ctx context.Context, dur time.Duration, threads int) benchmarks.Result

type TestSpec struct {
	Name string
	Run  testFn
}

type RunResult struct {
	Results  []benchmarks.Result
	Canceled bool
}

func RunSuiteDialog(parent *qt.QWidget, mode string, threads int, dur time.Duration, tests []TestSpec) RunResult {
	dlg := qt.NewQDialog(parent)
	dlg.SetWindowTitle("Running Benchmarks — " + mode)

	v := qt.NewQVBoxLayout(dlg.QWidget)
	title := qt.NewQLabel3(fmt.Sprintf("%s — %d tests", mode, len(tests)))
	cur := qt.NewQLabel3("Ready…")
	overall := qt.NewQProgressBar(nil)
	overall.SetRange(0, len(tests)*100)
	per := qt.NewQProgressBar(nil)
	per.SetRange(0, 100)
	per.SetFormat("Current: %p%")
	log := qt.NewQPlainTextEdit(nil)
	log.SetReadOnly(true)
	log.SetMinimumHeight(140)

	btns := qt.NewQHBoxLayout2()
	btnCancel := qt.NewQPushButton3("Cancel")
	btns.AddStretch()
	btns.AddWidget(btnCancel.QWidget)

	v.AddWidget(title.QWidget)
	v.AddWidget(cur.QWidget)
	v.AddWidget(overall.QWidget)
	v.AddWidget(per.QWidget)
	v.AddWidget(log.QWidget)
	v.AddLayout(btns.QLayout)
	dlg.Resize(560, 360)
	dlg.Show()

	res := RunResult{}
	finished := false
	
	ctx, cancel := context.WithCancel(context.Background())

	btnCancel.OnClicked(func() {
		if finished {
			dlg.Accept() // <- actually close the modal dialog
			return
		}
		cancel() // <- signal cancellation
	})

	go func() {
		defer cancel()

		for i, t := range tests {
			if ctx.Err() != nil {
				break
			}

			mainthread.Wait(func() {
				cur.SetText(fmt.Sprintf("Test %d/%d — %s", i+1, len(tests), t.Name))
				log.AppendPlainText(fmt.Sprintf("> %s", t.Name))
				per.SetValue(0)
			})

			// Run the test in a separate goroutine so we can update progress
			testDone := make(chan benchmarks.Result, 1)
			go func() {
				testDone <- t.Run(ctx, dur, threads)
			}()

			start := time.Now()
			tick := time.NewTicker(100 * time.Millisecond)

		loop:
			for {
				select {
				case r := <-testDone:
					res.Results = append(res.Results, r)
					mainthread.Wait(func() {
						per.SetValue(100)
						overall.SetValue((i + 1) * 100)
					})
					break loop
				case <-tick.C:
					p := int(float64(time.Since(start)) / float64(dur) * 100)
					if p > 99 { p = 99 }
					if p < 0 { p = 0 }
					
					mainthread.Wait(func() {
						per.SetValue(p)
						overall.SetValue(i*100 + p)
					})
				case <-ctx.Done():
					break loop
				}
			}
			tick.Stop()
		}

		mainthread.Wait(func() {
			if ctx.Err() != nil {
				res.Canceled = true
				log.AppendPlainText("Canceled by user.")
				cur.SetText("Canceled.")
			} else {
				log.AppendPlainText("Completed.")
				cur.SetText("Finished.")
			}
			btnCancel.SetText("Close")
			finished = true
		})
	}()

	dlg.Exec()
	return res
}
