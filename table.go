package table

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectionMode defines how selections work (can be combined with bitwise OR)
type SelectionMode int

const (
	SelectionOff    SelectionMode = 0
	SelectionRow    SelectionMode = 1 << 0 // 1
	SelectionColumn SelectionMode = 1 << 1 // 2
	SelectionCell   SelectionMode = 1 << 2 // 4
)

type Theme struct {
	Header       lipgloss.Style
	Cell         lipgloss.Style
	Border       lipgloss.Style
	SelectedRow  lipgloss.Style
	SelectedCell lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("63")),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Border: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		SelectedRow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")),
		SelectedCell: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("129")),
	}
}

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Home     key.Binding
	End      key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Home: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to start"),
		),
		End: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to end"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pageup", "ctrl+u"),
			key.WithHelp("PgUp/ctrl+u", "scroll up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pagedown", "ctrl+d"),
			key.WithHelp("PgDn/ctrl+d", "scroll down"),
		),
	}
}

// Model represents the table model
type Model struct {
	// Data
	headers []string
	rows    [][]string

	// Dimensions
	width  int
	height int

	// Column widths
	columnWidths []int

	// Scrolling
	offsetX int // Horizontal scroll offset
	offsetY int // Vertical scroll offset (row index)

	// Selection
	selectionMode SelectionMode
	selectedRow   int
	selectedCol   int

	// Styling
	theme        Theme
	columnStyles map[int]lipgloss.Style
	rowStyles    map[int]lipgloss.Style

	// Options
	showHeaders bool
	showBorders bool

	// Keymap
	keyMap KeyMap
}

// New creates a new table model
func New() Model {
	return Model{
		headers:       []string{},
		rows:          [][]string{},
		columnWidths:  []int{},
		selectedRow:   0,
		selectedCol:   0,
		selectionMode: SelectionRow,
		showHeaders:   true,
		showBorders:   true,
		columnStyles:  make(map[int]lipgloss.Style),
		rowStyles:     make(map[int]lipgloss.Style),
		theme:         DefaultTheme(),
		keyMap:        DefaultKeyMap(),
	}
}

// SetTheme sets the theme for the table
func (m *Model) SetTheme(theme Theme) {
	m.theme = theme
}

// SetKeyMap sets the key bindings for the table
func (m *Model) SetKeyMap(keyMap KeyMap) {
	m.keyMap = keyMap
}

// ShowHeaders sets whether to show the headers
func (m *Model) ShowHeaders(show bool) {
	m.showHeaders = show
}

// ShowBorders sets whether to show borders between cells
func (m *Model) ShowBorders(show bool) {
	m.showBorders = show
}

// SetHeaders sets the table headers
func (m *Model) SetHeaders(headers []string) {
	m.headers = headers
	m.calculateColumnWidths()
}

// SetRows sets the table rows
func (m *Model) SetRows(rows [][]string) {
	m.rows = rows
	m.calculateColumnWidths()
}

// SetSize sets the viewport size
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetSelectionMode sets how selection works
func (m *Model) SetSelectionMode(mode SelectionMode) {
	m.selectionMode = mode
}

// AddSelectionMode adds a selection mode to the current mode
func (m *Model) AddSelectionMode(mode SelectionMode) {
	m.selectionMode |= mode
}

// RemoveSelectionMode removes a selection mode from the current mode
func (m *Model) RemoveSelectionMode(mode SelectionMode) {
	m.selectionMode &^= mode
}

// ToggleSelectionMode toggles a selection mode
func (m *Model) ToggleSelectionMode(mode SelectionMode) {
	if m.selectionMode&mode != 0 {
		m.RemoveSelectionMode(mode)
	} else {
		m.AddSelectionMode(mode)
	}
}

// HasSelectionMode checks if a selection mode is active
func (m *Model) HasSelectionMode(mode SelectionMode) bool {
	return m.selectionMode&mode != 0
}

// SetColumnStyle sets a custom style for a specific column
func (m *Model) SetColumnStyle(col int, style lipgloss.Style) {
	m.columnStyles[col] = style
}

// SetRowStyle sets a custom style for a specific row
func (m *Model) SetRowStyle(row int, style lipgloss.Style) {
	m.rowStyles[row] = style
}

