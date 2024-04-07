package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Pulse int

const (
	PULSE_LOW  Pulse = iota
	PULSE_HIGH Pulse = iota
)

type PendingPulse struct {
	fromModule        string
	pulse             Pulse
	destinationModule string
}

func CreatePendingPulses(fromModule string, pulse Pulse, destinationModules []string) []PendingPulse {
	pendingPulses := make([]PendingPulse, 0, len(destinationModules))
	for _, toModule := range destinationModules {
		pendingPulses = append(pendingPulses, PendingPulse{fromModule, pulse, toModule})
	}
	return pendingPulses
}

type Module interface {
	GetName() string
	GetDestinationModules() []string
	RegisterInput(string)
	SendPulse(Pulse, string) []PendingPulse
}

type Broadcaster struct {
	destinationModules []string
}

func (broadcaster *Broadcaster) GetName() string {
	return "broadcaster"
}

func (broadcaster *Broadcaster) GetDestinationModules() []string {
	return broadcaster.destinationModules
}

func (broadcaster *Broadcaster) RegisterInput(inputModule string) {
	// Do nothing.
}

func (broadcaster *Broadcaster) SendPulse(pulse Pulse, fromModule string) []PendingPulse {
	return CreatePendingPulses("broadcaster", pulse, broadcaster.destinationModules)
}

type FlipFlop struct {
	name               string
	on                 bool
	destinationModules []string
}

func (flipFlop *FlipFlop) GetName() string {
	return flipFlop.name
}

func (flipFlop *FlipFlop) GetDestinationModules() []string {
	return flipFlop.destinationModules
}

func (flipFlop *FlipFlop) RegisterInput(inputModule string) {
	// Do nothing.
}

func (flipFlop *FlipFlop) SendPulse(pulse Pulse, fromModule string) []PendingPulse {
	if pulse == PULSE_LOW {
		newPulse := PULSE_LOW
		if !flipFlop.on {
			newPulse = PULSE_HIGH
		}
		flipFlop.on = !flipFlop.on
		return CreatePendingPulses(flipFlop.name, newPulse, flipFlop.destinationModules)
	}

	return make([]PendingPulse, 0, 0)
}

type Conjunction struct {
	name               string
	memory             map[string]Pulse
	destinationModules []string
}

func (conjunction *Conjunction) GetName() string {
	return conjunction.name
}

func (conjunction *Conjunction) GetDestinationModules() []string {
	return conjunction.destinationModules
}

func (conjunction *Conjunction) RegisterInput(inputModule string) {
	conjunction.memory[inputModule] = PULSE_LOW
}

func (conjunction *Conjunction) SendPulse(pulse Pulse, fromModule string) []PendingPulse {
	conjunction.memory[fromModule] = pulse

	newPulse := PULSE_LOW
	for _, pulse := range conjunction.memory {
		if pulse == PULSE_LOW {
			newPulse = PULSE_HIGH
		}
	}
	return CreatePendingPulses(conjunction.name, newPulse, conjunction.destinationModules)
}

func getModules(scanner *bufio.Scanner) map[string]Module {
	modules := make(map[string]Module)
	for scanner.Scan() {
		line := scanner.Text()
		splitModule := strings.Split(line, " -> ")
		if len(splitModule) != 2 {
			log.Fatal("Invalid line: ", line)
		}

		destinationModules := strings.Split(splitModule[1], ", ")

		if splitModule[0] == "broadcaster" {
			modules[splitModule[0]] = &Broadcaster{destinationModules}
			continue
		}

		moduleType := splitModule[0][0]
		name := splitModule[0][1:]
		switch moduleType {
		case '%':
			modules[name] = &FlipFlop{name, false, destinationModules}
		case '&':
			modules[name] = &Conjunction{name, make(map[string]Pulse), destinationModules}
		default:
			log.Fatal("Invalid line: ", line)
		}
	}

	for _, module := range modules {
		for _, destinationModuleName := range module.GetDestinationModules() {
			destinationModule, found := modules[destinationModuleName]
			// Ignore missing outputs.
			if found {
				destinationModule.RegisterInput(module.GetName())
			}
		}
	}

	return modules
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buttonPresses := 0
	pulseQueue := make(chan PendingPulse, 1024)
	modules := getModules(bufio.NewScanner(file))
	dhHighPulses := make(map[string]int)

	for {
		buttonPresses += 1
		pulseQueue <- PendingPulse{"button", PULSE_LOW, "broadcaster"}

	PULSE_LOOP:
		for {
			select {
			case pendingPulse := <-pulseQueue:
				// Via inspection of the puzzle, `dh` is a conjunction with 4
				// conjunction inputs. We record the first time each one sends
				// a high signal; once they have all sent a high signal, we
				// break out of the loop and multiply them all together to get
				// the answer.
				if pendingPulse.destinationModule == "dh" && pendingPulse.pulse == PULSE_HIGH {
					_, found := dhHighPulses[pendingPulse.fromModule]
					if !found {
						dhHighPulses[pendingPulse.fromModule] = buttonPresses
					}
				}

				module, found := modules[pendingPulse.destinationModule]
				if found {
					newPulses := module.SendPulse(pendingPulse.pulse, pendingPulse.fromModule)
					for _, pulse := range newPulses {
						pulseQueue <- pulse
					}
				}
			default:
				break PULSE_LOOP
			}
		}
		if len(dhHighPulses) == 4 {
			break
		}
	}

	total := 1
	for _, count := range dhHighPulses {
		total = total * count
	}
	fmt.Println(total)
}
