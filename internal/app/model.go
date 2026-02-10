package app
import(
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/Senasphy/unisocket/internal/connections"
	"time"


)


type tickMsg time.Time
type DataArrivalMsg []table.Row
type SavedMessageString string

type Model struct{
	table table.Model
	statusMessage string
	filename textinput.Model
	lineCount int
	isNamingFile bool
	filters connections.Filters
}


//INIT METHOD FOR THE MODEL
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fetchDataCmd(), tick())
}


func New(t table.Model, filter connections.Filters)tea.Model{
	return Model{
		table: t,
		filters: filter,

	}
}


func tick() tea.Cmd{
	return tea.Tick(time.Second * 2, func (t time.Time) tea.Msg{
					return tickMsg(t)
	})
}


func (m Model) fetchDataCmd() tea.Cmd {
    return func() tea.Msg {
        rows := connections.GetConnections(m.filters)
        return DataArrivalMsg(rows)
    }
}


