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
	"fmt"
	"math"
	"sort"

	"github.com/mappu/miqt/qt"
)

type Bar struct {
	Label string
	Value float64
}

type BarChart struct {
	*qt.QWidget
	data       []Bar
	max        float64
	hoverIndex int
}

func NewBarChart(parent *qt.QWidget) *BarChart {
	w := qt.NewQWidget(parent)
	bc := &BarChart{QWidget: w, hoverIndex: -1}
	w.SetMouseTracking(true)
	w.SetMinimumHeight(200)

	w.OnMouseMoveEvent(func(super func(*qt.QMouseEvent), e *qt.QMouseEvent) {
		idx := bc.hitTest(e.Pos().X(), e.Pos().Y())
		if idx != bc.hoverIndex {
			bc.hoverIndex = idx
			bc.Update()
		}
	})

	w.OnLeaveEvent(func(super func(*qt.QEvent), e *qt.QEvent) {
		bc.hoverIndex = -1
		bc.Update()
	})

	w.OnPaintEvent(func(super func(*qt.QPaintEvent), e *qt.QPaintEvent) {
		p := qt.NewQPainter()
		if !p.Begin(w.QPaintDevice) {
			return
		}
		defer p.End()

		r := w.Rect()
		margin := 12
		base := qt.NewQRect4(r.X()+margin, r.Y()+margin, r.Width()-2*margin, r.Height()-2*margin)

		if len(bc.data) == 0 || bc.max <= 0 {
			p.DrawText6(r, int(qt.AlignCenter), "No data")
			return
		}

		// chart area with bottom space for labels
		labelH := 22
		leftPad := 36 // y-axis labels
		topPad := 10
		chart := qt.NewQRect4(base.X()+leftPad, base.Y()+topPad, base.Width()-leftPad, base.Height()-labelH-topPad-4)

		// axes
		p.DrawRectWithRect(chart)
		y0 := chart.Y() + chart.Height()

		// y ticks (5)
		ticks := 5
		for i := 0; i <= ticks; i++ {
			t := float64(i) / float64(ticks)
			y := y0 - int(t*float64(chart.Height()))
			p.DrawLine2(chart.X(), y, chart.X()+chart.Width(), y)
			label := fmt.Sprintf("%.0f", t*bc.max)
			p.DrawText7(chart.X()-34, y-8, 32, 16, int(qt.AlignRight|qt.AlignVCenter), label)
		}

		bars := len(bc.data)
		gap := 10
		barW := int(math.Max(1, float64(chart.Width()-gap*(bars+1))/float64(bars)))
		x := chart.X() + gap

		// draw bars
		for i, b := range bc.data {
			ratio := b.Value / bc.max
			if ratio > 1 {
				ratio = 1
			}
			h := int(float64(chart.Height()) * ratio)
			y := y0 - h
			rect := qt.NewQRect4(x, y, barW, h)
			// fill
			if i == bc.hoverIndex {
				p.FillRect3(rect, qt.NewQBrush11(qt.NewQColor3(180, 210, 255), 1))
			} else {
				p.FillRect3(rect, qt.NewQBrush11(qt.NewQColor3(120, 170, 220), 1))
			}
			p.DrawRectWithRect(rect)

			// x label
			p.DrawText7(x, base.Y()+base.Height()-labelH+2, barW, labelH, int(qt.AlignHCenter|qt.AlignTop), b.Label)

			// value label on top
			p.DrawText7(x, y-18, barW, 16, int(qt.AlignHCenter|qt.AlignBottom), fmt.Sprintf("%.0f", b.Value))
			x += barW + gap
		}

		// hover tooltip
		if bc.hoverIndex >= 0 && bc.hoverIndex < len(bc.data) {
			// draw a simple tooltip box near top
			text := fmt.Sprintf("%s: %.0f", bc.data[bc.hoverIndex].Label, bc.data[bc.hoverIndex].Value)
			wtxt := p.FontMetrics().HorizontalAdvance(text) + 12
			rect := qt.NewQRect4(chart.X()+8, chart.Y()+8, wtxt, 22)
			p.FillRect3(rect, qt.NewQBrush11(qt.NewQColor3(50, 50, 50), 1))
			p.SetPen(qt.NewQColor3(230, 230, 230))
			p.DrawText6(rect, int(qt.AlignCenter), text)
			p.SetPen(qt.NewQColor3(0, 0, 0)) // reset pen
		}
	})
	return bc
}

func (bc *BarChart) SetData(items []Bar) {
	sort.Slice(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	bc.data = items
	bc.max = 0
	for _, it := range items {
		if it.Value > bc.max {
			bc.max = it.Value
		}
	}
	if bc.max < 1 {
		bc.max = 1
	}
	bc.Update()
}

func (bc *BarChart) hitTest(x, y int) int {
	r := bc.Rect()
	margin := 12
	base := qt.NewQRect4(r.X()+margin, r.Y()+margin, r.Width()-2*margin, r.Height()-2*margin)
	labelH := 22
	leftPad := 36
	topPad := 10
	chart := qt.NewQRect4(base.X()+leftPad, base.Y()+topPad, base.Width()-leftPad, base.Height()-labelH-topPad-4)

	if len(bc.data) == 0 {
		return -1
	}
	bars := len(bc.data)
	gap := 10
	barW := int(math.Max(1, float64(chart.Width()-gap*(bars+1))/float64(bars)))
	bx := chart.X() + gap
	for i := 0; i < bars; i++ {
		br := qt.NewQRect4(bx, chart.Y(), barW, chart.Height())
		if br.Contains3(x, y, false) {
			return i
		}
		bx += barW + gap
	}
	return -1
}
