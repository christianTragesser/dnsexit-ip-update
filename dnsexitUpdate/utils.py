import socket
import requests
from dnsexitUpdate import logs
from dns.resolver import Resolver

log = logs.logger('utils')
dnsexit_nameservers = ['ns1.dnsexit.com', 'ns2.dnsexit.com', 'ns3.dnsexit.com', 'ns4.dnsexit.com']
nameservers = [socket.gethostbyname(ns) for ns in dnsexit_nameservers]
resolve = Resolver()
resolve.nameservers = nameservers


def get_update_url(data_url):
    update_data = requests.get(data_url)
    update_domain = update_data.text.split()[0].split('//')[-1]
    log.info('INFO: DNSExit update URL is https://' + update_domain)
    return 'https://' + update_domain


def dns_lookup(domain):
    try:
        answers = resolve.query(domain)
        return tuple((rdata.address for rdata in answers))
    except Exception as e:
        log.error('ERROR: {}'.format(e))
        return ()


def evaluate_ip_sync(domain):
    egress_ip = requests.get('https://api.ipify.org').text
    dns_ips = dns_lookup(domain)

    if dns_ips == ():
        return None

    if isinstance(dns_ips, tuple):
        dns_list = ''
        for ip in dns_ips:
            dns_list = dns_list + ip + ' '
    else:
        dns_list = dns_ips

    log.info('INFO: Evaluating DNS A record for {}, egress {} - dns {}'.format(domain, egress_ip, dns_list))

    if egress_ip in dns_ips:
        log.info('INFO: DNS A record for {} is up to date.'.format(domain))
        return True
    else:
        log.info('INFO: Updating {} DNS A record.'.format(domain))
        return False


def update_dns_a_record(update_fqdn, user, password, domain):
    ip = requests.get('https://api.ipify.org').text
    update_query = '{}?login={}&password={}&host={}&myip={}'.format(update_fqdn, user, password, domain, ip)
    # this needs better response handling, update endpoint always returns 200
    r = requests.get(update_query)
    log.info('INFO: DNSExit IP Update service has been notified to use IP address {} for domain {}.'.format(ip, domain))
    return r.status_code


def validate_credentials(login, password):
    creds_validation_url = 'https://update.dnsexit.com/ipupdate/account_validate.jsp?login={}&password={}'.format(login, password)
    r = requests.get(creds_validation_url)
    if '0=OK' in r.text:
        log.info('INFO: DNSExit IP Update credentials are valid.')
        return True
    else:
        log.error('ERROR: The provided DNSExit IP Update credentials are not valid, exiting.')
        return False


def validate_domain(login, domain):
    domain_validation_url = 'https://update.dnsexit.com/ipupdate/domains.jsp?login={}'.format(login)
    r = requests.get(domain_validation_url)
    if '0=' in r.text and domain in r.text:
        log.info('INFO: {} domain is valid.'.format(domain))
        return True
    else:
        log.error('ERROR: {0:s} domain is invalid, {0:s} not found in {1:s} account.'.format(domain, login))
        return False
