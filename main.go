package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

var agentSlice []AgentPosition
var PosAX float64
var PosAY float64
var PosBX float64
var PosBY float64

type Agent struct {
	Name   string
	Status string
}

type AgentPosition struct {
	Name string
	PosX float64
	PosY float64
}

type Relation struct {
	AgentA string
	AgentB string
}

func main() {
	http.Handle("/circle", http.HandlerFunc(circle))
	http.Handle("/text", http.HandlerFunc(text))
	err := http.ListenAndServe("localhost:2003", nil)
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
		case "degraded":
			statusColor = "fill:yellow;stroke:yellow"
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
		s.Circle(int(sin), int(cos), 50, statusColor)
		appendAgents(agent.Name, sin, cos)
		angle++
	}

	s.Circle(400, 400, 10, "fill:black;stroke:black") // this is to have a visual center reference

	s.End()
}

func text(w http.ResponseWriter, req *http.Request) {
	agentSlice = nil
	var sin float64
	var cos float64
	s := svg.New(w)
	height := 800
	widht := 800
	aperture := len(Agents())
	angle := 0
	w.Header().Set("Content-Type", "image/svg+xml")
	s.Start(height, widht)
	defineStyle(s)
	for _, agent := range Agents() {
		if angle == 0 {
			sin = math.Sin(0)
			cos = math.Cos(0)
		} else {
			sin = math.Sin(float64(2 * math.Pi * float64(angle) / float64(aperture)))
			cos = math.Cos(float64(2 * math.Pi * float64(angle) / float64(aperture)))
		}
		sin = sin*350 + 400
		cos = cos*350 + 400
		appendAgents(agent.Name, sin, cos)
		angle++
	}
	drawConnections(s)
	drawAgents(s)
	s.Style("text/css", "#textgreen:hover {fill: black;}")
	// s.Circle(400, 400, 10, "fill:black;stroke:black") // this is to have a visual center reference
	s.End()
}

func defineStyle(s *svg.SVG) {
	s.Style("text/css", "text {text-anchor:middle;font:1.2em effra,sans-serif;}")
	s.Style("text/css", ".agentonline {fill:green}")
	s.Style("text/css", ".agentoffline {fill:red}")
	s.Style("text/css", ".agentdegraded {fill:yellow}")
	s.Style("text/css", "text:hover {font-size:xxx-large;fill:black}")
}

func drawAgents(s *svg.SVG) {

	agentClass := ""
	for _, element := range agentSlice {
		for _, agent := range Agents() {
			if element.Name == agent.Name {
				switch agent.Status {
				case "online":
					agentClass = "class=\"agentonline\""
				case "offline":
					agentClass = "class=\"agentoffline\""
				case "degraded":
					agentClass = "class=\"agentdegraded\""
				}
				s.Text(int(element.PosX), int(element.PosY), element.Name, agentClass)

			}
		}
	}
}

func drawConnections(s *svg.SVG) {
	for _, relation := range Relations() {
		for _, agentA := range Agents() {
			if agentA.Name == relation.AgentA && agentA.Status != "offline" {
				for _, agentB := range Agents() {
					if agentB.Name == relation.AgentB && agentB.Status != "offline" {
						for _, agents := range agentSlice {
							if agents.Name == relation.AgentA {
								PosAX = agents.PosX
								PosAY = agents.PosY
							}
						}
						for _, agents := range agentSlice {
							if agents.Name == relation.AgentB {
								PosBX = agents.PosX
								PosBY = agents.PosY
							}
						}
						color, _ := randomHex()
						s.Line(int(PosAX), int(PosAY), int(PosBX), int(PosBY), fmt.Sprintf("stroke-width:2;stroke:#%s", color))
					}
				}
			}
		}
	}
}

func Agents() []Agent {
	return []Agent{
		{"agent01", "online"},
		{"agent02", "online"},
		{"agent03", "degraded"},
		{"agent04", "offline"},
		{"agent05", "offline"},
		{"agent06", "online"},
	}
}

func Relations() []Relation {
	return []Relation{
		{"agent01", "agent02"},
		{"agent01", "agent03"},
		{"agent01", "agent04"},
		{"agent02", "agent03"},
		{"agent02", "agent05"},
		{"agent03", "agent06"},
		{"agent03", "agent04"},
		{"agent04", "agent06"},
		{"agent05", "agent01"},
		{"agent06", "agent02"},
	}
}

func appendAgents(name string, x float64, y float64) {
	agentSlice = append(agentSlice, AgentPosition{name, x, y})
}

func randomHex() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
