# boxter

boxter is a commandline application which helps track and sync multiple voxbox playbooks between machines.

## Dependencies

- `rsync` boxter relies on rsync to update playbooks

## Commands

sync

```
NAME:
   boxter sync - syncs new or specified version of playbook

USAGE:
   boxter sync [command options] [arguments...]

OPTIONS:
   --config value               path to the configuration file
   --host value                 the name of the host machine
   --remote-playbook-dir value  the directory to sync playbooks to in a remote host
   --ssh value                  ssh connection url
```

The `--config` points to a simple configuration file which has the following formal

```json
{
	"hostsFile": "path/to/hosts/ini/file",
	"localPlaybookDIr": "/path/to the/local/playbooks",
	"remotePlaybookDIr": "/path/to/remote/playbooks"
}
```

You can also specify the version of the playbook you want to sync as the first argument to the sync command

```
boxter sync --config /path/to/config  0.1.0
```

This way, `0.1.0` will take precedence over  the version you specified on the ini file.

boxter uses a ini configuration file to determine which versions of the voxbox playbook should be provisioned for specific hosts.
In the sample config above `hostsFile` is the path to this configuration file.

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
