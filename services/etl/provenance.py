"""
Generate extracts a set of entities from the model and corresponding
commit relevant to the ETL conventions document. The set of entities
consist of ETL conventions data for the model, tables, and fields as
well as provenance information including causual or associated events
and the people or systems that were involved or are responsible.

All information produced here is stateless since historical information
is not available.

The structure of an entity is as follows:

- `domain` - The domain the entity applies to.
- `name` - The unique identity value of the entity. An identity should
be provided if the entity is a continuant in the domain so all occurrents
can be associated.
- `labels` - A set of labels for the entity. Each label provides a
secondary index for the entity.
- `attrs` - A set of key-value pairs where the key is a identity value
and the value is a literal or identity value.
- `refs` - A set of references to entities that are related.

An identity value is a reference to another entity. A reference value
consists of a domain and name. A string encoded reference value is
`domain:name`, that is the two parts delimited by a colon. Alternately,
an ident can represented as an object with keys `domain` and `name`.
"""
from copy import deepcopy
from version import __version__
from dateutil import parse_date


SPEC_VERSION = '1.0.0'


def validate(message):
    if not message.get('timestamp'):
        raise ValueError('timestamp is required')

    if not message.get('spec'):
        raise ValueError('spec is required')

    if 'attrs' not in message:
        return message

    for key, value in message['attrs'].items():
        if not isinstance(key, str):
            raise TypeError('attribute keys must be strings')

        # Reference value.
        if isinstance(value, dict):
            for k, v in value.items():
                if k != 'domain' and k != 'ident':
                    raise KeyError('only domain and ident keys are allowed in '
                                   'reference values')

                if not isinstance(v, str):
                    raise TypeError('domain and ident values must be strings '
                                    'not {}'.format(type(v)))

        elif not isinstance(value, (bool, str, float, int)):
            raise TypeError('unknown type {}. supported types are bool, '
                            'str, float, and int'.format(type(value)))

    return message


class entity():
    def __init__(self, data, timestamp, batch=None):
        self.batch = batch or ''
        self.timestamp = timestamp
        self.spec = data.get('spec', SPEC_VERSION)
        self.domain = data.get('domain', '')
        self.name = data.get('name', '')
        self.labels = data.get('labels', ())
        self.attrs = data.get('attrs', {})
        self.refs = data.get('refs', {})

    def __getitem__(self, key):
        if key.startswith('_'):
            raise KeyError

        return self.__dict__[key]

    def __setitem__(self, key, value):
        self.attrs[key] = value

    @property
    def ident(self):
        return {
            'domain': self.domain,
            'name': self.name,
        }

    def json(self):
        return validate({
            'batch': self.batch,
            'timestamp': self.timestamp,
            'spec': self.spec,
            'domain': self.domain,
            'name': self.name,
            'labels': tuple(set(self.labels)),
            'attrs': self.attrs,
            'refs': self.refs,
        })


def generate(file_name, domain, model, commit):
    # The unique batch ID is the commit SHA.
    batch = commit['sha']

    # Timestamp of when the entity became available.
    timestamp = parse_date(commit['commit']['committer']['date'])

    entities = []

    # Committer and author may be the same person, but entities
    # are deduped downstream.
    committer = commit['commit']['committer']

    committer = entity({
        'domain': domain,
        'name': committer['email'],
        'labels': ['Person', 'Agent'],
        'attrs': committer,
    }, timestamp, batch=batch)

    author = commit['commit']['author']

    author = entity({
        'domain': domain,
        'name': author['email'],
        'labels': ['Person', 'Agent'],
        'attrs': author,
    }, timestamp, batch=batch)

    # The git commit that corresponds to the current state. The metadata
    # of the commit does not technically need to be copied here. A consumer
    # of this data could lookup the commit and fetch the data manually.
    commit = entity({
        'domain': domain,
        'name': commit['sha'],
        'labels': ['Commit'],
        'attrs': {
            'sha': commit['sha'],
            'url': commit['commit']['url'],
            'message': commit['commit']['message'],
            'commit_time': commit['commit']['committer']['date'],
            'author_time': commit['commit']['author']['date'],
        },
        'refs': {
            'committer': committer.ident,
            'author': author.ident,
        }
    }, timestamp, batch=batch)

    # The source file containing the data the model, tables, and fields
    # were extracted from.
    source_file = entity({
        'domain': domain,
        'name': file_name,
        'labels': ['File'],
        'attrs': {
            'path': file_name,
        },
        'refs': {
            'commit': commit.ident,
        },
    }, timestamp, batch=batch)

    service = entity({
        'domain': domain,
        'name': 'pedsnet/etlconv',
        'labels': ['Agent', 'Service'],
        'attrs': {
            'name': 'PEDSnet ETL Conventions Service',
            'version': __version__,
        },
    }, timestamp, batch=batch)

    # Append the entities.
    entities.append(service)
    entities.append(source_file)
    entities.append(commit)
    entities.append(author)
    entities.append(committer)

    # The base event is copied for each entity. This event describes
    # the relationship between the relationship between the Git commit
    # and the entity.
    base_event = {
        'domain': domain,
        'name': 'event_%s' % commit['name'],
        'labels': ['Event', 'EntitiesExtracted'],
        'attrs': {
            'event': 'EntitiesExtracted',
        },
        'refs': {
            'file': source_file.ident,
            'service': service.ident,
        },
    }

    model = dict(model)
    tables = model.pop('tables')

    model = entity({
        'domain': domain,
        'name': model['name'],
        'labels': ['Model'],
        'attrs': model,
    }, timestamp, batch=batch)

    entities.append(model)

    event = entity(deepcopy(base_event), timestamp, batch=batch)
    event['refs']['entity'] = model.ident
    entities.append(event)

    for table_name, attrs in tables.items():
        attrs = dict(attrs)
        fields = attrs.pop('fields')
        attrs['name'] = table_name

        table_id = table_name

        table = entity({
            'domain': domain,
            'name': table_id,
            'labels': ['Table'],
            'attrs': attrs,
            'refs': {
                'model': model.ident,
            },
        }, timestamp, batch=batch)

        entities.append(table)

        event = entity(deepcopy(base_event), timestamp, batch=batch)
        event['refs']['entity'] = table.ident
        entities.append(event)

        for field_name, attrs in fields.items():
            attrs = dict(attrs)
            attrs['name'] = field_name

            field_id = '{}.{}'.format(table_name, field_name)

            field = entity({
                'domain': domain,
                'name': field_id,
                'labels': ['Field'],
                'attrs': attrs,
                'refs': {
                    'table': table.ident,
                    'model': model.ident,
                },
            }, timestamp, batch=batch)

            entities.append(field)

            event = entity(deepcopy(base_event), timestamp, batch=batch)
            event['refs']['entity'] = field.ident
            entities.append(event)

    return entities
