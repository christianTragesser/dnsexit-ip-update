import requests
import socket
import logs

log = logs.logger('utils')


def get_update_url(data_url):
    update_data = requests.get(data_url)
    update_domain = update_data.text.split()[0].split('//')[-1]
    log.info('DNSExit update URL is https://'+update_domain)
    return 'https://'+update_domain


def evaluate_ip_sync(domain):
    egress_ip = requests.get('https://api.ipify.org').text
    dns_ip = socket.gethostbyname(domain)
    log.info('Evaluating DNS A record for {}: egress {} - dns {}'.format(domain, egress_ip, dns_ip))
    if egress_ip == dns_ip:
        log.info('DNS A record for {} is up to date.'.format(domain))
        return True
    else:
        log.info('Updating {} DNS A record.'.format(domain))
        return False


def update_dns_a_record(update_fqdn, user, password, domain):
    ip = requests.get('https://api.ipify.org').text
    update_query = '{}?login={}&password={}&host={}&myip={}'.format(update_fqdn, user, password, domain, ip)
    r = requests.get(update_query)
    log.info('{} DNS A record has been updated to {}.'.format(domain, ip))
    return r.status_code


def validate_credentials(login, password):
    creds_validation_url = 'https://update.dnsexit.com/ipupdate/account_validate.jsp?login={}&password={}'.format(login, password)
    r = requests.get(creds_validation_url)
    if '0=OK' in r.text:
        log.info('DNSExit IP Update credentials are valid.')
        return True
    else:
        log.error('The provided DNSExit IP Update credentials are not valid, exiting.')
        return False
