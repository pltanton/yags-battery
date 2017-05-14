This is a plugin for [`YAGS`](https://github.com/pltanton/yags).

It displays battery level.

# Configuration
Available configuration keys:

| Key    | Description                 | Default value |
| ---    | ---                         | ---           |
| name   | UPower name of battery      | BAT0          |
| high   | format for _lvl_ > 75       | {lvl}         |
| medium | format for _lvl_ > 35       | {lvl}         |
| low    | format for _lvl_ > 12       | {lvl}         |
| empty  | format for _lvl_ <= 12      | {lvl}         |
| ac     | format for conneted adapted | {lvl}         |

Available variables for interpolation:

| Variable name | Description               |
| ---           | ---                       |
| lvl           | battery level in precents |
