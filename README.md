![orbi](https://github.com/bronzdoc/orbi/blob/master/logo.png)
[![Build Status](https://travis-ci.org/bronzdoc/orbi.svg?branch=master)](https://travis-ci.org/bronzdoc/orbi)

> Project structure generator.

Generate project structures using yaml and golang templating.

# Install
> NOTE: orbi just works with \*nix operating systems, windows is not supported for now.

### Binaries

- **linux** [386](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-linux-386) / [amd64](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-linux-amd64) / [arm](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-linux-arm) / [arm64](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-linux-arm64)
- **darwin** [386](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-darwin-386) / [amd64](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-darwin-amd64)
- **freebsd** [386](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-freebsd-386) / [amd64](https://github.com/bronzdoc/orbi/releases/download/v0.1.1/orbi-freebsd-amd64)

### Via Go

```shell
$ go get github/bronzdoc/orbi
```

# Usage

Orbi defines project structures using a `definition.yml`.

```yaml
---
context: .
resources:
  - dir:
     name: dir_1
     files:
      - file_a

  - dir:
     name: dir_2
     files:
      - file_b
     dir:
      name: dir_3
      files:
        - file_c
        - file_d

  - files:
     - file_e
     - file_f
```

A `context` is where your `resources` structure will be created.

The way orbi organize definitions is with something called a `plan`, you can create a new plan by doing:

```shell
$ orbi plan new my_plan
```

This command will generate the following:

```shell
$ tree ~/.orbi/plans/my_plan

my_plan
├── definition.yml
└── templates
```

You can notice a `templates` directory, this is where your templates should be.

In order to template a file all you need to do is create a file named the same as a file resource.

```yaml
context: .
resources:
  files:
    - file_a
```

```shell
├── definition.yml
└── templates
    └── file_a
```

Ok, we have a plan with a definition.yml and templates, so... how we create all that stuff we defined?

All you need to do is tell orbi to execute a `plan`:

```shell
$ orbi exec my_plan
```

this command will generate the file structure defined in your plan definition.yml.

If your plan templates happen to have variables, you can pass values to those variables with the `--vars` flag.

in `.orbi/plans/my_plan/templates/file_a`

```
{{.name}} is awesome
```

```shell
$ orbi exec my_plan --vars="name=Tarantino"
```

that command will generate the file named `file_a` with content `Tarantino is awesome`.
> NOTE: you can also pass a KEY=VALUE variables file with `--vars-file`


### Sharing plans
orbi let you download a plan from a repository with the `orbi plan get` command

**ssh:**
```shell
$ orbi plan get my_plan git@github.com:user/plan_name.git
```

**https:**
```shell
$ orbi plan get my_plan https://user@github.com/user/plan_name.git
```

## Contributing

Contributions are greatly appreciated and encouraged, see [CONTRIBUTING](https://github.com/bronzdoc/orbi/blob/master/CONTRIBUTING.md)
