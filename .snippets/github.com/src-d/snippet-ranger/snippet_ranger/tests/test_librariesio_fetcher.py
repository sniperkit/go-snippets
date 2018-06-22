import os

import unittest
import tempfile

from snippet_ranger.librariesio_fetcher import LibrariesIOFetcher
from snippet_ranger.tests.models import DATA_DIR


class LibrariesIOFetcherTests(unittest.TestCase):
    LIBRARIES_IO_DATA_PATH = os.path.join(DATA_DIR, "test_librariesio_data/")

    def setUp(self):
        self.libio = LibrariesIOFetcher(LibrariesIOFetcherTests.LIBRARIES_IO_DATA_PATH)

    def get_url_str(self, lib, platform):
        with tempfile.NamedTemporaryFile() as f:
            self.libio.get_dependent_rep_urls(lib, platform, f.name)
            with open(f.name) as f_urls:
                urls = f_urls.read()
        return set(urls.split("\n")[:-1])

    def test_get_dependent_rep_urls(self):
        url1 = set(["https://github.com/repo1/repo1"])
        url2 = set(["https://github.com/repo2/repo2"])

        urls = self.get_url_str({"lib1": "lib1.url"}, "")
        self.assertEqual(urls, url1)

        urls = self.get_url_str({"lib2": "lib2.url"}, "Platform2")
        self.assertEqual(urls, url1 | url2)
        with self.assertRaises(ValueError) as _:
            self.get_url_str({"lib2": "lib2.url"}, "wrong_platform")

        urls = self.get_url_str({"lib1": "lib1.url", "lib2": "lib2.url"}, "")
        self.assertEqual(urls, url1 | url2)

        with tempfile.TemporaryDirectory() as tmpdir:
            self.libio.get_dependent_rep_urls(
                {"lib1": "lib1.url", "lib2": "lib2.url"}, "", tmpdir)

            with open(os.path.join(tmpdir, "lib1.txt")) as f_urls:
                urls = f_urls.read()
            self.assertEqual(set(urls.split("\n")[:-1]), url1)

            with open(os.path.join(tmpdir, "lib2.txt")) as f_urls:
                urls = f_urls.read()
            self.assertEqual(set(urls.split("\n")[:-1]), url1 | url2)

if __name__ == "__main__":
    unittest.main()
