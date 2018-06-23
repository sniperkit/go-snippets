guntalina
=========

**guntalina** is the tool for creating and executing command list based on:

* list of modified/overwritten/touched/chmoded/chowned/created files;
* action definitions;
* list of rules.

It's intended for use together with **gunter**, **guntalina**'s brother.

Consider the example, when configuration manager creates or changes
following files:

```
/etc/nginx/conf.d/sites-available.conf
/etc/nginx/nginx.conf
```

Somehow it should tell nginx to reload configuration to make use of new
changes.

**guntalina** can reload nginx for you based on changes in specified files.
All you need is to write configuration rule.

## Configuration

**guntalina** configuration file is written in the YAML language and consists
of two basic directives:

- `actions` - an array of action definitions, for example, `nginx-reload`, or
    `nginx-restart`, if you want.
- `rules` - an array of rule definitions, which should give an answer on the
    question like "*When some action should be invoked?*"

#### Actions

Action, in the **guntalina**'s meaning, is the named list of commands which
can be executed by rules. Commands should be declared in the `commands:`
section.

Commands will be executed sequentially in order. If command returned non-zero
exit code (e.g. failed), then execution will be aborted.

Let's write some typical actions for `nginx reload` and `nginx restart`:

```yaml
actions:
    nginx-reload:
        commands:
            - systemctl reload nginx

    nginx-restart:
        commands:
            - nginx -t # let's force check nginx config before real restart
            - systemctl restart nginx
```

> Note: actions section should be described as associative array
> (object), that's mean action name must not starts with `-`.

#### Rules

Basically, rules consists of two parts:

* list of actions, which should be executed; specified under `workflow`
  section;
* list of file path glob patterns, which will be looked for changes; specified
  under `masks` section;

Using `masks` you can handle files by glob patterns like `/etc/nginx/conf.d/*`
or `/etc/nginx/*.conf` and do specified `workflow`.

Let's proceed to write rules for our arbitary example with `nginx`:

```yaml
rules:
    - masks:
        - /etc/nginx/conf.d/*
        - /etc/nginx/nginx.conf
      workflow:
        - nginx-reload
```

So, action `nginx-reload` will be triggered if `/etc/nginx/nginx.conf` or some
file in the `/etc/nginx/conf.d/` directory will be changed.

When **guntalina** passes through the list of modified files and list of rules,
it will disable rules after their use, so if you have two or more modified
files which hit single rule, **guntalina** will execute matched workflow only once.

Actually, if rule contains the `group` directive, then **guntalina** will disable
all rules within same group.

Also, if you have some rules, which has repeatable workflow actions, **guntalina**
will execute actions only once.

Let's write more complex config, when **guntalina** should make a decision for
restart or reload for example software named **exampled**.

**exampled** should be restarted when files like as
`/etc/exampled/conf.d/*_cache_zone` changed. If any other file in the directory
`/etc/exampled/conf.d/` changed, then **exampled** should be reloaded.

For this case `group` directive should be used:

```yaml
actions:
    exampled-restart:
        commands:
            - exampled --check-config
            - examplectl stop
            - examplectl start

    exampled-reload:
        commands:
            - examplectl reload

rules:
    - group: exampled
      masks:
        - /etc/nginx/conf.d/*_cache_zone
      workflow:
        - exampled-restart

    - group: exampled
      masks:
        - /etc/nginx/conf.d/*
      workflow:
        -  exampled-reload
```

How you can see, both rules uses *group* directive which has `exampled`
value. It means, that **guntalina** will not run `exampled-reload`, if
`exampled-restart` already triggered, because after triggering
`exampled-restart` group `exampled` will be disabled.

> Note: rules section should be described as non-indexed array, that's mean
> rule item must starts with `-`.

#### Includes

**guntalina** has it's own YAML with Blackjack and includes. Files can be included
using `!include path/to/another/config`.

**guntalina**'s default configuration:

```yaml
actions:
    !include conf.d/*/actions

rules:
    !include conf.d/*/rules
```

Therefore, I recommend configuration structure, which looks like this:

```
.
├── conf.d
│   ├── haproxy
│   │   ├── actions
│   │   └── rules
│   └── nginx
│       ├── actions
│       └── rules
└── guntalina.conf
```

It's best way to provide `haproxy` and `nginx` directories as packages like
`haproxy-guntalina` and `nginx-guntalina` from your local repository, which
can be optional dependencies for `haproxy` and `nginx` upstream packages.

## Usage

#### Synopsis

- `-s <source>` - Specify source file, which is the list of
     modified/overwrited/created files. **required**
- `-c <config>` - Specify configuration file, as described in
     [Configuration](#Configuration) section.
- `-r --dry-run` - Dry-run mode, in this mode commands will not be executed,
     but printed on the stderr.
- `-f --force` - Do not stop execution if any command failed.
- `-v --version` - Show **guntalina**'s version.
- `-h --help` - Show help message.

**guntalina** has only one required argument: `-s <source>`, configuration
data will be read from `/etc/guntalina/guntalina.conf` file by default.