// calculateColumnWidths calculates the width of each column
func (m *Model) calculateColumnWidths() {
	if len(m.headers) == 0 && len(m.rows) == 0 {
		return
	}

	// Initialize column widths
	numCols := len(m.headers)
	if len(m.rows) > 0 && len(m.rows[0]) > numCols {
		numCols = len(m.rows[0])
	}

	m.columnWidths = make([]int, numCols)

	// Check header widths
	for i, header := range m.headers {
		if len(header) > m.columnWidths[i] {
			m.columnWidths[i] = len(header)
		}
	}

	// Check row widths
	for _, row := range m.rows {
		for i, cell := range row {
			if i < numCols && len(cell) > m.columnWidths[i] {
				m.columnWidths[i] = len(cell)
			}
		}
	}

	// Add padding
	for i := range m.columnWidths {
		m.columnWidths[i] += 2
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			m.moveSelection(-1, 0)
		case key.Matches(msg, m.keyMap.Down):
			m.moveSelection(1, 0)
		case key.Matches(msg, m.keyMap.Left):
			m.moveSelection(0, -1)
		case key.Matches(msg, m.keyMap.Right):
			m.moveSelection(0, 1)
		case key.Matches(msg, m.keyMap.Home):
			m.selectedRow = 0
			m.selectedCol = 0
			m.ensureVisible()
		case key.Matches(msg, m.keyMap.End):
			m.selectedRow = len(m.rows) - 1
			m.selectedCol = len(m.headers) - 1
			m.ensureVisible()
		case key.Matches(msg, m.keyMap.PageUp):
			m.moveSelection(-10, 0)
		case key.Matches(msg, m.keyMap.PageDown):
			m.moveSelection(10, 0)
		}
	}

	return m, nil
}

// moveSelection moves the selection by the given delta
func (m *Model) moveSelection(rowDelta, colDelta int) {
	if len(m.rows) == 0 || len(m.headers) == 0 {
		return
	}

	// Always update position based on the delta
	if rowDelta != 0 {
		m.selectedRow += rowDelta
		if m.selectedRow < 0 {
			m.selectedRow = 0
		}
		if m.selectedRow >= len(m.rows) {
			m.selectedRow = len(m.rows) - 1
		}
	}

	if colDelta != 0 {
		m.selectedCol += colDelta
		if m.selectedCol < 0 {
			m.selectedCol = 0
		}
		if m.selectedCol >= len(m.headers) {
			m.selectedCol = len(m.headers) - 1
		}
	}

	m.ensureVisible()
}

// ensureVisible ensures the selected cell is visible
func (m *Model) ensureVisible() {
	// Calculate visible rows
	visibleRows := m.height
	if m.showHeaders {
		visibleRows -= 1
		if m.showBorders {
			visibleRows -= 1
		}
	}
	if m.showBorders && visibleRows > 1 {
		// Account for borders between rows
		visibleRows = (visibleRows + 1) / 2
	}

	// Vertical scrolling
	if m.selectedRow < m.offsetY {
		m.offsetY = m.selectedRow
	} else if m.selectedRow >= m.offsetY+visibleRows {
		m.offsetY = m.selectedRow - visibleRows + 1
	}

	// Horizontal scrolling
	// Calculate total width needed up to selected column
	totalWidth := 0
	for i := 0; i <= m.selectedCol && i < len(m.columnWidths); i++ {
		totalWidth += m.columnWidths[i]
		if m.showBorders && i > 0 {
			totalWidth += 1 // Border between columns
		}
	}

	// Adjust horizontal offset
	if totalWidth-m.columnWidths[m.selectedCol] < m.offsetX {
		// Selected column is too far left
		m.offsetX = totalWidth - m.columnWidths[m.selectedCol]
		if m.selectedCol > 0 && m.showBorders {
			m.offsetX -= 1
		}
	} else if totalWidth > m.offsetX+m.width {
		// Selected column is too far right
		m.offsetX = totalWidth - m.width
	}

	if m.offsetX < 0 {
		m.offsetX = 0
	}
}

// View renders the table
func (m Model) View() string {
	var lines []string

	// Render headers
	if m.showHeaders {
		headerLine := m.renderRow(m.headers, -1, true)
		lines = append(lines, headerLine)

		if m.showBorders {
			lines = append(lines, m.renderBorder())
		}
	}

	// Calculate visible rows
	availableHeight := m.height - len(lines)
	visibleRows := availableHeight
	if m.showBorders && visibleRows > 1 {
		visibleRows = (visibleRows + 1) / 2
	}

	// Render visible rows
	for i := 0; i < visibleRows && m.offsetY+i < len(m.rows); i++ {
		rowIdx := m.offsetY + i
		rowLine := m.renderRow(m.rows[rowIdx], rowIdx, false)
		lines = append(lines, rowLine)

		// Add border between rows
		if m.showBorders && i < visibleRows-1 && rowIdx < len(m.rows)-1 {
			lines = append(lines, m.renderBorder())
		}
	}

	return strings.Join(lines, "\n")
}

