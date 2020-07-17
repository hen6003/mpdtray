package main

import (
	"log"
	"strings"

	"github.com/dawidd6/go-appindicator"
	"github.com/fhs/gompd/mpd"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var status mpd.Attrs = nil
var songNameItem *gtk.MenuItem = nil
var pausedItem *gtk.CheckMenuItem = nil

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

	songName := song["file"]
	songNameSplit := strings.Split(songName, "/")
	songName = songNameSplit[len(songNameSplit)-1]
	songNameSplit = strings.Split(songName, ".")
	songName = songNameSplit[0]

	runes := []rune(songName)
	if len(runes) > 20 {
		songName = string(runes[:20])
	}

	songNameItem.SetLabel(songName)

	if status["state"] == "play" {
		pausedItem.SetStateFlags(gtk.STATE_FLAG_NORMAL, true)
	} else {
		pausedItem.SetStateFlags(gtk.STATE_FLAG_CHECKED, true)
	}

	return true
}

func main() {
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

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

	update(conn)
	glib.TimeoutAdd(5000, update, conn)

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
