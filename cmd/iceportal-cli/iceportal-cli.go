package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	iceportalapi "github.com/craftamap/iceportal-api"
)

type model struct {
	Speed               float64
	NextStopIn          string
	NextStopTimeActual  string
	NextStopTimePlanned string
	Track               string
	Progress            float64
	NextStationName     string
	Width               int64
	Prog                progress.Model
}

type tickMsg struct{}
type fetchMsg struct {
	status iceportalapi.Status
	trip   iceportalapi.TripInfo
}

func initialModel() model {
	return model{
		Speed:               0,
		NextStopIn:          "0 min",
		NextStopTimeActual:  "00:00",
		NextStopTimePlanned: "00:00",
		NextStationName:     "",
		Track:               "",
		Progress:            0,
		Prog:                progress.NewModel(progress.WithSolidFill("#f20c00")),
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
		m.Prog.Width = msg.Width - 4
		m.Width = int64(msg.Width)
	case tickMsg:
		return m, fetch
	case fetchMsg:
		m.Speed = msg.status.Speed
		stop := findNextArrival(msg.trip)
		timeDiff := stop.Timetable.ActualArrivalTime - msg.status.ServerTime
		m.NextStopIn = fmt.Sprintf("%d min", (timeDiff/1000)/60)
		uPlanned := time.Unix(stop.Timetable.ScheduledArrivalTime/1000, 0)
		m.NextStopTimePlanned = fmt.Sprintf("%02d:%02d", uPlanned.Local().Hour(), uPlanned.Local().Minute())
		uActual := time.Unix(stop.Timetable.ActualArrivalTime/1000, 0)
		m.NextStopTimeActual = fmt.Sprintf("%02d:%02d", uActual.Local().Hour(), uActual.Local().Minute())
		m.Track = stop.Track.Actual
		m.NextStationName = stop.Station.Name
		m.Progress = (float64(msg.trip.Trip.ActualPosition) / float64(msg.trip.Trip.TotalDistance))
		return m, fetch
	}

	return m, nil
}

func (m model) View() string {
	nextStopTime := ""
	if m.NextStopTimeActual == m.NextStopTimePlanned {
		nextStopTime = m.NextStopTimePlanned
	} else {
		striked := lipgloss.NewStyle().Foreground(lipgloss.Color("#f20c00")).Strikethrough(true)
		nextStopTime = fmt.Sprintf("%s %s", striked.Render(m.NextStopTimePlanned), m.NextStopTimeActual)
	}
	statusLineLeft := fmt.Sprintf("%3.f km/h • Next arrival in: %s (%s)", m.Speed, m.NextStopIn, nextStopTime)
	statusLineRight := fmt.Sprintf("Track %s • %s", m.Track, m.NextStationName)
	// we have 2 spaces padding left and right = 4
	spaceInMiddle := int(m.Width) - 4 - len(statusLineLeft) - len(statusLineRight)
	if spaceInMiddle < 0 {
		spaceInMiddle = 1
	}
	spacesInMiddle := strings.Repeat(" ", spaceInMiddle)

	return fmt.Sprintf("\n\n  %s%s%s  \n"+
		"  %s  \n\n", statusLineLeft, spacesInMiddle, statusLineRight, m.Prog.ViewAs(float64(m.Progress)))
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
