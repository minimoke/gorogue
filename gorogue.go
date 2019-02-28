package main

import (
	"fmt"
	"gorogue/player"
	"gorogue/world"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const pageheader = `<html>
<head>
<title>GO Rogue</title>
<link rel="stylesheet" href="static/styles.css">
</head>
<body>
<table class="main">
<tr>
`

const (
	aidist = 2
)

type WebContent struct {
	Gamemap     string
	Textarea    string
	PlayerStats string
}

type SaveLoader interface {
	Save(fname string) (err bool)
	Load(fname string) (err bool)
}

func processplayerinput(playerinput string) {

	switch playerinput {
	case "Q":
		return
	case "A 2", "A 3", "A 4":
		pslice := strings.Fields(playerinput)
		npcnum, _ := strconv.Atoi(pslice[1])
		pidx := npcnum - 2
		pid := npc[pidx].PlayerID
		pname := npc[pidx].Name
		objname := objmap[player1.Objects[0]].Name
		for i := range world1[player1.PXloc][player1.PYloc].Players {
			if world1[player1.PXloc][player1.PYloc].Players[i] == npcnum {
				_, npcdmg, pdmg := player.NpcAttack(&npc[pidx], objmap[player1.Objects[0]], &player1)
				webcontent.Textarea = fmt.Sprintf("Player %d-%s attacked with %s for %d damage taking %d damage\n", pid, pname, objname, npcdmg, pdmg)
				break
			} else {
				webcontent.Textarea = fmt.Sprintln("Player not at location")
			}
		}

	case "L":
		webcontent.Textarea = fmt.Sprintf("%s\n", world1[player1.PXloc][player1.PYloc].Desc)
		objects := world1[player1.PXloc][player1.PYloc].Objects
		for i := range objects {
			obji := world1[player1.PXloc][player1.PYloc].Objects[i]
			if obji != 0 {
				webcontent.Textarea += fmt.Sprintf("Object %d-%s\n", i+1, worldobjects1[obji].ObjType)
			}
		}
		players := world1[player1.PXloc][player1.PYloc].Players
		for i := range players {
			playeri := world1[player1.PXloc][player1.PYloc].Players[i]
			if playeri != 0 && playeri != 1 {
				pid := npc[playeri-2].PlayerID
				pname := npc[playeri-2].Name
				webcontent.Textarea += fmt.Sprintf("Player %d-%s\n", pid, pname)
			}
		}
	case "?":
		webcontent.Textarea = fmt.Sprintln("Please enter a direction to move (n,s,e,w) or l to look?")
	case "N", "S", "E", "W":
		updateplayerloc(&player1, playerinput)
		chsem <- 1
	case "P 1", "P 2", "P 3", "P 4", "P 5":
		pslice := strings.Fields(playerinput)
		objnum, _ := strconv.Atoi(pslice[1])
		obji := world1[player1.PXloc][player1.PYloc].Objects[objnum-1]
		if obji != 0 {
			if ok := player.Playerpickup(&player1, worldobjects1[obji].ObjID); ok {
				webcontent.Textarea = fmt.Sprintf("%d-%s picked up\n", worldobjects1[obji].ObjID, worldobjects1[obji].ObjType)
				world.Removeobject(&world1, player1.PXloc, player1.PYloc, worldobjects1[obji].ObjID)
			} else {
				webcontent.Textarea = fmt.Sprintln("You can't carry anymore items")
			}
		}
	case "I":
		empty := true
		objects := player1.Objects
		webcontent.Textarea += fmt.Sprintln("\nYou are carrying:")
		for i := range objects {
			obji := objects[i]
			if obji != 0 {
				webcontent.Textarea += fmt.Sprintf("Slot %d-%s\n", i+1, worldobjects1[obji].ObjType)
				empty = false
			}
		}
		if empty {
			webcontent.Textarea += fmt.Sprintln("Nothing")
		}
	case "D 1", "D 2", "D 3", "D 4", "D 5":
		pslice := strings.Fields(playerinput)
		objnum, _ := strconv.Atoi(pslice[1])
		obji := player1.Objects[objnum-1]
		if obji != 0 {
			if ok := player.Playerdrop(&player1, worldobjects1[obji].ObjID); ok {
				webcontent.Textarea = fmt.Sprintf("%d-%s dropped\n", worldobjects1[obji].ObjID, worldobjects1[obji].ObjType)
				world.Putobject(&world1, player1.PXloc, player1.PYloc, worldobjects1[obji].ObjID)
			} else {
				webcontent.Textarea = fmt.Sprintln("You can't drop anymore items at this location")
			}
		}
	case "SAVE":
		gamesave(&player1, &world1)
		webcontent.Textarea = fmt.Sprintln("Game saved")
	case "LOAD":
		if ok := gameload(&player1, &world1); ok {
			webcontent.Textarea = fmt.Sprintln("Game loaded")
		} else {
			webcontent.Textarea = fmt.Sprintln("Error - Game not loaded")
		}
	default:
		webcontent.Textarea = fmt.Sprintln("Invalid command entered")
	}
}

func updateplayerloc(ch *player.Character, direction string) (err bool) {
	oldPXloc := ch.PXloc
	oldPYloc := ch.PYloc
	if ok := player.Moveplayer(ch, direction); ok {
		world.Removeplayer(&world1, oldPXloc, oldPYloc, ch.PlayerID)
		world.Putplayer(&world1, ch.PXloc, ch.PYloc, ch.PlayerID)
		if ch.PlayerID == 1 {
			webcontent.Textarea = fmt.Sprintf("Old Player location %d,%d\n", oldPXloc, oldPYloc)
			webcontent.Textarea += fmt.Sprintf("Player location %d,%d\n", ch.PXloc, ch.PYloc)
			webcontent.Textarea += fmt.Sprintf("%s\n", world1[ch.PXloc][ch.PYloc].Desc)
		}
		return true
	} else {
		if ch.PlayerID == 1 {
			webcontent.Textarea = fmt.Sprintln("You can't move in that direction")
		} else {
			return false
		}
	}
	return false
}

func mainairoutine(npc *player.Character, npcsem chan int) {
	for {
		<-chsem
		dir := calcaimove(npc)
		switch dir {
		case 0:
			_ = updateplayerloc(npc, "N")
		case 1:
			_ = updateplayerloc(npc, "E")
		case 2:
			_ = updateplayerloc(npc, "S")
		case 3:
			_ = updateplayerloc(npc, "W")
		}
	}
}

func calcaimove(npc *player.Character) (dir int) {
	xdiff := npc.PXloc - player1.PXloc
	ydiff := npc.PYloc - player1.PYloc
	switch {
	case xdiff > aidist:
		dir = 3
	case xdiff < (-aidist):
		dir = 1
	case ydiff > aidist:
		dir = 2
	case ydiff < (-aidist):
		dir = 0
	default:
		dir = -1
	}
	//fmt.Printf("NPC %d: [%d,%d,%d]\n", npc.PlayerID, xdiff, ydiff, dir)
	return dir
}

func gamesave(p, w SaveLoader) {
	_ = w.Save("worldsave1.bin")
	_ = p.Save("playersave1.bin")
}

func gameload(p, w SaveLoader) (err bool) {
	if ok := w.Load("worldsave1.bin"); !ok {
		return false
	}

	if ok := p.Load("playersave1.bin"); !ok {
		return false
	}
	return true

}

func getpageinput(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		webcontent.Textarea = fmt.Sprintf("Welcome %s the %s, you stand on a grassy plain with mountains all around you.\n", player1.Name, player1.Class)
		renderpage(&w)
	} else {
		r.ParseForm()
		forminput := r.Form["Command"][0]
		playerinput := strings.ToUpper(forminput)
		if playerinput != "Q" {
			processplayerinput(playerinput)
		} else {
			os.Exit(0)
		}
		renderpage(&w)
	}
}

