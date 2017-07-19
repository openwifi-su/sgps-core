# sgps-core

Is a restful API backend for WIFI based location service.


## configfile

The config file is a .toml file.

following config parameters have to been defined:

```
# A sample TOML config file.
[database]
db_user = "<USER>"
db_password = "<PASSWORD>"
db_name = "<DATABASE NAME>"

[MLS]
apikey = "<API KEY>"

[old_api]
path = "<PART TO API>"
port = <PORT>

[new_api]
path = "<PART TO API>"
```

## Information for developers

The main development process takes place on:
[git.ffnw.de](https://git.ffnw.de/GSoC/sgps-core "ffnw gitlab repo")

Discussions can be taken on the Mailinglist [openwifi@lists.ffnw.de](https://lists.ffnw.de/mailman/listinfo/openwifi "Mailinglist for openwifi.su")

You can create and send Patches with `git send-email` as well to the above mentioned ML.
