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
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/e1z0/Benchy/internal/benchmarks"
	"github.com/e1z0/Benchy/internal/scoring"
	"github.com/e1z0/Benchy/internal/sysinfo"
	"github.com/e1z0/Benchy/internal/ui"

	"github.com/mappu/miqt/qt"
)

type tabWidgets struct {
	table    *qt.QTableWidget
	overall  *qt.QLabel
	chart    *ui.BarChart
	tileCPU  *ui.Tile
	tileMem  *ui.Tile
	tileStor *ui.Tile
	tileImg  *ui.Tile
	lastJSON []byte
}

func newTab(title string, parent *qt.QTabWidget) *tabWidgets {
	w := qt.NewQWidget(nil)
	v := qt.NewQVBoxLayout(w)

	tiles := qt.NewQHBoxLayout2()
	tileCPU := ui.NewTile("CPU")
	tileMem := ui.NewTile("Memory")
	tileStor := ui.NewTile("Storage")
	tileImg := ui.NewTile("Image")
	tiles.AddWidget(tileCPU.Box.QWidget)
	tiles.AddWidget(tileMem.Box.QWidget)
	tiles.AddWidget(tileStor.Box.QWidget)
	tiles.AddWidget(tileImg.Box.QWidget)

	overall := qt.NewQLabel5("Overall: —", nil)
	f := overall.Font()
	f.SetPointSize(f.PointSize() + 8)
	f.SetBold(true)
	overall.SetFont(f)

	tbl := qt.NewQTableWidget4(0, 6, nil)
	tbl.SetHorizontalHeaderLabels([]string{"Test", "Threads", "Duration (s)", "Throughput", "Score", "Notes"})
	tbl.HorizontalHeader().SetStretchLastSection(true)

	chart := ui.NewBarChart(nil)

	v.AddLayout(tiles.QLayout)
	v.AddWidget(overall.QWidget)
	v.AddWidget(tbl.QWidget)
	v.AddWidget(chart.QWidget)

	parent.AddTab(w, title)
	return &tabWidgets{
		table: tbl, overall: overall, chart: chart,
		tileCPU: tileCPU, tileMem: tileMem, tileStor: tileStor, tileImg: tileImg,
	}
}

func main() {
	app := qt.NewQApplication(os.Args)
	ui.EnableDark(app)

	win := qt.NewQMainWindow(nil)
	win.SetWindowTitle("Benchy")

	central := qt.NewQWidget(nil)
	root := qt.NewQVBoxLayout(central)

	si := sysinfo.Collect()
	info := qt.NewQPlainTextEdit(nil)
	info.SetReadOnly(true)
	info.SetPlainText(si.String())

	durLbl := qt.NewQLabel3("Duration (s):")
	dur := qt.NewQSpinBox(nil)
	dur.SetRange(1, 60)
	dur.SetValue(5)
	run := qt.NewQPushButton3("Run Both")
	exp1 := qt.NewQPushButton3("Export Single-Core JSON")
	exp1.SetEnabled(false)
	expm := qt.NewQPushButton3("Export Multi-Core JSON")
	expm.SetEnabled(false)

	tabs := qt.NewQTabWidget(nil)
	single := newTab("Single-Core", tabs)
	multi := newTab("Multi-Core", tabs)

	opts := qt.NewQHBoxLayout(nil)
	opts.AddWidget(durLbl.QWidget)
	opts.AddWidget(dur.QWidget)
	opts.AddSpacing(8)
	opts.AddWidget(run.QWidget)
	opts.AddStretch()
	opts.AddWidget(exp1.QWidget)
	opts.AddWidget(expm.QWidget)

	root.AddWidget(info.QWidget)
	root.AddLayout(opts.QLayout)
	root.AddWidget(tabs.QWidget)
	win.SetCentralWidget(central)

	// Ordered test list
	tests := []ui.TestSpec{
		{Name: "CPU SHA-256", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunCPUSHA256(ctx, d, th)
		}},
		{Name: "AES-CTR", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunCPUAES(ctx, d, th, 32)
		}},
		{Name: "Zstd Compress", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunCPUZstd(ctx, d, th, 3)
		}},
		{Name: "Gzip Compress", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunCPUGzip(ctx, d, th, -1)
		}},
		{Name: "JSON Parse", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunCPUJSONParse(ctx, d, th)
		}},
		{Name: "MatMul", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunMatMul(ctx, d, th, 256)
		}},
		{Name: "Memory copy", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunMemCopy(ctx, d, th)
		}},
		{Name: "Gaussian Blur 1080p", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunImageBlur(ctx, d, th)
		}},
		{Name: "Disk seq R/W", Run: func(ctx context.Context, d time.Duration, th int) benchmarks.Result {
			return benchmarks.RunDiskSeq(ctx, d, filepath.Join(os.TempDir(), "benchyqt.seq"))
		}},
	}

	run.OnClicked(func() {
		run.SetEnabled(false)
		exp1.SetEnabled(false)
		expm.SetEnabled(false)

		d := time.Duration(dur.Value()) * time.Second

		// Single-Core first
		sres := ui.RunSuiteDialog(win.QWidget, "Single-Core", 1, d, tests)
		populateTab(single, sres.Results)

		// Multi-Core second
		mres := ui.RunSuiteDialog(win.QWidget, "Multi-Core", sysinfo.Collect().LogicalCPUs, d, tests)
		populateTab(multi, mres.Results)

		run.SetEnabled(true)
		exp1.SetEnabled(true)
		expm.SetEnabled(true)
	})

	exp1.OnClicked(func() {
		if len(single.lastJSON) == 0 {
			return
		}
		fn := filepath.Join(userHome(), fmt.Sprintf("benchyqt-single-%d.json", time.Now().Unix()))
		_ = os.WriteFile(fn, single.lastJSON, 0644)
		info.AppendPlainText("Saved: " + fn)
	})
	expm.OnClicked(func() {
		if len(multi.lastJSON) == 0 {
			return
		}
		fn := filepath.Join(userHome(), fmt.Sprintf("benchyqt-multi-%d.json", time.Now().Unix()))
		_ = os.WriteFile(fn, multi.lastJSON, 0644)
		info.AppendPlainText("Saved: " + fn)
	})

	win.Resize(1200, 760)
	win.Show()
	qt.QApplication_Exec()
}

