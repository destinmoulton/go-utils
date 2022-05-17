package lib

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

type DBusUtils struct{}

var DBus DBusUtils

func (b *DBusUtils) Msg(title string, contents string) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Errorf("cannot connect to dbus: %v", err)
	}
	defer conn.Close()

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
		"", title, contents, []string{},
		map[string]dbus.Variant{}, int32(5000))
	if call.Err != nil {
		log.Errorf("unable to send dbus message: %v", call.Err)
	}
}
