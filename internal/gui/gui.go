package gui

import (
	"log"
	"strings"

	"github.com/burgr033/chaosplayr/internal/file"
	"github.com/burgr033/chaosplayr/internal/mpv"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

var (
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
	borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).BorderForeground(lipgloss.Color("62"))
)

type Item struct {
	title      string
	link       string
	pubDate    string
	author     string
	keywords   string
	summary    string
	duration   string
	identifier string
	isFavorite bool
}

// convert feed into custom item and parse the keywords
func convertFeed(feed []*gofeed.Item) []Item {
	var items []Item
	for n, v := range feed {
		keywords := strings.ReplaceAll(v.ITunesExt.Keywords, " ", "")
		keywords = strings.ReplaceAll(keywords, ",", " #")
		keywords = "#" + keywords
		i := Item{
			title:    v.Title,
			pubDate:  v.Published,
			link:     v.GUID,
			keywords: keywords,
			summary:  v.ITunesExt.Summary,
			author:   v.ITunesExt.Author,
			duration: v.ITunesExt.Duration,
		}
		if n == 584 {
			items = append(items, i)
		} else {
			items = append(items, i)
		}
	}
	return items
}

func (i Item) Title() string { return i.title }
func (i Item) Description() string {
	return i.keywords + "\n" + "🎙️" + i.author + " • 🕛️" + i.duration + " • 🗓️" + i.pubDate
}
func (i Item) FilterValue() string { return i.title + " " + i.author + " " + i.keywords }
func (i Item) Summary() string     { return i.summary }

// model custom model struct
type model struct {
	list    list.Model
	showing bool
	summary string
}

// NewModel gerneates the list from the yaml file and manipulates the items
func NewModel(items []Item) model {
	listItems := []list.Item{}
	for _, v := range items {
		// parsing here
		listItems = append(listItems, v)
	}
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3)
	l := list.New(listItems, delegate, 0, 0)
	l.Title = "chaosplayr"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("*"), key.WithHelp("*", "toggle watchlist")),
			key.NewBinding(key.WithKeys("."), key.WithHelp(".", "information about the video")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "opens stream")),
		}
	}
	return model{list: l}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update handles the default model interactions and also adds the handling of keys
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "*" {
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				log.Printf("would put it on the list %v\n", i.link)
			}
		}
		if msg.String() == "." {
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.showing = !m.showing
				if m.showing {
					m.summary = i.Summary()
				} else {
					m.summary = ""
				}
			}
		}
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				mpv.LaunchMPV(i.link)
			}
			return m, nil
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View overriding the default bubbletea View method for enabling the markdown window
func (m model) View() string {
	if m.showing {
		out, err := glamour.Render(m.summary, "dark")
		if err != nil {
			out = "Error rendering markdown"
		}
		return borderStyle.Render(out)
	}
	return docStyle.Render(m.list.View())
}

// CreateProgram loads the resources, creates the model and launches the actual program
func CreateProgram(rssURL string) (*tea.Program, error) {
	items := convertFeed(file.GetFeed(rssURL))

	m := NewModel(items)
	return tea.NewProgram(m, tea.WithAltScreen()), nil
}
