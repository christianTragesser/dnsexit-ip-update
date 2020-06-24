import requests
import socket


def get_update_url(data_url):
    update_data = requests.get(data_url)
    update_domain = update_data.text.split()[0].split('//')[-1]
    print('DNSExit update URL is https://'+update_domain+'\n')
    return 'https://'+update_domain


def evaluate_ip_sync(domain):
    egress_ip = requests.get('https://api.ipify.org').text
    dns_ip = socket.gethostbyname(domain)
    print('Evaluating DNS A record for {}'.format(domain))
    print('\tEgress IP is {}, DNS IP {}'.format(egress_ip, dns_ip))
    if egress_ip == dns_ip:
        print('\tDNS A record for {0:s} is up to date.'.format(domain))
        return True
    else:
        print('\tUpdating {0:s} DNS A record.'.format(domain))
        return False


def update_dns_a_record(update_fqdn, user, password, domain):
    ip = requests.get('https://api.ipify.org').text
    update_query = '{0:s}?login={1:s}&password={2:s}&host={3:s}&myip={4:s}'.format(update_fqdn, user, password, domain, ip)
    r = requests.get(update_query)
    print('{0:s} DNS A record has been updated to {1:s}.'.format(domain, ip))
    return r.status_code
