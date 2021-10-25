package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	iceportalapi "github.com/craftamap/iceportal-api"
)

type model struct {
	Speed           float64
	NextStopIn      string
	Track           string
	Progress        float64
	NextStationName string
	Prog            progress.Model
}

type tickMsg struct{}
type fetchMsg struct {
	status iceportalapi.Status
	trip   iceportalapi.TripInfo
}

func initialModel() model {
	return model{
		Speed:           0,
		NextStopIn:      "0 min",
		NextStationName: "",
		Track:           "",
		Progress:        0,
		Prog:            progress.NewModel(progress.WithSolidFill("#f20c00")),
	}
}

func (m model) Init() tea.Cmd {
	cmd := tea.Every(5*time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
	return cmd
}

func fetch() tea.Msg {
	status := make(chan iceportalapi.Status)
	trip := make(chan iceportalapi.TripInfo)
	go func() {
		s, _ := iceportalapi.FetchStatus()
		status <- s
	}()
	go func() {
		t, _ := iceportalapi.FetchTrip()
		trip <- t
	}()
	return fetchMsg{
		status: <-status,
		trip:   <-trip,
	}
}

func findNextArrival(tripInfo iceportalapi.TripInfo) iceportalapi.Stop {
	var nextArrival iceportalapi.Stop
	for _, s := range tripInfo.Trip.Stops {
		if !s.Info.Passed {
			nextArrival = s
			break
		}
	}

	return nextArrival
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Prog.Width = msg.Width
	case tickMsg:
		return m, fetch
	case fetchMsg:
		m.Speed = msg.status.Speed
		stop := findNextArrival(msg.trip)
		timeDiff := stop.Timetable.ActualArrivalTime - msg.status.ServerTime
		m.NextStopIn = fmt.Sprintf("%d min", (timeDiff/1000)/60)
		m.Track = stop.Track.Actual
		m.NextStationName = stop.Station.Name
		m.Progress = (float64(msg.trip.Trip.ActualPosition) / float64(msg.trip.Trip.TotalDistance))
		return m, fetch
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("%3.f km/h • Next arrival in: %s • Track %s • %s\n%s", m.Speed, m.NextStopIn, m.Track, m.NextStationName, m.Prog.ViewAs(float64(m.Progress)))
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
