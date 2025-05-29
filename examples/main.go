package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	table "github.com/ionut-t/gotable"
)

type model struct {
	table       table.Model
	width       int
	height      int
	infoStyle   lipgloss.Style
	showBorders bool
	showHeaders bool
}

func initialModel() model {

	t := table.New()
	t.SetHeaders([]string{
		"ID", "Name", "Department", "Salary", "Location", "Years", "Status",
		"Email", "Phone", "Manager", "Hire Date", "Bonus", "Level", "Project", "Remote", "Performance", "Team", "Extension",
	})

	// Set data
	rows := [][]string{
		{"001", "John Doe", "Engineering", "$120,000", "New York", "5", "Active", "john.doe@example.com", "555-1234", "Alice", "2018-01-15", "$10,000", "Senior", "Apollo", "Yes", "A", "Alpha", "1001"},
		{"002", "Jane Smith", "Marketing", "$95,000", "Los Angeles", "3", "Active", "jane.smith@example.com", "555-2345", "Bob", "2020-03-22", "$7,500", "Mid", "Gemini", "No", "B", "Beta", "1002"},
		{"003", "Bob Johnson", "Sales", "$110,000", "Chicago", "7", "Active", "bob.johnson@example.com", "555-3456", "Charlie", "2016-07-10", "$12,000", "Senior", "Orion", "Yes", "A", "Gamma", "1003"},
		{"004", "Alice Brown", "HR", "$85,000", "Houston", "2", "Active", "alice.brown@example.com", "555-4567", "Diana", "2021-09-01", "$5,000", "Junior", "Pegasus", "No", "C", "Delta", "1004"},
		{"005", "Charlie Wilson", "Engineering", "$135,000", "Seattle", "8", "Active", "charlie.wilson@example.com", "555-5678", "Ethan", "2015-11-30", "$15,000", "Lead", "Apollo", "Yes", "A", "Alpha", "1005"},
		{"006", "Diana Martinez", "Finance", "$105,000", "Boston", "4", "Active", "diana.martinez@example.com", "555-6789", "Fiona", "2019-05-18", "$8,000", "Mid", "Gemini", "No", "B", "Beta", "1006"},
		{"007", "Ethan Davis", "IT", "$115,000", "San Francisco", "6", "On Leave", "ethan.davis@example.com", "555-7890", "George", "2017-02-14", "$11,000", "Senior", "Orion", "Yes", "B", "Gamma", "1007"},
		{"008", "Fiona Garcia", "Legal", "$140,000", "Washington DC", "9", "Active", "fiona.garcia@example.com", "555-8901", "Hannah", "2014-08-25", "$16,000", "Lead", "Pegasus", "Yes", "A", "Delta", "1008"},
		{"009", "George Miller", "Operations", "$90,000", "Phoenix", "1", "Active", "george.miller@example.com", "555-9012", "Ian", "2022-04-12", "$4,000", "Junior", "Apollo", "No", "C", "Alpha", "1009"},
		{"010", "Hannah Lee", "Engineering", "$125,000", "Denver", "5", "Active", "hannah.lee@example.com", "555-0123", "Julia", "2018-12-03", "$10,500", "Senior", "Gemini", "Yes", "A", "Beta", "1010"},
		{"011", "Ian Thompson", "Marketing", "$92,000", "Atlanta", "3", "Active", "ian.thompson@example.com", "555-1122", "Kevin", "2020-06-20", "$7,000", "Mid", "Orion", "No", "B", "Gamma", "1011"},
		{"012", "Julia White", "Sales", "$118,000", "Miami", "6", "Active", "julia.white@example.com", "555-2233", "Linda", "2017-10-09", "$11,500", "Senior", "Pegasus", "Yes", "A", "Delta", "1012"},
		{"013", "Kevin Harris", "Engineering", "$130,000", "Portland", "7", "Active", "kevin.harris@example.com", "555-3344", "Michael", "2016-03-27", "$13,000", "Lead", "Apollo", "Yes", "A", "Alpha", "1013"},
		{"014", "Linda Clark", "HR", "$88,000", "Dallas", "2", "Active", "linda.clark@example.com", "555-4455", "John", "2021-01-19", "$5,500", "Junior", "Gemini", "No", "C", "Beta", "1014"},
		{"015", "Michael Lewis", "Finance", "$112,000", "Philadelphia", "5", "Active", "michael.lewis@example.com", "555-5566", "Jane", "2018-05-05", "$9,500", "Mid", "Orion", "Yes", "B", "Gamma", "1015"},
		{"016", "Natalie Young", "Engineering", "$122,000", "Austin", "4", "Active", "natalie.young@example.com", "555-6677", "Oscar", "2019-08-11", "$10,200", "Senior", "Apollo", "Yes", "A", "Alpha", "1016"},
		{"017", "Oscar King", "Marketing", "$97,000", "San Diego", "2", "Active", "oscar.king@example.com", "555-7788", "Paul", "2021-04-23", "$7,800", "Mid", "Gemini", "No", "B", "Beta", "1017"},
		{"018", "Paula Adams", "Sales", "$113,000", "Detroit", "5", "Active", "paula.adams@example.com", "555-8899", "Quinn", "2018-10-17", "$9,800", "Senior", "Orion", "Yes", "A", "Gamma", "1018"},
		{"019", "Quinn Baker", "HR", "$86,000", "Charlotte", "3", "Active", "quinn.baker@example.com", "555-9900", "Rachel", "2020-12-29", "$5,200", "Junior", "Pegasus", "No", "C", "Delta", "1019"},
		{"020", "Rachel Evans", "Engineering", "$138,000", "San Jose", "9", "Active", "rachel.evans@example.com", "555-1010", "Steve", "2013-07-14", "$16,500", "Lead", "Apollo", "Yes", "A", "Alpha", "1020"},
		{"021", "Steve Foster", "Finance", "$108,000", "Columbus", "6", "Active", "steve.foster@example.com", "555-2020", "Tina", "2017-03-05", "$8,500", "Mid", "Gemini", "No", "B", "Beta", "1021"},
		{"022", "Tina Green", "IT", "$117,000", "Indianapolis", "7", "Active", "tina.green@example.com", "555-3030", "Uma", "2016-09-21", "$11,700", "Senior", "Orion", "Yes", "B", "Gamma", "1022"},
		{"023", "Uma Hall", "Legal", "$142,000", "Jacksonville", "10", "Active", "uma.hall@example.com", "555-4040", "Victor", "2012-05-30", "$17,000", "Lead", "Pegasus", "Yes", "A", "Delta", "1023"},
		{"024", "Victor Ingram", "Operations", "$93,000", "Fort Worth", "2", "Active", "victor.ingram@example.com", "555-5050", "Wendy", "2021-11-13", "$4,500", "Junior", "Apollo", "No", "C", "Alpha", "1024"},
		{"025", "Wendy Johnson", "Engineering", "$127,000", "El Paso", "6", "Active", "wendy.johnson@example.com", "555-6060", "Xavier", "2017-06-18", "$10,800", "Senior", "Gemini", "Yes", "A", "Beta", "1025"},
		{"026", "Xavier King", "IT", "$119,000", "Memphis", "5", "Active", "xavier.king@example.com", "555-7070", "Yvonne", "2018-09-12", "$10,900", "Senior", "Orion", "Yes", "B", "Gamma", "1026"},
		{"027", "Yvonne Lewis", "Legal", "$143,000", "Baltimore", "11", "Active", "yvonne.lewis@example.com", "555-8080", "Zach", "2011-04-28", "$17,500", "Lead", "Pegasus", "Yes", "A", "Delta", "1027"},
		{"028", "Zachary Moore", "Operations", "$94,000", "Louisville", "3", "Active", "zachary.moore@example.com", "555-9090", "Abby", "2021-12-15", "$4,700", "Junior", "Apollo", "No", "C", "Alpha", "1028"},
		{"029", "Abby Nelson", "Engineering", "$129,000", "Milwaukee", "7", "Active", "abby.nelson@example.com", "555-1111", "Ben", "2016-07-19", "$13,200", "Lead", "Apollo", "Yes", "A", "Alpha", "1029"},
		{"030", "Ben Owens", "Finance", "$109,000", "Albuquerque", "6", "Active", "ben.owens@example.com", "555-2222", "Cara", "2017-08-23", "$8,700", "Mid", "Gemini", "No", "B", "Beta", "1030"},
		{"031", "Cara Parker", "IT", "$118,000", "Tucson", "8", "Active", "cara.parker@example.com", "555-3333", "Derek", "2015-10-11", "$11,800", "Senior", "Orion", "Yes", "B", "Gamma", "1031"},
		{"032", "Derek Quinn", "Legal", "$144,000", "Fresno", "12", "Active", "derek.quinn@example.com", "555-4444", "Ella", "2010-03-17", "$18,000", "Lead", "Pegasus", "Yes", "A", "Delta", "1032"},
		{"033", "Ella Roberts", "Operations", "$95,000", "Sacramento", "4", "Active", "ella.roberts@example.com", "555-5555", "Frank", "2021-10-20", "$4,900", "Junior", "Apollo", "No", "C", "Alpha", "1033"},
		{"034", "Frank Scott", "Engineering", "$131,000", "Kansas City", "8", "Active", "frank.scott@example.com", "555-6666", "Grace", "2015-12-29", "$13,500", "Lead", "Apollo", "Yes", "A", "Alpha", "1034"},
		{"035", "Grace Turner", "Marketing", "$98,000", "Mesa", "3", "Active", "grace.turner@example.com", "555-7777", "Hank", "2020-02-14", "$8,000", "Mid", "Gemini", "No", "B", "Beta", "1035"},
		{"036", "Hank Underwood", "Sales", "$114,000", "Atlanta", "6", "Active", "hank.underwood@example.com", "555-8888", "Ivy", "2017-11-05", "$10,000", "Senior", "Orion", "Yes", "A", "Gamma", "1036"},
		{"037", "Ivy Vincent", "HR", "$89,000", "Omaha", "2", "Active", "ivy.vincent@example.com", "555-9999", "Jack", "2021-03-22", "$5,700", "Junior", "Pegasus", "No", "C", "Delta", "1037"},
		{"038", "Jack Walker", "Finance", "$113,000", "Colorado Springs", "5", "Active", "jack.walker@example.com", "555-0001", "Kara", "2018-06-30", "$9,900", "Mid", "Orion", "Yes", "B", "Gamma", "1038"},
		{"039", "Kara Xu", "Engineering", "$123,000", "Raleigh", "4", "Active", "kara.xu@example.com", "555-0002", "Liam", "2019-09-18", "$10,400", "Senior", "Apollo", "Yes", "A", "Alpha", "1039"},
		{"040", "Liam Young", "Marketing", "$99,000", "Long Beach", "2", "Active", "liam.young@example.com", "555-0003", "Mona", "2021-05-25", "$8,200", "Mid", "Gemini", "No", "B", "Beta", "1040"},
		{"041", "Mona Zane", "Sales", "$115,000", "Virginia Beach", "5", "Active", "mona.zane@example.com", "555-0004", "Nate", "2018-11-13", "$10,200", "Senior", "Orion", "Yes", "A", "Gamma", "1041"},
		{"042", "Nate Allen", "HR", "$90,000", "Oakland", "3", "Active", "nate.allen@example.com", "555-0005", "Olga", "2020-01-07", "$5,900", "Junior", "Pegasus", "No", "C", "Delta", "1042"},
		{"043", "Olga Brown", "Finance", "$114,000", "Minneapolis", "6", "Active", "olga.brown@example.com", "555-0006", "Pete", "2017-07-21", "$10,100", "Mid", "Orion", "Yes", "B", "Gamma", "1043"},
		{"044", "Pete Carter", "Engineering", "$132,000", "Tulsa", "9", "Active", "pete.carter@example.com", "555-0007", "Quinn", "2014-10-02", "$14,000", "Lead", "Apollo", "Yes", "A", "Alpha", "1044"},
		{"045", "Quinn Davis", "Marketing", "$100,000", "Arlington", "4", "Active", "quinn.davis@example.com", "555-0008", "Rita", "2019-02-19", "$8,400", "Mid", "Gemini", "No", "B", "Beta", "1045"},
		{"046", "Rita Evans", "Sales", "$116,000", "New Orleans", "7", "Active", "rita.evans@example.com", "555-0009", "Sam", "2016-05-16", "$10,400", "Senior", "Orion", "Yes", "A", "Gamma", "1046"},
		{"047", "Sam Foster", "HR", "$91,000", "Wichita", "2", "Active", "sam.foster@example.com", "555-0010", "Tina", "2021-08-11", "$6,100", "Junior", "Pegasus", "No", "C", "Delta", "1047"},
		{"048", "Tina Grant", "Finance", "$115,000", "Cleveland", "5", "Active", "tina.grant@example.com", "555-0011", "Uma", "2018-03-27", "$10,300", "Mid", "Orion", "Yes", "B", "Gamma", "1048"},
		{"049", "Uma Hall", "Engineering", "$124,000", "Bakersfield", "4", "Active", "uma.hall2@example.com", "555-0012", "Victor", "2019-10-05", "$10,600", "Senior", "Apollo", "Yes", "A", "Alpha", "1049"},
		{"050", "Victor Ingram", "Marketing", "$101,000", "Aurora", "3", "Active", "victor.ingram2@example.com", "555-0013", "Wendy", "2020-06-14", "$8,600", "Mid", "Gemini", "No", "B", "Beta", "1050"},
	}
	t.SetRows(rows)

	// Set initial size (will be updated on window size message)
	t.SetSize(80, 20)

	// Configure custom column styles
	// Highlight salary column
	salaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("82")).
		Bold(true)
	t.SetColumnStyle(3, salaryStyle)

	// Highlight status column
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))
	t.SetColumnStyle(6, statusStyle)

	// Set custom row style for on-leave employee
	onLeaveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true)
	t.SetRowStyle(6, onLeaveStyle) // Row 7 (index 6)

	// Create model
	m := model{
		table: t,
		infoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1),
		showBorders: true,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve space for info line
		m.table.SetSize(msg.Width, msg.Height-3)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "r":
			// Toggle row selection mode
			m.table.ToggleSelectionMode(table.SelectionRow)

		case "c":
			// Toggle column selection mode
			m.table.ToggleSelectionMode(table.SelectionColumn)

		case "x":
			// Toggle cell selection mode
			m.table.ToggleSelectionMode(table.SelectionCell)

		case "n":
			// Clear all selection modes
			m.table.SetSelectionMode(table.SelectionOff)

		case "a":
			// Select all (row + column)
			m.table.SetSelectionMode(table.SelectionRow | table.SelectionColumn)

		case "b":
			// Toggle borders
			m.showBorders = !m.showBorders
			m.table.ShowBorders(m.showBorders)

		case "h":
			// Toggle headers
			m.showHeaders = !m.showHeaders
			m.table.ShowHeaders(m.showHeaders)
			m.table.SetSize(m.width, m.height-3)
		}
	}

	// Update table
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string {
	// Main table view
	tableView := m.table.View()

	// Info line
	var selectionModes []string
	if m.table.HasSelectionMode(table.SelectionRow) {
		selectionModes = append(selectionModes, "Row")
	}
	if m.table.HasSelectionMode(table.SelectionColumn) {
		selectionModes = append(selectionModes, "Column")
	}
	if m.table.HasSelectionMode(table.SelectionCell) {
		selectionModes = append(selectionModes, "Cell")
	}

	selectionMode := "None"
	if len(selectionModes) > 0 {
		selectionMode = strings.Join(selectionModes, "+")
	}

	selectedCell, _ := m.table.GetSelectedCell()

	info := fmt.Sprintf(
		"Selection: %s | Row: %d, Col: %d | Cell: %s | Keys: (r)ow (c)ol (x)cell (a)ll (n)one (b)orders (h)eaders (q)uit",
		selectionMode,
		m.table.GetSelectedRow()+1,
		m.table.GetSelectedColumn()+1,
		selectedCell,
	)

	return tableView + "\n" + m.infoStyle.Render(info)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
