# boxter

boxter uses a ini configuration file to determine which versions of the voxbox playbook should be provisioned for specific hosts.

This is a sample of the ini file

```ini
[latest]
box2 =
box3  =
box4 =

[0.1.0]
habarimaalum=

[0.2.1]
mambojambo=
```

The sections `[all], [0.1.0] etc ` defines the versions of the playbook. The keys defined under these sections represent the host names.

__IMPORTANT__ don't forget to add the `=` sign after the host name. The configuration is nini we need to add this to be compliant with the parser but it has no effect whatsoeve.

So that sample says provision voxbox version `0.1.0` for a host named `habarimaalum`


## Commands

sync

```
NAME:
   boxter sync - syncs new or specified version of playbook

USAGE:
   boxter sync [command options] [arguments...]

OPTIONS:
   --config value  path to the configuration file
   --host value    the name of the host machine
   --ssh value     ssh connection url
```