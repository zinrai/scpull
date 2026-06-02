# scpull

Fetch a predefined set of files from one or more hosts via scp.

## What it does

scpull reads a YAML configuration file that defines named profiles. Each
profile lists a login user and a set of absolute remote paths. Given a
profile name and one or more hosts, scpull runs scp for each file and
saves it under a directory named after the host, preserving the remote
path structure.

The destination is always determined by the host, never by the profile.
The same host fetched through different profiles accumulates files under
one directory per host.

## Configuration

scpull reads `$HOME/.scpull.yaml` by default. Override it with `-config`.

    profiles:
      switch:
        user: admin
        paths:
          - /var/log/messages
          - /flash/syslog.txt
      server:
        user: deploy
        paths:
          - /etc/nginx/conf.d/site.conf

## Usage

Fetch the files in a profile from a single host:

    scpull switch 192.168.2.5

Fetch from several hosts in sequence:

    scpull switch 192.168.2.5 192.168.2.6 192.168.2.7

Use an alternate configuration file:

    scpull -config ./lab.yaml server web01

Given the `switch` profile above, fetching from `192.168.2.5`
produces:

    192.168.2.5/var/log/messages
    192.168.2.5/flash/syslog.txt

scp runs interactively, so authentication relies on your existing ssh
setup (keys, agent, ~/.ssh/config). scpull does not handle credentials.

## License

This project is licensed under the [MIT License](./LICENSE).
