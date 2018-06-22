import json
import httplib2
import sys
import time
import socket


# User downloader
class UserGet(object):
    def __init__(self, oauth_id, oauth_secret, since, stop, outdir):
        self.oauth_id = oauth_id
        self.oauth_secret = oauth_secret

        self.since = since
        self.stop = since
        self.first_user = since
        self.url = "https://api.github.com/users"

        self.log_file = open(outdir + "/getter.log", "w")
        self.output_dir = outdir

    # Get one page of users with id > since
    def get_page(self, since):
        """Gets a page of users such that id > since. Waits for"""
        h = httplib2.Http()  # (".cache")

        url = self.url + "?since=%d" % (since)

        if self.oauth_id and self.oauth_secret:
            url += "&client_id=%s&client_secret=%s" % (self.oauth_id,
                                                       self.oauth_secret)

        # self.log("Querying " + url)
        r, content = h.request(url, "GET")

        return r, content

    # Gets all pages
    def get_all(self):
        self.log("Starting...")
        users = []

        last_user = self.since

        try:  # Catches KeyboardInterrupts to shutdown gracefully
            while last_user < stop:
                # Send request, repeat if network problem
                try:
                    r, content = self.get_page(self.since)
                except httplib2.HttpLib2Error as e:
                    self.log("Httplib2 error: %s" % str(e))
                    self.log("Trying again...")
                    continue
                except socket.error as e:
                    self.log("Socket error %d: %s" % (e.errno, e.strerror))
                    self.log("Trying again...")
                    continue

                # Check the response status code to see if the request was
                # successful
                if r['status'] == '200':
                    jcontent = json.loads(content)

                    # If we don't get new users, we stop
                    if len(jcontent) == 0 or self.since == jcontent[-1]['id']:
                        self.log("Last request didn't return new users. \
                                 Stopping!")
                        self.dump(users)
                        return
                    else:
                        self.since = jcontent[-1]['id']
                        last_user = self.since

                    users.extend(jcontent)

                else:
                    # If the request was not succesful, print headers and wait
                    # a little bit
                    self.log("Received return status: %s" % r['status'])
                    self.log(str(r))

                    time.sleep(3)

                # Check the number remaining API calls
                remaining_calls = int(r['x-ratelimit-remaining'])

                if remaining_calls == 0:
                    waittime = int(r['x-ratelimit-reset']) - time.time()
                    self.log("Waiting %d minutes for more API calls"
                             % (waittime / 60))
                    time.sleep(waittime)

                # Dump users if we have more than 5000
                if len(users) > 5000:
                    self.log("Remaining API calls: %d \tLast user obtained: %d"
                             % (remaining_calls, self.since))

                    self.dump(users)
                    users = []

        except KeyboardInterrupt:
            # Ignore exception and jump to finally
            pass

        finally:
            # Close gracefully
            self.log("Writing files...")

            self.dump(users)

            if len(users) > 0:
                self.log("Last user written: %d" % users[-1]['id'])
            else:
                self.log("No user was fetched")

            self.log_file.flush()
            self.log_file.close()

    # Writes msg both to stdout and to the log file
    def log(self, msg):
        print(("[%d] UserGetter: " % (int(time.time()))) + msg)
        self.log_file.write("[" + str(time.time()) + "] " + msg + "\n")

    # Dumps the list of users to a file
    def dump(self, users):
        if len(users) > 0:
            out_name = self.output_dir + "/users"
            out_name += "_" + str(users[0]['id']).zfill(7)
            out_name += "_" + str(users[-1]['id']).zfill(7)
            out_name += ".json"

            output = open(out_name, "a")
            output.write(json.dumps(users))
            output.flush()
            output.close()


if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("Usage: %s <first user> <last user> <output dir> "
              "[oauth id] [oauth secret]" % (sys.argv[0]))
        print("Gets all the users such that "
              "[first user < user id <= last user id ]")
        print("When it runs out of API calls it waits")
        sys.exit(-1)
    else:
        since = int(sys.argv[1])
        stop = int(sys.argv[2])
        outdir = sys.argv[3]

    if len(sys.argv) > 5:
        oauth_id = sys.argv[4]
        oauth_secret = sys.argv[5]
    else:
        oauth_id = None
        oauth_secret = None

    getter = UserGet(oauth_id, oauth_secret, since, stop, outdir)
    getter.get_all()
