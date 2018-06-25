from datetime import datetime


DATETIME_FORMATS = (
    '%Y-%m-%dT%H:%M:%SZ',
    '%Y-%m-%dT%H:%M:%S',
    '%Y-%m-%d',
)


def parse_date(s):
    for fmt in DATETIME_FORMATS:
        try:
            return datetime.strptime(s, fmt).timestamp()
        except ValueError:
            pass
