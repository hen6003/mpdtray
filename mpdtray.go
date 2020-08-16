package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dawidd6/go-appindicator"
	"github.com/fhs/gompd/mpd"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var status mpd.Attrs = nil
var songNameItem *gtk.MenuItem = nil
var pausedItem *gtk.CheckMenuItem = nil
var randItem *gtk.CheckMenuItem = nil

func update(mpdConn *mpd.Client) bool {
	var err error = nil
	status, err = mpdConn.Status()
	if err != nil {
		log.Fatalln(err)
	}

	song, err := mpdConn.CurrentSong()
	if err != nil {
		log.Fatalln(err)
	}

	songName := song["Title"]
	if songName == "" {
		songName = song["file"]
		songNameSplit := strings.Split(songName, "/")
		songName = songNameSplit[len(songNameSplit)-1]
		songNameSplit = strings.Split(songName, ".")
		songName = songNameSplit[0]
	}

	runes := []rune(songName)
	if len(runes) > 20 {
		songName = string(runes[:20])
	}

	songNameItem.SetLabel(songName)

	if status["state"] == "play" {
		pausedItem.SetStateFlags(gtk.STATE_FLAG_INCONSISTENT, true)
	} else {
		pausedItem.SetStateFlags(gtk.STATE_FLAG_CHECKED, true)
	}

	return true
}

func indicator(port string) {
	conn, err := mpd.Dial("tcp", port)
	if err != nil {
		fmt.Println("MPD port wrong, Usage: mpdtray <mpdPort>, default: 'localhost:6600'")
		log.Fatalln(err)
	}

	gtk.Init(nil)

	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}

	songNameItem, err = gtk.MenuItemNew()
	if err != nil {
		log.Fatal(err)
	}

	nextItem, err := gtk.MenuItemNewWithLabel("Next")
	if err != nil {
		log.Fatal(err)
	}

	pausedItem, err = gtk.CheckMenuItemNewWithLabel("Paused")
	if err != nil {
		log.Fatal(err)
	}

	seperatorItem, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		log.Fatal(err)
	}

	quitItem, err := gtk.MenuItemNewWithLabel("Exit")
	if err != nil {
		log.Fatal(err)
	}

	glib.TimeoutAdd(5000, update, conn)

	update(conn)
	indicator := appindicator.New("mpd_tray", "mpd", appindicator.CategoryApplicationStatus)
	indicator.SetTitle("mpd_tray")
	indicator.SetLabel("MPD", "")
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)
	songNameItem.SetSensitive(false)

	_, err = pausedItem.Connect("activate", func() {
		if status["state"] == "play" {
			conn.Pause(true)
		} else {
			conn.Pause(false)
		}

		update(conn)
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = nextItem.Connect("activate", func() {
		conn.Next()

		update(conn)
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = quitItem.Connect("activate", func() { gtk.MainQuit() })
	if err != nil {
		log.Fatal(err)
	}

	menu.Add(songNameItem)
	menu.Add(nextItem)
	menu.Add(pausedItem)
	menu.Add(seperatorItem)
	menu.Add(quitItem)
	menu.ShowAll()

	gtk.Main()
}

func main() {

	args := os.Args

	mpdPort := "localhost:6600"
	if len(args) != 1 {
		if strings.Contains(args[1], "help") {
			fmt.Println("Usage:\n\nmpdtray help: print this help message\nmpdtray <port>: choose what port mpd is on, default localhost:6600\n\nif no ip is defined it defaults to localhost.\nexample: mpdtray 4.4.4.4:6666 would be 4.4.4.4:6666\nexample: mpdtray 6666 would be localhost:6666")
			return
		}

		if strings.Contains(args[1], ":") {
			mpdPort = args[1]
		} else {
			mpdPort = "localhost:" + args[1]
		}
	}

	indicator(mpdPort)
}
