#! /usr/bin/python3
import trayicon, os, argparse
import tkinter as tk
from mpd import MPDClient

icon_path = os.path.join(os.path.dirname(__file__), 'icon.png')
global client
client = MPDClient()

def make_icon(gui):
    def update():
        client.disconnect()
        client.connect("localhost", port)

        if client.status()['random'] == '1':
            icon.menu.set_item_value(2, True)
        else:
            icon.menu.set_item_value(2, False)
        
        if client.status()['state'] == 'pause':
            icon.menu.set_item_value(1, True)
        else:
            icon.menu.set_item_value(1, False)

    def random():
        client.disconnect()
        client.connect("localhost", port)

        if icon.menu.get_item_value(2):
            client.random(1)
        else:
            client.random(0)

    def pause():
        client.disconnect()
        client.connect("localhost", port)

        if icon.menu.get_item_value(1):
            client.pause(1)
        else:
            client.pause(0)

    def next():
        client.next()
        update()

    if gui == 'qt':
        module = trayicon.qticon
    elif gui == 'tk':
        module = trayicon.tkicon
    else:
        module = trayicon.gtkicon

    icon = module.TrayIcon('icon.png', fallback_icon_path=icon_path)

    icon.menu.add_command(label = 'Next', command = next)
    #icon.menu.add_checkbutton(command = pause)

    icon.menu.add_checkbutton(label = 'Pause', command=pause)

    # checkbutton
    icon.menu.add_checkbutton(label = 'Random', command=random)

    #icon.menu.disable_item(0)
    # separator
    icon.menu.add_separator()

    icon.menu.add_command(label='Quit', command=root.destroy)

    icon.bind_left_click(command = update)

    # start icon event loop
    update()
    icon.loop(root)


toolkits = trayicon.get_available_gui_toolkits()

root = tk.Tk()
            # network timeout in seconds (floats allowed), default: None
client.idletimeout = None          # timeout for fetching the result of the idle command is handled seperately, default: None
 
parser = argparse.ArgumentParser(description = 'A system tray widget for controlling mpd')

parser.add_argument("-s", "--style", type = str, help = 'What style you want (gtk/qt/tk), defaults to gtk')
parser.add_argument("-p", "--port", type = int, help = 'What port mpd is on, defaults to 6600')

args = parser.parse_args()

global port
port = args.port


if port == None:
    port = 6600

#client.connect("localhost", port)

root.withdraw()

try:
    style = args.style
except IndexError:
    style = 'gtk'

make_icon(style)

root.mainloop()