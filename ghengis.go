package ghengis

import (
	"github.com/KaptajnKold/antwar"
	"math/rand"
)

type pos struct {
	x, y int
}

var home = pos{0, 0}

type ghengis struct {
	pos, destination pos
	hasSeenFood      bool
}

func (here pos) horizontalDirectionTo(there pos) antwar.Action {
	if here.x < there.x {
		return antwar.EAST
	} else if here.x > there.x {
		return antwar.WEST
	}
	return antwar.HERE
}

func (here pos) verticalDirectionTo(there pos) antwar.Action {
	if here.y < there.y {
		return antwar.NORTH
	} else if here.y > there.y {
		return antwar.SOUTH
	}
	return antwar.HERE
}

func (p pos) inverse() pos {
	return pos{-p.x, -p.y}
}

func flipACoin() bool {
	if rand.Intn(2) == 0 {
		return true
	}
	return false
}

func divisibleByThree(n int) bool {
	return n%3 == 0
}

func (me *ghengis) setRandomDestination() {
	me.destination.x = 3 * (me.pos.x + rand.Intn(100) - 50)
	me.destination.y = 3 * (me.pos.y + rand.Intn(100) - 50)
}

func (here pos) directionTo(there pos) antwar.Action {
	if here.x == there.x {
		return here.verticalDirectionTo(there)
	}
	if here.y == there.y {
		return here.horizontalDirectionTo(there)
	}
	if (divisibleByThree(here.x) && divisibleByThree(here.y)) || (!divisibleByThree(here.x) && !divisibleByThree(here.y)) {
		if flipACoin() {
			return here.horizontalDirectionTo(there)
		} else {
			return here.verticalDirectionTo(there)
		}
	}
	if divisibleByThree(here.x) {
		return here.verticalDirectionTo(there)
	}
	return here.horizontalDirectionTo(there)
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (here *pos) distanceTo(there pos) int {
	return Abs(here.x-there.x) + Abs(here.y-there.y)
}

func (me *ghengis) directionHome() antwar.Action {
	return me.pos.directionTo(home)
}

func (me *ghengis) updatePosition(decision antwar.Action) {
	switch decision {
	case antwar.EAST:
		me.pos.x++
	case antwar.SOUTH:
		me.pos.y--
	case antwar.WEST:
		me.pos.x--
	case antwar.NORTH:
		me.pos.y++
	}
}

func (me *ghengis) directionToDestination() antwar.Action {
	return me.pos.directionTo(me.destination)
}

func (me *ghengis) tilePositions(env *antwar.Tile) map[*antwar.Tile]pos {
	return map[*antwar.Tile]pos{
		env.Here():  me.pos,
		env.North(): pos{me.pos.x, me.pos.y + 1},
		env.East():  pos{me.pos.x + 1, me.pos.y},
		env.South(): pos{me.pos.x, me.pos.y - 1},
		env.West():  pos{me.pos.x - 1, me.pos.y},
	}
}

func (me *ghengis) distanceHome() int {
	return me.pos.distanceTo(home)
}

func (me *ghengis) distanceToDestinationAndHome() int {
	return me.pos.distanceTo(me.destination) + me.destination.distanceTo(home)
}

func moreFoodThanAnts(tile *antwar.Tile) bool {
	return tile.FoodCount() > tile.AntCount() // && tile.Team() != "ghengis"
}

func (me *ghengis) Decide(env *antwar.Tile, brains []antwar.AntBrain) (decision antwar.Action, bringFood bool) {
	tilePositions := me.tilePositions(env)
	bringFood = false
	if len(brains) > 0 {
		for _, brain := range brains {
			other, ok := brain.(*ghengis)
			if !ok {
				panic("Not a ghengis")
			}

			// Find the shortest way home
			if other.distanceHome() < me.distanceHome() {
				me.pos = other.pos
			}

			if !me.hasSeenFood && other.hasSeenFood {
				me.destination = other.destination
				me.hasSeenFood = true
			}

			if me.hasSeenFood && me.distanceToDestinationAndHome() > other.distanceToDestinationAndHome() {
				me.destination = other.destination
			}
		}
	}
	for tile, tilePos := range tilePositions {
		if moreFoodThanAnts(tile) {
			me.hasSeenFood = true
			me.destination = tilePos
		}
	}
	if env.FoodCount() > 0 {
		bringFood = true
		decision = me.directionHome()
	} else {
		decision = me.directionToDestination()
		if decision == antwar.HERE {
			me.setRandomDestination()
			me.hasSeenFood = false
			decision = me.directionToDestination()
		}
	}
	me.updatePosition(decision)
	return
}

func Spawn() antwar.AntBrain { return new(ghengis) }
