# migration-passwd

This is a Golang CLI mostly based on the
[migrationtools](https://gitlab.com/future-ad-laboratory/migrationtools) code,
more specifically the `migrate_common.ph`, `migrate_group.pl` and
`migrate_passwd.pl`. Although the code from there is pretty archaic (even for
Perl long tradition of backwards compatibility), the definitions for OpenLDAP
seems to be pretty solid though.

This CLI **does not** supports Samba or NIS handling, only the regular files to
manage users and groups (`/etc/passwd` and `/etc/group`).

This CLI **does** supports `/etc/gshadow`, which is not included in the
[migrationtools](https://gitlab.com/future-ad-laboratory/migrationtools) already
mentioned files.

## Requirements

The minimal expected schemas to be available in the OpenLDAP server:

- cosine
- nis
- inetorgperson
- misc

## Development

The requirements are:

- Golang: see the `go.mod` file for details.
- GNU Make

See also the `Makefile` for the available targets.

## References

- [migrationtools project repository](https://gitlab.com/future-ad-laboratory/migrationtools)
