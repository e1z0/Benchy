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

import "github.com/mappu/miqt/qt"

type Tile struct {
	Box   *qt.QGroupBox
	Value *qt.QLabel
}

func NewTile(title string) *Tile {
	g := qt.NewQGroupBox4(title, nil)
	v := qt.NewQVBoxLayout(g.QWidget)

	lab := qt.NewQLabel5("—", nil)
	f := lab.Font()
	f.SetBold(true)
	f.SetPointSize(f.PointSize() + 6)
	lab.SetFont(f)

	v.AddWidget(lab.QWidget)
	v.AddStretch()

	return &Tile{Box: g, Value: lab}
}
