package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/godbus/dbus"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/plugins"
	"github.com/pltanton/yags/utils"
)

type battery struct {
	conf           *viper.Viper
	batName        string
	out            chan string
	acTickDuration time.Duration
	acTimer        *time.Timer
}

// New returns new instance of battery plugin by given name
func New(conf *viper.Viper) plugins.Plugin {
	conf = setDefaults(conf)
	tickDuration, err := time.ParseDuration(conf.GetString("animationTick"))
	if err != nil {
		log.Panic("Can't parse tick duration", err)
	}
	return battery{
		out:            make(chan string, 1),
		conf:           conf,
		acTickDuration: tickDuration,
		acTimer:        time.NewTimer(tickDuration),
	}
}

// Chan returns a strings channel with battery state monitor
func (b battery) Chan() chan string {
	return b.out
}

// StartMonitor starts monitoring for battery changing events
func (b battery) StartMonitor() {
	lvl, state := b.parseBatLevel()
	b.out <- b.formatMessage(lvl, state)
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(fmt.Errorf("cannot connect dbus session: %s", err.Error()))
	}

	arg := fmt.Sprintf(
		"type='signal',path='%s',interface='%s',member='%s',sender='%s'",
		fmt.Sprintf("/org/freedesktop/UPower/devices/battery_%s", b.conf.GetString("name")),
		"org.freedesktop.DBus.Properties",
		"PropertiesChanged",
		"org.freedesktop.UPower",
	)

	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		arg,
	)

	// Init channel to listen system resume
	resumeChan := make(chan *dbus.Signal, 1)
	utils.GetResumeDbusConn().Signal(resumeChan)

	c := make(chan *dbus.Signal, 1)
	conn.Signal(c)
	for {
		select {
		case <-c:
		case <-resumeChan:
		case <-b.acTimer.C:
		}

		lvl, state := b.parseBatLevel()
		if state == 1 || state == 4 {
			b.acTimer.Reset(b.acTickDuration)
		}

		b.out <- b.formatMessage(lvl, state)
	}
}

// formatMessage formats message for printing
func (b battery) formatMessage(lvl int, state uint32) string {
	var format, icon string
	if state != 2 {
		format = b.conf.GetString("acFormat")
		icon = b.getAnimationIcon()
	} else {
		format = b.conf.GetString("format")
		icon = b.getIcon(lvl)
	}
	format = utils.ReplaceVar(format, "icon", icon)
	return utils.ReplaceVar(format, "lvl", strconv.Itoa(lvl))
}

func (b battery) getAnimationIcon() string {
	iconSet := b.conf.Get("icons").([]interface{})
	tickDuration, _ := time.ParseDuration(b.conf.GetString("animationTick"))
	return iconSet[(time.Now().UnixNano()/tickDuration.Nanoseconds())%int64(len(iconSet))].(string)
}

func (b battery) getIcon(lvl int) string {
	iconSet := b.conf.Get("icons").([]interface{})
	full := b.conf.GetInt("full")
	if lvl >= full {
		return iconSet[len(iconSet)-1].(string)
	}

	var delta float32
	delta = (float32(full) / float32(len(iconSet)))
	for i := 1; i <= len(iconSet); i++ {
		if float32(lvl) < delta*float32(i) {
			return iconSet[i-1].(string)
		}
	}
	return ""
}

// parseBatLevel connects to the system bus and get the State and Percentage
// properties from the UPower's BAT object. It returns the level in percents
// and integer status, which means:
//
//  0: Unknown
//  1: Charging
//  2: Discharging
//  3: Empty
//  4: Fully charged
//  5: Pending charge
//  6: Pending discharge
//
func (b battery) parseBatLevel() (int, uint32) {
	conn, _ := dbus.SystemBus()
	pth := fmt.Sprintf("/org/freedesktop/UPower/devices/battery_%s", b.conf.GetString("name"))
	object := conn.Object(
		"org.freedesktop.UPower",
		dbus.ObjectPath(pth),
	)
	lvl, _ := object.GetProperty("org.freedesktop.UPower.Device.Percentage")
	state, _ := object.GetProperty("org.freedesktop.UPower.Device.State")
	return int(lvl.Value().(float64)), state.Value().(uint32)
}
