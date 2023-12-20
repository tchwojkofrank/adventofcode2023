package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func readInput(fname string) string {
	content, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	return string(content)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Enter input file name.\n")
		return
	}
	params := os.Args[1]
	inputName := strings.Split(params, " ")[0]
	text := readInput(inputName)
	start := time.Now()
	run(text)
	end := time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
	start = time.Now()
	run2(text)
	end = time.Now()
	fmt.Printf("Running time: %v\n", end.Sub(start))
}

const (
	Low  = 0
	High = 1
)

type Module interface {
	Execute(pulse int, source string) []Pulse
	Name() string
	Type() string
	Destinations() []string
}

type FlipFlop struct {
	name         string
	t            string
	state        bool
	destinations []string
}

type Conjunction struct {
	name         string
	t            string
	state        map[string]int
	destinations []string
}

type Pulse struct {
	source      string
	destination string
	pulse       int
}

type Broadcast struct {
	name         string
	t            string
	destinations []string
}

type Machine struct {
	modules        map[string]Module
	lowPulseCount  int
	highPulseCount int
}

func (ff *FlipFlop) Execute(pulse int, source string) []Pulse {
	if pulse == High {
		return nil
	}
	pulses := make([]Pulse, 0, len(ff.destinations))
	var newPulse int
	if !ff.state {
		newPulse = High
	} else {
		newPulse = Low
	}
	ff.state = !ff.state
	for _, d := range ff.destinations {
		pulses = append(pulses, Pulse{ff.name, d, newPulse})
	}
	return pulses
}

func (ff *FlipFlop) Name() string {
	return ff.name
}

func (ff *FlipFlop) Destinations() []string {
	return ff.destinations
}

func (ff *FlipFlop) Type() string {
	return ff.t
}

func (conj *Conjunction) Execute(pulse int, source string) []Pulse {
	conj.state[source] = pulse
	newPulses := make([]Pulse, 0)
	for _, v := range conj.state {
		if v == Low {
			for _, d := range conj.destinations {
				newPulses = append(newPulses, Pulse{conj.name, d, High})
			}
			return newPulses
		}
	}
	for _, d := range conj.destinations {
		newPulses = append(newPulses, Pulse{conj.name, d, Low})
	}
	return newPulses
}

func (conj *Conjunction) Name() string {
	return conj.name
}

func (conj *Conjunction) Destinations() []string {
	return conj.destinations
}

func (conj *Conjunction) Type() string {
	return conj.t
}

func (b *Broadcast) Execute(pulse int, source string) []Pulse {
	newPulses := make([]Pulse, 0)
	for _, d := range b.destinations {
		newPulses = append(newPulses, Pulse{b.name, d, pulse})
	}
	return newPulses
}

func (b *Broadcast) Name() string {
	return b.name
}

func (b *Broadcast) Destinations() []string {
	return b.destinations
}

func (b *Broadcast) Type() string {
	return b.t
}

func getModule(moduleString string) Module {
	parts := strings.Split(moduleString, " -> ")
	destinations := strings.Split(parts[1], ", ")
	switch parts[0][0] {
	case '%':
		var f FlipFlop
		f.name = parts[0][1:]
		f.t = "%"
		f.state = false
		f.destinations = destinations
		return &f
	case '&':
		var c Conjunction
		c.destinations = destinations
		c.name = parts[0][1:]
		c.t = "&"
		c.state = make(map[string]int)
		return &c
	case 'b':
		var b Broadcast
		b.destinations = destinations
		b.name = parts[0]
		b.t = "b"
		return &b
	}
	return nil
}

func sendPulses(m *Machine, pulses []Pulse) []Pulse {
	nextPulses := make([]Pulse, 0)
	for _, pulse := range pulses {
		if pulse.pulse == Low {
			m.lowPulseCount++
		} else {
			m.highPulseCount++
		}
		// // pulseName := "low"
		// if pulse.pulse == High {
		// 	pulseName = "high"
		// }
		// fmt.Printf("%v -%v -> %v\n", pulse.source, pulseName, pulse.destination)
		module, ok := m.modules[pulse.destination]
		if ok {
			pulses = module.Execute(pulse.pulse, pulse.source)
			nextPulses = append(nextPulses, pulses...)
		}
	}
	return nextPulses
}

func pushButton(m *Machine, p int) (int, int) {
	pulse := Pulse{"", "broadcaster", p}
	pulses := make([]Pulse, 1)
	pulses[0] = pulse
	for len(pulses) > 0 {
		pulses = sendPulses(m, pulses)
	}
	return m.lowPulseCount, m.highPulseCount
}

func run(input string) string {
	moduleList := strings.Split(input, "\n")
	var machine Machine
	machine.modules = make(map[string]Module)
	machine.lowPulseCount = 0
	machine.highPulseCount = 0
	for _, v := range moduleList {
		module := getModule(v)
		machine.modules[module.Name()] = module
	}
	for name, module := range machine.modules {
		for _, destination := range module.Destinations() {
			module, ok := machine.modules[destination]
			if ok && module.Type() == "&" {
				machine.modules[destination].(*Conjunction).state[name] = Low
			}
		}
	}
	for i := 0; i < 1000; i++ {
		pushButton(&machine, Low)
	}
	fmt.Printf("Low: %d, High: %d\n", machine.lowPulseCount, machine.highPulseCount)
	fmt.Printf("Result: %d\n", machine.highPulseCount*machine.lowPulseCount)
	return fmt.Sprintf("%d", machine.highPulseCount*machine.lowPulseCount)
}

type Final struct {
	done bool
}

func (f *Final) Execute(pulse int, source string) []Pulse {
	if !f.done && pulse == Low {
		f.done = true
	}
	return nil
}

func (f *Final) Name() string {
	return "rx"
}

func (f *Final) Destinations() []string {
	return nil
}

func (f *Final) Type() string {
	return "f"
}

func printLS(m *Machine) {
	ls := m.modules["ls"].(*Conjunction)
	fmt.Printf("ls: %v\n", ls.state)
}

func run2(input string) string {
	moduleList := strings.Split(input, "\n")
	var machine Machine
	machine.modules = make(map[string]Module)
	machine.lowPulseCount = 0
	machine.highPulseCount = 0
	for _, v := range moduleList {
		module := getModule(v)
		machine.modules[module.Name()] = module
	}
	var f Final
	machine.modules[f.Name()] = &f
	for name, module := range machine.modules {
		for _, destination := range module.Destinations() {
			module, ok := machine.modules[destination]
			if ok && module.Type() == "&" {
				machine.modules[destination].(*Conjunction).state[name] = Low
			}
		}
	}
	done := false
	count := 0
	for !done {
		pushButton(&machine, Low)
		count++
		if count%100000 == 0 {
			fmt.Printf("Button pushes: %d\n", count)
			printLS(&machine)

		}
		done = machine.modules[f.Name()].(*Final).done
	}
	fmt.Printf("Final Button pushes: %d\n", count)
	return fmt.Sprintf("%d", count)
}
