package player

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/minimoke/gorogue/world"
)

type Character struct {
	PlayerID     int
	Name         string
	Class        string
	Sex          string
	XP           int
	Strength     int
	Constitution int
	Dexterity    int
	Intelligence int
	Wisdom       int
	Charisma     int
	Level        int
	Hp           int
	PXloc        int
	PYloc        int
	Objects      [maxobjscarried]int
}

type Player interface {
	init(id, xloc, yloc int, name, class, sex string)
	move(dir string) (err bool)
	pickup(objid int) (err bool)
	drop(objid int) (err bool)
	show(output *string)
	attack(object world.ObjectType, player *Character) (err bool, npcdmg, playerdmg int)
}

const (
	maxstat        int = 18
	maxobjscarried int = 10
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (ch *Character) init(id, xloc, yloc int, name, class, sex string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ch.PlayerID = id
	ch.Name = name
	ch.Class = class
	ch.Sex = sex
	ch.XP = 5
	ch.Strength = r.Intn(maxstat)
	ch.Constitution = r.Intn(maxstat)
	ch.Dexterity = r.Intn(maxstat)
	ch.Intelligence = r.Intn(maxstat)
	ch.Wisdom = r.Intn(maxstat)
	ch.Charisma = r.Intn(maxstat)
	ch.Hp = 20 + r.Intn(maxstat)
	ch.Level = 1
	ch.PXloc = xloc
	ch.PYloc = yloc
	ch.Objects[0] = 0
}

func (ch *Character) move(dir string) (err bool) {
	switch dir {
	// North
	case "N":
		if ch.PYloc < world.Maxworldy-1 {
			ch.PYloc++
			return true
		}
	// East
	case "E":
		if ch.PXloc < world.Maxworldx-1 {
			ch.PXloc++
			return true
		}
	// South
	case "S":
		if ch.PYloc > 0 {
			ch.PYloc--
			return true
		}
	// West
	case "W":
		if ch.PXloc > 0 {
			ch.PXloc--
			return true
		}
	default:
		return false
	}
	return false
}

func (ch *Character) pickup(objid int) (err bool) {
	for i := range ch.Objects {
		if ch.Objects[i] == 0 {
			ch.Objects[i] = objid
			return true
		} else {
			i++
		}
	}
	return false
}

func (ch *Character) drop(objid int) (err bool) {
	for i := range ch.Objects {
		if ch.Objects[i] == objid {
			ch.Objects[i] = 0
			return true
		} else {
			i++
		}
	}
	return false
}

func (ch *Character) Save(savefname string) (err bool) {
	file, error := os.Create(savefname)
	check(error)
	defer file.Close()
	if error == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(ch)
		return true
	}
	return false

}

func (ch *Character) Load(loadfname string) (err bool) {
	// Set location and objects to Zero before decode
	ch.PXloc, ch.PYloc = 0, 0
	for i := range ch.Objects {
		ch.Objects[i] = 0
	}
	file, error := os.Open(loadfname)
	defer file.Close()
	if error == nil {
		decoder := gob.NewDecoder(file)
		decoder.Decode(ch)
		return true
	}
	return false
}

func (ch *Character) attack(object world.ObjectType, player *Character) (err bool, npcdmg, playerdmg int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if object.ObjID != 0 {
		playerdmg = r.Intn(ch.Strength) / 4
		npcdmg = r.Intn(object.Attack)
		ch.Hp -= npcdmg
		player.Hp -= playerdmg
		return true, npcdmg, playerdmg
	}
	return
}

func (ch *Character) show(output *string) {
	*output = fmt.Sprintln("<table class=\"pstats\">")
	*output += fmt.Sprintf("<tr><td>Name</td><td>%s</td></tr>\n", ch.Name)
	*output += fmt.Sprintf("<tr><td>Class</td><td>%s</td></tr>\n", ch.Class)
	*output += fmt.Sprintf("<tr><td>Sex</td><td>%s</td></tr>\n", ch.Sex)
	*output += fmt.Sprintf("<tr><td>Level</td><td>%d</td></tr>\n", ch.Level)
	*output += fmt.Sprintf("<tr><td>XP</td><td>%d</td></tr>\n", ch.XP)
	*output += fmt.Sprintf("<tr><td>HP</td><td>%d</td></tr>\n", ch.Hp)
	*output += fmt.Sprintf("<tr><td>STR</td><td>%d</td></tr>\n", ch.Strength)
	*output += fmt.Sprintf("<tr><td>CON</td><td>%d</td></tr>\n", ch.Constitution)
	*output += fmt.Sprintf("<tr><td>DEX</td><td>%d</td></tr>\n", ch.Dexterity)
	*output += fmt.Sprintf("<tr><td>INT</td><td>%d</td></tr>\n", ch.Intelligence)
	*output += fmt.Sprintf("<tr><td>WIS</td><td>%d</td></tr>\n", ch.Wisdom)
	*output += fmt.Sprintf("<tr><td>CHR</td><td>%d</td></tr>\n", ch.Charisma)
	*output += fmt.Sprintln("</table>")
}

func Showplayer(p Player, output *string) {
	p.show(output)
}

func Createplayer(p Player, id, xloc, yloc int, name, class, sex string) {
	p.init(id, xloc, yloc, name, class, sex)
}

func Moveplayer(p Player, dir string) (err bool) {
	return p.move(dir)
}

func Playerpickup(p Player, objid int) (err bool) {
	return p.pickup(objid)
}

func Playerdrop(p Player, objid int) (err bool) {
	return p.drop(objid)
}

func NpcAttack(npc Player, object world.ObjectType, player *Character) (err bool, npcdmg, playerdmg int) {
	return npc.attack(object, player)
}