func shortName(s string) string {
	switch s {
	case "CPU SHA-256":
		return "SHA256"
	case "Gaussian Blur 1080p":
		return "Blur1080p"
	case "Memory copy":
		return "MemCopy"
	case "Disk seq R/W":
		return "DiskSeq"
	case "Zstd Compress":
		return "Zstd"
	case "Gzip Compress":
		return "Gzip"
	case "JSON Parse":
		return "JSON"
	case "AES-CTR":
		return "AES"
	case "MatMul":
		return "MatMul"
	}
	return s
}

func userHome() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return h
}

func geo(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	prod := 1.0
	n := 0
	for _, v := range vals {
		if v > 0 {
			prod *= v
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return math.Pow(prod, 1.0/float64(n))
}

func populateTab(t *tabWidgets, results []benchmarks.Result) {
	t.table.SetRowCount(0)

	var bars []ui.Bar
	var cpu, mem, stor, img []float64

	for _, r := range results {
		row := t.table.RowCount()
		t.table.InsertRow(row)

		score := scoring.Score(r)

		t.table.SetItem(row, 0, qt.NewQTableWidgetItem2(r.Name))
		t.table.SetItem(row, 1, qt.NewQTableWidgetItem2(fmt.Sprintf("%d", r.Threads)))
		t.table.SetItem(row, 2, qt.NewQTableWidgetItem2(fmt.Sprintf("%.2f", r.Duration.Seconds())))
		t.table.SetItem(row, 3, qt.NewQTableWidgetItem2(r.ThroughputString()))
		t.table.SetItem(row, 4, qt.NewQTableWidgetItem2(fmt.Sprintf("%.0f", score)))
		t.table.SetItem(row, 5, qt.NewQTableWidgetItem2(r.Notes))

		bars = append(bars, ui.Bar{Label: shortName(r.Name), Value: score})

		switch r.Name {
		case "CPU SHA-256", "AES-CTR", "Zstd Compress", "Gzip Compress", "JSON Parse", "MatMul":
			cpu = append(cpu, score)
		case "Memory copy":
			mem = append(mem, score)
		case "Disk seq R/W":
			stor = append(stor, score)
		case "Gaussian Blur 1080p":
			img = append(img, score)
		}
	}

	// overall + section tiles
	overall := scoring.Aggregate(results)
	t.overall.SetText(fmt.Sprintf("Overall: %.0f", overall))

	t.tileCPU.Value.SetText(fmt.Sprintf("%.0f", geo(cpu)))
	t.tileMem.Value.SetText(fmt.Sprintf("%.0f", geo(mem)))
	t.tileStor.Value.SetText(fmt.Sprintf("%.0f", geo(stor)))
	t.tileImg.Value.SetText(fmt.Sprintf("%.0f", geo(img)))

	// chart
	t.chart.SetData(bars)

	// export blob
	b, _ := json.MarshalIndent(struct {
		System   sysinfo.Info        `json:"system"`
		Results  []benchmarks.Result `json:"results"`
		Overall  float64             `json:"overall"`
		Sections map[string]float64  `json:"sections"`
	}{
		System:  sysinfo.Collect(),
		Results: results,
		Overall: overall,
		Sections: map[string]float64{
			"CPU": geo(cpu), "Memory": geo(mem), "Storage": geo(stor), "Image": geo(img),
		},
	}, "", "  ")
	t.lastJSON = b
}
