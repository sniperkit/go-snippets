# OpenLDAP SLAPD on Alpine Linux

[![Docker Stars](https://img.shields.io/docker/stars/dweomer/openldap.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/dweomer/openldap.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/dweomer/openldap/latest.svg)][layers]
[![Image Layers](https://img.shields.io/imagelayers/layers/dweomer/openldap/latest.svg)][layers]

The [_`memberOf`_](http://www.openldap.org/software/man.cgi?query=slapo-memberof&sektion=5) (with [_refint_](http://www.openldap.org/software/man.cgi?query=slapo-refint&sektion=5)) overlay is installed and configured for [`groupOfUniqueNames`](https://tools.ietf.org/html/rfc4519#section-3.6).

## Usage
```
docker run -itd -p 389:389 dweomer/openldap
```

## Modifying Init/Runtime Behavior
### Environment
#### `declare -x LDAP_DOMAIN`
Defaults to `example.com` if not overridden.

#### `declare -x LDAP_DOMAIN_OBJECTCLASS`
Default value is [`domain`](https://tools.ietf.org/html/rfc4524#section-3.4). Supports any [`objectClass`](https://tools.ietf.org/html/rfc4512#section-3.3) or combination thereof that allows for the [`dc` aka `domainComponent`](https://tools.ietf.org/html/rfc4519#section-2.4) and [`o` aka `organizationName`](https://tools.ietf.org/html/rfc4519#section-2.19) attributes, e.g.

```
export LDAP_DOMAIN_OBJECTCLASS="organization
objectClass: dcObject"
```

#### `declare -x LDAP_SUFFIX`
By default this is generated from `LDAP_DOMAIN`, e.g. `dc=example,dc=com` for the default value of such, but one could as easily pass `o=example.com` or any other legitimate [`distinguishedName`](https://tools.ietf.org/html/rfc4512#section-2.3.2) supported by the [`objectClass`](https://tools.ietf.org/html/rfc4512#section-3.3).

#### `declare -x LDAP_ORGANIZATION`
Defaults to the value of `LDAP_DOMAIN` if not overridden.

#### `declare -x LDAP_PASSWORD`
If not specified this is the string 'lderp!' concatenated with the first [`domainComponent`](https://tools.ietf.org/html/rfc4519#section-2.4) from `LDAP_DOMAIN`. So, `lderp!example` for the default value of `LDAP_DOMAIN`. This will be the password for the bind [`dn`](https://tools.ietf.org/html/rfc4512#section-2.3.2) of `cn=admin,${LDAP_SUFFIX}`, e.g. `cn=admin,dc=example,dc=com`.

### Database
The first time the container starts up it will look for any files under the `/srv/openldap.d/` hierarchy and process them in lexical order. Right now the only supported file extensions are `.sh` and `.ldif`, everything else is ignored.

If not already present, `/srv/openldap.d/000-domain.ldif` is created prior to such processing so that it will be picked up. This will create an [`organizationalUnit`](https://tools.ietf.org/html/rfc4519#section-3.11) for user accounts with [RDN](https://tools.ietf.org/html/rfc4512#section-2.3.1) `cn=users` and for user groups with [RDN](https://tools.ietf.org/html/rfc4512#section-2.3.1) `cn=groups`.

## License

See the `LICENSE` file in this repository.

[hub]: https://hub.docker.com/r/dweomer/openldap/
[issues]: https://github.com/dweomer/dockerfiles-openldap/issues
[layers]: https://imagelayers.io/?images=dweomer/openldap:latest
