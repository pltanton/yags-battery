package main

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("icons", []string{"", "", "", "", ""})
	v.SetDefault("full", 100)
	v.SetDefault("animationTick", "1000ms")
	v.SetDefault("format", "{icon} {lvl}")
	v.SetDefault("acFormat", "{icon} {lvl}")
	v.SetDefault("name", "BAT0")
	return v
}
