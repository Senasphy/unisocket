package app

import(
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/Senasphy/unisocket/internal/snapshot"
	"fmt"
	"strings"
)


var baseStyle = lipgloss.NewStyle().	
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd){

 var cmd tea.Cmd;
 if m.isNamingFile{
	 switch msg := msg.(type){
	 case tea.KeyMsg:
		 switch msg.String(){
		 case "enter":
			 filename := m.filename.Value();
			 if filename == "" { filename = "portList.json"}
				if !strings.HasSuffix(strings.ToLower(filename), ".json") {
											filename += ".json"
									}
				m.isNamingFile = false;
				return m, snapshot.SaveStateToFile(m.table.Rows(), filename)
		 case "esc":
			 m.isNamingFile = false
			 return m, nil
		 }
	 }
	 m.filename, cmd = m.filename.Update(msg)
	 return m, cmd
 }
 switch msg := msg.(type){
 case tickMsg:
	 return m, tea.Batch(m.fetchDataCmd(), tick())
 case DataArrivalMsg:
	 m.table.SetRows(msg)
	 return m, nil

 case tea.KeyMsg:
	 switch msg.String(){
	 case "ctrl+c", "q", "esc":
		 return m, tea.Quit
	 case "s":
						m.isNamingFile = true
            m.filename = textinput.New()
            m.filename.Placeholder = "filename.json"
            m.filename.Focus()
            return m, nil   
	
	 }
	case snapshot.SavedMessageString:
		m.statusMessage = string(msg)
		return m, nil

 }

 m.table, cmd = m.table.Update(msg)
 return m, cmd
}




//THE VIEW METHOD FOR THE MODEL
func (m Model) View() string {

    s := baseStyle.Render(m.table.View()) + "\n"

    if m.statusMessage != "" {
        s += "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(m.statusMessage)
    }
		if m.isNamingFile {
        s+= baseStyle.Render(
            fmt.Sprintf("Enter filename to save snapshot:\n\n%s\n\n(Enter to save, Esc to cancel)", 
            m.filename.View()),
        )
    }
    s += "\n  ↑/↓: Scroll • s: Save as... • q: Quit\n"
    return s
	}

