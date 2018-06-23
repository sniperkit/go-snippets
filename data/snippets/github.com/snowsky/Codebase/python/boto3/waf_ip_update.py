"""
Update an AWS IP set with the the top 10 IP addresses from Apache logs
"""
import sys
import re
import boto3
from collections import Counter
from tabulate import tabulate

def get_top_n_ip_addresses(num):
    """
    Parses apache logs to find the top n X-Forwarded-For ip addresses
    """
    all_ip_addresses = []

    with open('/var/log/httpd/access_log') as content:
        for line in content:
            # the first "cell" surrounded with brackets is the X-Forwarded-For
            regex = re.search(r'\((.*?)\)', line)
            if regex:
                # format is X-Forwarded-For: client, proxy1, proxy2
                # we want the right most IP, the IP that hit cloudfront
                forwarded_for = re.findall(r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', regex.group(0))
                if forwarded_for:
                    all_ip_addresses.append(forwarded_for[-1])

    return Counter(all_ip_addresses).most_common(num)

def main():
    """
    Grab top 10 X-Forwarded-For ip addresses and send to a WAF ip list
    """
    top_ip_addresses = get_top_n_ip_addresses(10)

    print
    print "Top 10 IP Addresses"
    print "==================="
    print
    print tabulate(top_ip_addresses, headers=["IP", "Count"])
    print

    updates_list = [{
        'Action': 'INSERT',
        'IPSetDescriptor': {
            'Type': 'IPV4',
            'Value': "%s/32" % ip[0]
        }
    } for ip in top_ip_addresses]

    waf = boto3.client('waf')
    waf_ip_sets = waf.list_ip_sets(
        Limit=100
    )['IPSets']

    if len(waf_ip_sets) < 1:
        sys.exit('WAF IP sets appear to be misconfigured.  Expecting 1 IP set.')

    waf_ip_set_id = waf_ip_sets[0]['IPSetId']
    print "Updating IP set: ", waf_ip_sets[0]['Name']

    waf.update_ip_set(IPSetId=waf_ip_set_id,
                      ChangeToken=waf.get_change_token()['ChangeToken'],
                      Updates=updates_list)

    print "Done!"

if __name__ == "__main__":
    main()
