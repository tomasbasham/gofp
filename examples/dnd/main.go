// DnD is a build system simulation using the [gofp.Either] monad to represent
// different character builds.
package main

import (
	"fmt"
	"strings"

	"github.com/tomasbasham/gofp"
)

// Mage represents a build that is more focused on magical abilities.
type Mage struct {
	Name   string
	Mana   int
	Spells []string
}

// Warrior represents a build that is more focused on physical combat.
type Warrior struct {
	Name    string
	Stamina int
	Weapon  string
}

// Build represents a build that can be either a mage or warrior.
type Build = gofp.Either[Mage, Warrior]

func main() {
	traits := []string{"wise", "mystical", "quick"}

	build := classify(traits)
	build = specialize(build)
	build = nameCharacter("Aelin", build)

	fmt.Println(describe(build))
}

func classify(traits []string) Build {
	magic := 0
	martial := 0

	for _, t := range traits {
		switch strings.ToLower(t) {
		case "studious", "wise", "mystical", "intellectual":
			magic++
		case "strong", "agile", "tough", "quick":
			martial++
		}
	}

	if magic >= martial {
		return gofp.Left[Mage, Warrior](Mage{
			Mana:   100,
			Spells: []string{"Light"},
		})
	}

	return gofp.Right[Mage](Warrior{
		Stamina: 100,
		Weapon:  "Fists",
	})
}

func specialize(build Build) Build {
	return build.MapLeft(func(m Mage) Mage {
		m.Mana += 50
		m.Spells = append(m.Spells, "Fireball")
		return m
	}).Map(func(w Warrior) Warrior {
		w.Stamina += 30
		w.Weapon = "Sword"
		return w
	})
}

func nameCharacter(name string, build Build) Build {
	return build.MapLeft(func(m Mage) Mage {
		m.Name = name
		return m
	}).Map(func(w Warrior) Warrior {
		w.Name = name
		return w
	})
}

func describe(build Build) string {
	return gofp.EitherFold(
		build,
		func(m Mage) string {
			return fmt.Sprintf("Mage %s with %d mana and spells: %v", m.Name, m.Mana, m.Spells)
		},
		func(w Warrior) string {
			return fmt.Sprintf("Warrior %s with %d stamina and weapon: %s", w.Name, w.Stamina, w.Weapon)
		},
	)
}