func renderpage(w *http.ResponseWriter) {
	fmt.Fprintf(*w, "%s", pageheader)
	fmt.Fprintf(*w, "<td align=\"center\">\n")
	world.Showworld(&world1, &webcontent.Gamemap)
	fmt.Fprintf(*w, "%s", webcontent.Gamemap)
	fmt.Fprintf(*w, "</td>\n")
	fmt.Fprintf(*w, "<td align=\"center\"><table><tr><td colspan=\"2\"><img src=\"static/fighterf1.jpeg\"></td>\n")
	fmt.Fprintf(*w, "<tr><td>%s</td><td>%s</td></tr></table>", objmap[player1.Objects[0]].Name, objmap[player1.Objects[1]].Name)
	fmt.Fprintf(*w, "<td align=\"center\">")
	player.Showplayer(&player1, &webcontent.PlayerStats)
	fmt.Fprintf(*w, "%s", webcontent.PlayerStats)
	fmt.Fprintf(*w, "</td>\n")
	fmt.Fprintf(*w, "<tr>")
	fmt.Fprintf(*w, "<td colspan=\"3\">")
	pagetempl.Execute(*w, webcontent)
}

var world1 world.Worldmap
var player1 player.Character
var npc [3]player.Character
var worldobjects1 world.Worldobjects
var webcontent WebContent
var chsem = make(chan int)
var objmap = make(map[int]world.ObjectType)
var pagetempl *template.Template

func main() {

	// populate world
	world.Createworld(&world1, "data/testmap1.csv")
	worldobjects1.Loadobjectlist("data/objectlist1.csv", objmap)

	// create player
	player.Createplayer(&player1, 1, 1, 2, "Tordek", "Fighter", "Female")
	world.Putplayer(&world1, player1.PXloc, player1.PYloc, player1.PlayerID)
	fmt.Println(player1)

	// create NPC 1
	player.Createplayer(&npc[0], 2, 9, 9, "Jack", "Thief", "Male")
	world.Putplayer(&world1, npc[0].PXloc, npc[0].PYloc, npc[0].PlayerID)
	fmt.Println(npc[0])
	go mainairoutine(&npc[0], chsem)

	// create NPC 2
	player.Createplayer(&npc[1], 3, 0, 9, "Cassandra", "Mage", "Female")
	world.Putplayer(&world1, npc[1].PXloc, npc[1].PYloc, npc[1].PlayerID)
	fmt.Println(npc[1])
	go mainairoutine(&npc[1], chsem)

	// create NPC 3
	player.Createplayer(&npc[2], 4, 9, 0, "Larry", "Ranger", "Male")
	world.Putplayer(&world1, npc[2].PXloc, npc[2].PYloc, npc[2].PlayerID)
	fmt.Println(npc[2])
	go mainairoutine(&npc[2], chsem)

	// Initialise page templage
	pagetempl = template.Must(template.ParseFiles("static/page1.gtpl"))

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/gorogue", getpageinput)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:	", err)
	}

}
