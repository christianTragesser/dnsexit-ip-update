import os
from time import sleep
from dnsexitUpdate import utils
from dnsexitUpdate import logs

log = logs.logger('main')

login = os.environ['LOGIN'] if 'LOGIN' in os.environ else ''
password = os.environ['PASSWORD'] if 'PASSWORD' in os.environ else ''
domain = os.environ['DOMAIN'] if 'DOMAIN' in os.environ else ''
interval = int(os.environ['CHECK_INTERVAL']) if 'CHECK_INTERVAL' in os.environ else 600
update_data_url = os.environ['UPDATE_DATA_URL'] if 'UPDATE_DATA_URL' in os.environ else 'https://www.dnsexit.com/ipupdate/dyndata.txt'

if login == '' or password == '' or domain == '':
    log.critical('DNSExit account info missing. Login, password, and domain must be provided as environment variables.')
    exit(1)


def main(update_url, login, password, domain):
    sync_result = utils.evaluate_ip_sync(domain)
    if sync_result is None:
        log.error('ERROR: dnsexit-ip-update is not able to resolve {0:s} and should only be used to update existing DNS A records. Skipping update for {0:s}'.format(domain))
    elif not sync_result:
        try:
            utils.update_dns_a_record(update_fqdn, login, password, domain)
        except Exception as e:
            log.error(e)
            exit(1)


if __name__ == '__main__':
    log.info('INFO: Using DNSExit info login:{} domain:{}'.format(login, domain))
    if not utils.validate_credentials(login, password):
        exit(1)

    if not utils.validate_domain(login, domain):
        exit(1)

    update_fqdn = utils.get_update_url(update_data_url)

    while True:
        main(update_fqdn, login, password, domain)
        sleep(interval)
