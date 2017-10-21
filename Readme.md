This is a plugin for [`YAGS`](https://github.com/pltanton/yags).

It displays battery level.

# Configuration
Available configuration keys:

| Key           | Description                                                                         | Default value             |
| ---           | ---                                                                                 | ---                       |
| name          | UPower name of battery                                                              | BAT0                      |
| format        | message format for dischargin state                                                 | {icon} {lvl}              |
| acFormat      | message format for charging state                                                   | {icon} {lvl}              |
| full          | level of fully charged battery                                                      | 100                       |
| animationTick | duration of animation update interval for ac state in `time.ParseDuration` notation | 1000ms                    |
| icons         | set of icons, used for animation and indication                                     | ["", "", "", "", ""] |

Available variables for interpolation:

| Variable name | Description               |
| ---           | ---                       |
| lvl           | battery level in precents |
| icon          | icon from `icons` set     |
