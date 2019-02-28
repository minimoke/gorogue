package world

import (
	"bufio"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ObjectType struct {
	ObjID   int
	ObjType string
	Name    string
	Desc    string
	Attack  int
	Defence int
	Speed   int
	Weight  int
}

type Location struct {
	Xloc    int
	Yloc    int
	Desc    string
	Objects [Maxobjects]int
	Players [Maxplayers]int
}

type Worldmap [Maxworldx][Maxworldy]Location

type Worldloader interface {
	populate(csvfname string)
	show(output *string)
	removeobject(xloc, yloc, objid int)
	putobject(xloc, yloc, objid int) (err bool)
	removeplayer(xloc, yloc, playerid int)
	putplayer(xloc, yloc, playerid int) (err bool)
}

type Worldobjects [Maxobjects]ObjectType

const (
	Maxworldx  int = 10
	Maxworldy  int = 10
	Maxobjects int = 5
	Maxplayers int = 5
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (myworld *Worldmap) populate(csvfname string) {
	csvfile, err := os.Open(csvfname)
	check(err)
	defer csvfile.Close()

	scanner := bufio.NewScanner(csvfile)

	// Read header line of the file
	scanner.Scan()
	r := csv.NewReader(strings.NewReader(scanner.Text()))
	_, err = r.Read()
	check(err)

	// Read the rest of the file
	for x := 0; x < Maxworldx; x++ {
		for y := 0; y < Maxworldy; y++ {
			scanner.Scan()
			r := csv.NewReader(strings.NewReader(scanner.Text()))
			record, err := r.Read()
			if err != io.EOF {
				myworld[x][y].Xloc, _ = strconv.Atoi(record[0])
				myworld[x][y].Yloc, _ = strconv.Atoi(record[1])
				myworld[x][y].Desc = record[2]
				for i := 0; i < Maxobjects; i++ {
					objid, _ := strconv.Atoi(record[3+i])
					myworld[x][y].Objects[i] = objid
				}
				for i := 0; i < Maxplayers; i++ {
					myworld[x][y].Players[i] = 0
				}
			}
		}
	}
}

func (myworld *Worldmap) removeobject(xloc, yloc, objid int) {
	for i := 0; i < Maxobjects; i++ {
		if myworld[xloc][yloc].Objects[i] == objid {
			myworld[xloc][yloc].Objects[i] = 0
		}
	}
}

func (myworld *Worldmap) putobject(xloc, yloc, objid int) (err bool) {
	for i := range myworld[xloc][yloc].Objects {
		if myworld[xloc][yloc].Objects[i] == 0 {
			myworld[xloc][yloc].Objects[i] = objid
			return true
		} else {
			i++
		}
	}
	return false
}

func (myworld *Worldmap) removeplayer(xloc, yloc, playerid int) {
	for i := 0; i < Maxplayers; i++ {
		if myworld[xloc][yloc].Players[i] == playerid {
			myworld[xloc][yloc].Players[i] = 0
		}
	}
}

func (myworld *Worldmap) putplayer(xloc, yloc, playerid int) (err bool) {
	for i := range myworld[xloc][yloc].Players {
		if myworld[xloc][yloc].Players[i] == 0 {
			myworld[xloc][yloc].Players[i] = playerid
			return true
		} else {
			i++
		}
	}
	return false
}

func (myworld *Worldmap) show(output *string) {
	var playerfound bool
	var playerid int
	*output = fmt.Sprintln("<table class=\"map\">")
	for y := Maxworldy - 1; y >= 0; y-- {
		*output += fmt.Sprintf("<tr>")
		playerfound = false
		for x := 0; x < Maxworldx; x++ {
			playerfound = false
			playerid = 0
			for i := range myworld[x][y].Players {
				if myworld[x][y].Players[i] != 0 {
					playerfound = true
					playerid += myworld[x][y].Players[i]
				}
			}
			if playerfound == true {
				*output += fmt.Sprintf("<td align=\"center\" style=\"background-color:gray; color:white\">%d</td>", playerid)
			} else {
				*output += fmt.Sprintf("<td align=\"center\">&nbsp;</td>")
			}
		}
		*output += fmt.Sprintf("</tr>\n")
		playerfound = false
	}
	*output += fmt.Sprintln("</table>")
}

func (w *Worldmap) Save(savefname string) (err bool) {
	file, error := os.Create(savefname)
	check(error)
	defer file.Close()
	if error == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(w)
		return true
	}
	return false

}

func (w *Worldmap) Load(loadfname string) (err bool) {
	// Set to Zero before Decode from File
	for y := Maxworldy - 1; y >= 0; y-- {
		for x := 0; x < Maxworldx; x++ {
			for i := 0; i < Maxobjects; i++ {
				w[x][y].Objects[i] = 0
			}
			for i := 0; i < Maxplayers; i++ {
				w[x][y].Players[i] = 0
			}
		}
	}
	file, error := os.Open(loadfname)
	defer file.Close()
	if error == nil {
		decoder := gob.NewDecoder(file)
		decoder.Decode(w)
		return true
	}
	return false
}

func (myobjectlist *Worldobjects) Loadobjectlist(csvfname string, objmap map[int]ObjectType) {

	csvfile, err := os.Open(csvfname)
	check(err)
	defer csvfile.Close()

	scanner := bufio.NewScanner(csvfile)

	// Read header line of the file
	scanner.Scan()
	r := csv.NewReader(strings.NewReader(scanner.Text()))
	_, err = r.Read()
	check(err)

	// Read the rest of the file
	for i := 0; i < Maxobjects; i++ {
		scanner.Scan()
		r := csv.NewReader(strings.NewReader(scanner.Text()))
		record, err := r.Read()
		if err != io.EOF {
			myobjectlist[i].ObjID, _ = strconv.Atoi(record[0])
			myobjectlist[i].ObjType = record[1]
			myobjectlist[i].Name = record[2]
			myobjectlist[i].Desc = record[3]
			myobjectlist[i].Attack, _ = strconv.Atoi(record[4])
			myobjectlist[i].Defence, _ = strconv.Atoi(record[5])
			myobjectlist[i].Speed, _ = strconv.Atoi(record[6])
			myobjectlist[i].Weight, _ = strconv.Atoi(record[7])
			objmap[myobjectlist[i].ObjID] = myobjectlist[i]
		}
	}
}

func Createworld(w Worldloader, worldcsvfname string) {
	w.populate(worldcsvfname)
}

func Showworld(w Worldloader, output *string) {
	w.show(output)
}

func Putobject(w Worldloader, xloc, yloc, objid int) {
	w.putobject(xloc, yloc, objid)
}

func Removeobject(w Worldloader, xloc, yloc, objid int) {
	w.removeobject(xloc, yloc, objid)
}

func Removeplayer(w Worldloader, xloc, yloc, playerid int) {
	w.removeplayer(xloc, yloc, playerid)
}

func Putplayer(w Worldloader, xloc, yloc, playerid int) {
	w.putplayer(xloc, yloc, playerid)
}