// renderRow renders a single row with proper horizontal scrolling
func (m Model) renderRow(row []string, rowIdx int, isHeader bool) string {
	var result strings.Builder
	currentPos := 0

	for colIdx := 0; colIdx < len(row) && colIdx < len(m.columnWidths); colIdx++ {
		colWidth := m.columnWidths[colIdx]

		// Add border before column (except first)
		if colIdx > 0 && m.showBorders {
			if currentPos >= m.offsetX {
				result.WriteString(m.theme.Border.Render("│"))
			}
			currentPos++
		}

		// Skip columns that are completely before the viewport
		if currentPos+colWidth <= m.offsetX {
			currentPos += colWidth
			continue
		}

		// Stop if we've gone past the viewport
		if currentPos >= m.offsetX+m.width {
			break
		}

		// Render the cell
		cell := row[colIdx]

		// Truncate if needed
		if len(cell) > colWidth-2 {
			cell = cell[:colWidth-3] + "…"
		}

		// Pad the cell
		cell = " " + cell + strings.Repeat(" ", colWidth-len(cell)-1)

		// Handle partial visibility at the start
		startOffset := 0
		if currentPos < m.offsetX {
			startOffset = m.offsetX - currentPos
		}

		// Handle partial visibility at the end
		endOffset := len(cell)
		if currentPos+colWidth > m.offsetX+m.width {
			endOffset = len(cell) - (currentPos + colWidth - m.offsetX - m.width)
		}

		// Extract visible portion
		if startOffset < endOffset {
			visibleCell := cell[startOffset:endOffset]

			// Apply style
			style := m.theme.Cell

			if isHeader {
				style = m.theme.Header
			} else {
				// Check for custom row style
				if rowStyle, ok := m.rowStyles[rowIdx]; ok {
					style = rowStyle
				}

				// Check for custom column style (takes precedence)
				if colStyle, ok := m.columnStyles[colIdx]; ok {
					style = colStyle
				}

				if m.HasSelectionMode(SelectionRow) && rowIdx == m.selectedRow {
					style = m.theme.SelectedRow
				}
				if m.HasSelectionMode(SelectionColumn) && colIdx == m.selectedCol {
					style = m.theme.SelectedCell
				}
				if m.HasSelectionMode(SelectionCell) && rowIdx == m.selectedRow && colIdx == m.selectedCol {
					style = m.theme.SelectedCell
				}
			}

			result.WriteString(style.Render(visibleCell))
		}

		currentPos += colWidth
	}

	return result.String()
}

// renderBorder renders a horizontal border line
func (m Model) renderBorder() string {
	var result strings.Builder
	currentPos := 0

	for colIdx := range m.columnWidths {
		colWidth := m.columnWidths[colIdx]

		// Add intersection before column (except first)
		if colIdx > 0 && m.showBorders {
			if currentPos >= m.offsetX && currentPos < m.offsetX+m.width {
				result.WriteString("┼")
			}
			currentPos++
		}

		// Skip columns that are completely before the viewport
		if currentPos+colWidth <= m.offsetX {
			currentPos += colWidth
			continue
		}

		// Stop if we've gone past the viewport
		if currentPos >= m.offsetX+m.width {
			break
		}

		// Calculate visible portion of the border
		startOffset := 0
		if currentPos < m.offsetX {
			startOffset = m.offsetX - currentPos
		}

		visibleWidth := colWidth - startOffset
		if currentPos+colWidth > m.offsetX+m.width {
			visibleWidth = m.offsetX + m.width - currentPos - startOffset
		}

		if visibleWidth > 0 {
			result.WriteString(strings.Repeat("─", visibleWidth))
		}

		currentPos += colWidth
	}

	return m.theme.Border.Render(result.String())
}

// GetSelectedRow returns the currently selected row index
func (m Model) GetSelectedRow() int {
	return m.selectedRow
}

// GetSelectedColumn returns the currently selected column index
func (m Model) GetSelectedColumn() int {
	return m.selectedCol
}

// GetSelectedCell returns the content of the currently selected cell
func (m Model) GetSelectedCell() (string, bool) {
	if m.selectedRow >= 0 && m.selectedRow < len(m.rows) &&
		m.selectedCol >= 0 && m.selectedCol < len(m.rows[m.selectedRow]) {
		return m.rows[m.selectedRow][m.selectedCol], true
	}

	return "", false
}

// GetCoordinates returns the coordinates of the selected cell
func (m Model) GetCoordinates() (int, int) {
	return m.selectedRow, m.selectedCol
}

// GetSelectionMode returns the current selection mode
func (m Model) GetSelectionMode() SelectionMode {
	return m.selectionMode
}

// ResetSelection resets the selection to the first cell
func (m *Model) ResetSelection() {
	m.selectedRow = 0
	m.selectedCol = 0
	m.offsetX = 0
	m.offsetY = 0
}

// SetSelectedCell sets the selected cell by coordinates
func (m *Model) SetSelectedCell(row, col int) {
	if row < 0 || row >= len(m.rows) || col < 0 || col >= len(m.headers) {
		return
	}

	m.selectedRow = row
	m.selectedCol = col
	m.ensureVisible()
}
