from logger import logger


# output the changes that occur between occurrents.
class Changelog():
    def __init__(self):
        # Cache of the most recent version of each entity being
        # evaluated. The key is the (domain, name) pair and the value
        # is the entity itself.
        self.cache = {}

    def evaluate_attrs(self, current, previous):
        if current is None:
            current = {}

        if previous is None:
            previous = {}

        attrs = {}

        # Evaluate current items.
        for ck, cv in current.items():
            # Exists in previous.
            if ck in previous:
                # Values match, ignore.
                if cv == previous[ck]:
                    logger.debug('[changelog] attr matches: %s', ck)
                    continue

                logger.debug('[changelog] attr changed: %s', ck)

                attrs[ck] = {
                    'action': 'change',
                    'value': cv,
                    'previous': previous[ck],
                }
            else:
                logger.debug('[changelog] new attr: %s', ck)

                attrs[ck] = {
                    'action': 'add',
                    'value': cv,
                }

        for pk, pv in previous.items():
            # Already evaluated.
            if pk in current:
                continue

            logger.debug('[changelog] removed attr: %s', pk)
            # Removed in current state.
            attrs[pk] = {
                'action': 'remove',
                'previous': pv,
            }

        if not attrs:
            return

        return attrs

    def evaluate_refs(self, current, previous):
        refs = self.evaluate_attrs(current, previous)

        if not refs:
            return

        # Check for changes due to a timestamp change for the
        # same continuant.
        for k, v in refs.items():
            if v['action'] == 'change':
                c = v['value']
                p = v['previous']

                cts = c.get('timestamp')
                pts = p.get('timestamp')

                # Neither have timestamps.
                if not cts and not pts:
                    continue

                ck = (c['domain'], c['name'])
                pk = (p['domain'], p['name'])

                if ck != pk:
                    continue

                # TODO: Determine if referenced entity changed between
                # timestamps. This requires querying the previous state..
                logger.debug('[changelog] ref timestamps differ')

        return refs

    def evaluate(self, entity):
        # Without a name, the entity cannot be evaluated.
        if not entity['name']:
            return

        key = (entity['domain'], entity['name'])

        logger.debug('[changelog] evaluating: %s', key)

        event = {
            'labels': ('Diff', 'Add'),
            'domain': 'pedsnet.etlconv.changelog',
            'name': 'event_%s_%s' % (entity['batch'], entity['name']),
            'timestamp': entity['timestamp'],
            'refs': {
                'current': {
                    'domain': entity['domain'],
                    'name': entity['name'],
                    'batch': entity.get('batch'),
                    'timestamp': entity['timestamp'],
                },
                'previous': None,
            }
        }

        # Compare with previous if it exists.
        if key in self.cache:
            prev_entity, prev_event = self.cache[key]

            logger.debug('[changelog] compare: %s (%s -> %s)', key,
                         entity['batch'], prev_entity['batch'])

            # Compare attributes
            attrs = self.evaluate_attrs(entity.get('attrs'),
                                        prev_entity.get('attrs'))

            # Compare relationships
            refs = self.evaluate_refs(entity.get('refs'),
                                      prev_entity.get('refs'))

            # Nothing change from the last state.
            if not attrs and not refs:
                return

            event['refs']['previous'] = {
                'domain': prev_entity['domain'],
                'name': prev_entity['name'],
                'batch': prev_entity.get('batch'),
                'timestamp': prev_entity['timestamp'],
            }

            event['refs']['next'] = {
                'domain': prev_event['domain'],
                'name': prev_event['name'],
            }

            event['labels'] = ('Diff', 'Change')

            event['attrs'] = {
                'attrs': attrs,
                'refs': refs,
            }

        # Update the cache to point to the latest state and corresponding
        # event.
        self.cache[key] = (entity, event)

        return event
