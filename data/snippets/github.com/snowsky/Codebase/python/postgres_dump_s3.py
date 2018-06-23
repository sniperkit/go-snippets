import os
import sys
import boto3
import tarfile
import subprocess
from optparse import OptionParser
from datetime import datetime

def main():
    parser = OptionParser()
    parser.add_option('-t', '--type', dest='backup_type',
                      help="Specify either 'hourly' or 'daily'.")
    parser.add_option('-d', '--dbhost', dest='db_host',
                      help="Database Host")
    parser.add_option('-n', '--dbname', dest='db_name',
                      help="Database Name")
    parser.add_option('-u', '--dbuser', dest='db_user',
                      help="Database User Name")
    parser.add_option('-p', '--dbpass', dest='db_pass',
                      help="Database User Password")
    parser.add_option('-b', '--s3', dest='s3_bucket',
                      help="S3 bucket name.")
    parser.add_option('-f', '--filename', dest='backup_file',
                      help="Absolute Path of Backup file ending with .tar.gz.")

    now = datetime.now()

    (options, args) = parser.parse_args()
    if options.backup_type == 'hourly':
        hour = str(now.hour).zfill(2)
    elif options.backup_type == 'daily':
        day_of_year = str(now.timetuple().tm_yday).zfill(3)
    else:
        parser.error('Invalid argument.')
        sys.exit(1)

    print 'Backing up %s database to %s' % (options.db_name, options.backup_file.split('.tar.gz')[0])
    ps = subprocess.Popen(
        ['pg_dump', 'postgresql://%s:%s@%s:5432/%s' % (options.db_user, options.db_pass, options.db_host, options.db_name), '-f', options.backup_file.split('.tar.gz')[0]],
        stdout=subprocess.PIPE
    )
    output = ps.communicate()[0]
    for line in output.splitlines():
        print line

    tar = tarfile.open(options.backup_file, "w:gz")
    tar.add(options.backup_file.split('.tar.gz')[0])
    tar.close()

    print 'Uploading %s to Amazon S3 bucket: %s...' % (options.backup_file, options.s3_bucket)
    upload_to_s3(options.backup_file, options.s3_bucket, os.path.basename(options.backup_file))


def upload_to_s3(source_path, s3_bucket, destination_filename):
    """
    Upload a file to an AWS S3 bucket.
    """
    s3 = boto3.resource('s3')
    s3.meta.client.upload_file(source_path, s3_bucket, destination_filename)


if __name__ == '__main__':
    main()
