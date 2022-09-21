package main

import (
	"log"
	"math"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

var agentSlice []AgentPosition

type Agent struct {
	Name   string
	Status string
}

type AgentPosition struct {
	Name string
	PosX float64
	PosY float64
}

func main() {
	http.Handle("/circle", http.HandlerFunc(circle))
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func circle(w http.ResponseWriter, req *http.Request) {
	statusColor := ""
	var sin float64
	var cos float64
	s := svg.New(w)
	height := 800
	widht := 800
	aperture := len(Agents())
	angle := 0
	w.Header().Set("Content-Type", "image/svg+xml")
	s.Start(height, widht)
	for _, agent := range Agents() {
		switch agent.Status {
		case "online":
			statusColor = "fill:green;stroke:green"
		case "offline":
			statusColor = "fill:red;stroke:red"
		case "busy":
			statusColor = "fill:yellow;stroke:yellow"
		case "blue":
			statusColor = "fill:blue;stroke:blue"
		}
		if angle == 0 {
			sin = math.Sin(0)
			cos = math.Cos(0)
		} else {
			sin = math.Sin(float64(2 * math.Pi * float64(angle) / float64(aperture)))
			cos = math.Cos(float64(2 * math.Pi * float64(angle) / float64(aperture)))
		}
		sin = sin*350 + 400
		cos = cos*350 + 400
		s.Circle(int(sin), int(cos), 10, statusColor)
		appendAgents(agent.Name, sin, cos)
		angle++
	}
	s.Circle(400, 400, 10, "fill:black;stroke:black")
	s.End()
}

func Agents() []Agent {
	return []Agent{
		{"agent01", "online"},
		{"agent02", "busy"},
		{"agent03", "offline"},
		{"agent04", "busy"},
		{"agent05", "online"},
		{"agent06", "busy"},
		{"agent07", "offline"},
		{"agent08", "busy"},
	}
}

func appendAgents(name string, x float64, y float64) {
	agentSlice = append(agentSlice, AgentPosition{name, x, y})
}
