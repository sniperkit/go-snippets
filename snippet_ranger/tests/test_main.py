import argparse
from contextlib import contextmanager
from io import StringIO
import sys
import unittest
import logging

import snippet_ranger.__main__ as main


@contextmanager
def captured_output():
    log = StringIO()
    log_handler = logging.StreamHandler(log)
    logging.getLogger().addHandler(log_handler)
    new_out, new_err = StringIO(), StringIO()
    old_out, old_err = sys.stdout, sys.stderr
    try:
        sys.stdout, sys.stderr = new_out, new_err
        yield sys.stdout, sys.stderr, log
    finally:
        sys.stdout, sys.stderr = old_out, old_err
        logging.getLogger().removeHandler(log_handler)


class MainTests(unittest.TestCase):
    def test_handlers(self):
        action2handler = {
            "source2func": "source2func_entry",
            "dependent_reps": "dependent_reps_entry",
            "snippet2df": "snippet2df_entry",
            "snippet2bow": "snippet2bow_entry",
            "snippet2fc_bow": "snippet2fc_bow_entry",
            "snippet2fc_df": "snippet2fc_df_entry",
            "pylib2uast": "pylib2uast_entry"

        }
        parser = main.get_parser()
        subcommands = set([x.dest for x in parser._subparsers._actions[2]._choices_actions])
        set_action2handler = set(action2handler)
        self.assertFalse(len(subcommands - set_action2handler),
                         "You forgot to add to this test {} subcommand(s) check".format(
                             subcommands - set_action2handler))

        self.assertFalse(len(set_action2handler - subcommands),
                         "You cover unexpected subcommand(s) {}".format(
                             set_action2handler - subcommands))

        called_actions = []
        args_save = sys.argv
        error_save = argparse.ArgumentParser.error
        try:
            argparse.ArgumentParser.error = lambda self, message: None

            for action, handler in action2handler.items():
                def handler_append(*args, **kwargs):
                    called_actions.append(action)

                handler_save = getattr(main, handler)
                try:
                    setattr(main, handler, handler_append)
                    sys.argv = [main.__file__, action]
                    main.main()
                finally:
                    setattr(main, handler, handler_save)
        finally:
            sys.argv = args_save
            argparse.ArgumentParser.error = error_save

        set_called_actions = set(called_actions)
        set_actions = set(action2handler)
        self.assertEqual(set_called_actions, set_actions)
        self.assertEqual(len(set_called_actions), len(called_actions))

    def test_empty(self):
        args = sys.argv
        error = argparse.ArgumentParser.error
        try:
            argparse.ArgumentParser.error = lambda self, message: None

            sys.argv = [main.__file__]
            with captured_output() as (stdout, _, _):
                main.main()
        finally:
            sys.argv = args
            argparse.ArgumentParser.error = error
        self.assertIn("usage:", stdout.getvalue())


if __name__ == "__main__":
    unittest.main()
