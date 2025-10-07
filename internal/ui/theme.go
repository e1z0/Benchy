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

func EnableDark(app *qt.QApplication) {
	pal := qt.NewQPalette()
	// basic dark palette
	pal.SetColor2(qt.QPalette__Window, qt.NewQColor3(37, 37, 38))
	pal.SetColor2(qt.QPalette__WindowText, qt.NewQColor3(230, 230, 230))
	pal.SetColor2(qt.QPalette__Base, qt.NewQColor3(30, 30, 30))
	pal.SetColor2(qt.QPalette__AlternateBase, qt.NewQColor3(45, 45, 45))
	pal.SetColor2(qt.QPalette__ToolTipBase, qt.NewQColor3(64, 64, 64))
	pal.SetColor2(qt.QPalette__ToolTipText, qt.NewQColor3(230, 230, 230))
	pal.SetColor2(qt.QPalette__Text, qt.NewQColor3(230, 230, 230))
	pal.SetColor2(qt.QPalette__Button, qt.NewQColor3(45, 45, 45))
	pal.SetColor2(qt.QPalette__ButtonText, qt.NewQColor3(230, 230, 230))
	pal.SetColor2(qt.QPalette__BrightText, qt.NewQColor3(255, 0, 0))
	pal.SetColor2(qt.QPalette__Highlight, qt.NewQColor3(38, 79, 120))
	pal.SetColor2(qt.QPalette__HighlightedText, qt.NewQColor3(255, 255, 255))
	qt.QApplication_SetPalette(pal)
}
