package main

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"os"

	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/google/uuid"
)

var PosAX float64
var PosAY float64
var PosBX float64
var PosBY float64
var globalmm MeshSnapshot

func main() {
	assignValues() // NOTE this assign its a test example the func should be removed in order to receive param when used as external package
	fmt.Println(UTF8(globalmm))
	fmt.Println(SVG(globalmm))
}

func assignValues() {
	agent01uuid := uuid.Must(uuid.NewRandom())
	agent02uuid := uuid.Must(uuid.NewRandom())
	agent03uuid := uuid.Must(uuid.NewRandom())
	agent04uuid := uuid.Must(uuid.NewRandom())
	agent05uuid := uuid.Must(uuid.NewRandom())
	agent06uuid := uuid.Must(uuid.NewRandom())

	agent01name := "agent01"
	agent02name := "agent02"
	agent03name := "agent03"
	agent04name := "agent04"
	agent05name := "agent05"
	agent06name := "agent06"

	globalmm.AgentNames = map[uuid.UUID]string{
		agent01uuid: agent01name,
		agent02uuid: agent02name,
		agent03uuid: agent03name,
		agent04uuid: agent04name,
		agent05uuid: agent05name,
		agent06uuid: agent06name}

	globalmm.Agents = map[uuid.UUID]string{
		agent01uuid: "online",
		agent02uuid: "offline",
		agent03uuid: "degraded",
		agent04uuid: "online",
		agent05uuid: "online",
		agent06uuid: "offline"}

	globalmm.Connections = map[uuid.UUID][]uuid.UUID{
		agent01uuid: {agent02uuid, agent03uuid},
		agent02uuid: {agent01uuid},
		agent03uuid: {agent01uuid},
		agent04uuid: {agent06uuid},
		agent05uuid: {},
		agent06uuid: {agent04uuid},
	}
}

type MeshSnapshot struct {
	// mapping of an agent id to a name
	AgentNames map[uuid.UUID]string
	// AgentNames map[string]string

	// mapping of the agent id to a status string
	Agents map[uuid.UUID]string

	// mapping of an agent and its connections to other agents
	Connections map[uuid.UUID][]uuid.UUID
}

// Returns a simple UTF8 representation of the mesh map
func UTF8(mm MeshSnapshot) string {
	output := ""
	for id, name := range mm.AgentNames {
		for id2, status := range mm.Agents {
			if id == id2 {
				output += "---------------------------------------------\n"
				output += fmt.Sprintf("Agent Name:\n\t%s (%s)\nAgent ID:\n\t%s\n", name, status, id)
				for id3, connections := range mm.Connections {
					if id2 == id3 {
						output += "Agent Connections:\n"
						for _, conn := range connections {
							for id4, name2 := range mm.AgentNames {
								if conn == id4 {
									output += fmt.Sprintf("\t%s\n", name2)
								}
							}
						}
					}
				}
			}
		}
		output += "---------------------------------------------\n\n\n"
	}
	return output
}

// Returns an SVG string that can be embedded in an HTML file
func SVG(mm MeshSnapshot) string {
	ac := make(map[uuid.UUID][]float64)
	f, err := os.Create("agents_svg_file.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := svg.New(f)
	s.Start(800, 800)

	defineStyle(s)
	assignCoordinates(ac)
	drawConnections(s, ac)
	drawAgents(ac, s)
	s.End()

	b, err := os.ReadFile("agents_svg_file.svg") // just pass the file name
    if err != nil {
        panic(err)
    }
	return string(b)
}

func defineStyle(s *svg.SVG) {
	s.Style("text/css", "text {text-anchor:middle;font:1.2em effra,sans-serif;}")
	s.Style("text/css", ".agentonline {fill:green}")
	s.Style("text/css", ".agentoffline {fill:red}")
	s.Style("text/css", ".agentdegraded {fill:yellow}")
	s.Style("text/css", "text:hover {font-size:xxx-large;fill:black}")
}

func assignCoordinates(ac map[uuid.UUID][]float64) {
	var sin float64
	var cos float64
	angle := 0
	aperture := len(globalmm.Agents)
	for uuid := range globalmm.AgentNames {
		if angle == 0 {
			sin = math.Sin(0)
			cos = math.Cos(0)
		} else {
			sin = math.Sin(float64(2 * math.Pi * float64(angle) / float64(aperture)))
			cos = math.Cos(float64(2 * math.Pi * float64(angle) / float64(aperture)))
		}
		sin = sin*350 + 400
		cos = cos*350 + 400
		appendAgents(ac, uuid, sin, cos)
		angle++
	}
}

func drawConnections(s *svg.SVG, ac map[uuid.UUID][]float64) {
	// iterate over connections
	for uuid1, connections := range globalmm.Connections {
		// filter offline agents (from)
		for uuidstatus, status := range globalmm.Agents {
			if uuid1 == uuidstatus && status != "offline" {
				for uuidAC, coordinates := range ac {
					if uuid1 == uuidAC {
						PosAX = coordinates[0]
						PosAY = coordinates[1]
					}
				}
				for _, uuidto := range connections {
					// filter offline agents (to)
					for connectionstatus, status := range globalmm.Agents {
						if uuidto == connectionstatus && status != "offline" {
							for uuidAC1, coordinates := range ac {
								if uuidto == uuidAC1 {
									PosBX = coordinates[0]
									PosBY = coordinates[1]
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
}

func appendAgents(ac map[uuid.UUID][]float64, uuid uuid.UUID, sin, cos float64) {
	ac[uuid] = append(ac[uuid], sin)
	ac[uuid] = append(ac[uuid], cos)
}

func drawAgents(ac map[uuid.UUID][]float64, s *svg.SVG) {
	agentClass := ""
	for uuid1, position := range ac {
		for uuid2, name := range globalmm.AgentNames {
			if uuid1 == uuid2 {
				for uuid3, status := range globalmm.Agents {
					if uuid2 == uuid3 {
						switch status {
						case "online":
							agentClass = "class=\"agentonline\""
						case "offline":
							agentClass = "class=\"agentoffline\""
						case "degraded":
							agentClass = "class=\"agentdegraded\""
						}
						s.Text(int(position[0]), int(position[1]), name, agentClass)
					}
				}
			}
		}
	}
}

func randomHex() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
