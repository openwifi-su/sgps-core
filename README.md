# sgps-core

Is a restful API backend for WIFI based location service.


## configfile

The config file is a .toml file.

following config parameters have to been defined:

```
# A sample TOML config file.
# For Postgresql
[database]
psql_user = "<USER>"
psql_password = "<PASSWORD>"
psql_name = "<DATABASE NAME>"
psql_tablename = "<TABLENAME>"
# For MySQL
[database]
msql_user = "<USER>"
msql_password = "<PASSWORD>"
msql_name = "<DATABASE NAME>"
msql_tablename = "<TABLENAME>"
# If bouth Postgresql and MySQL defined sgps will prefrare Postgresql

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
