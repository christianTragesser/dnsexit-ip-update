import os
from time import sleep
import utils
import logs

log = logs.logger('main')

login = os.environ['LOGIN'] if 'LOGIN' in os.environ else ''
password = os.environ['PASSWORD'] if 'PASSWORD' in os.environ else ''
domain = os.environ['DOMAIN'] if 'DOMAIN' in os.environ else ''
interval = int(os.environ['CHECK_INTERVAL']) if 'CHECK_INTERVAL' in os.environ else 600
update_data_url = os.environ['UPDATE_DATA_URL'] if 'UPDATE_DATA_URL' in os.environ else 'https://www.dnsexit.com/ipupdate/dyndata.txt'
missing_acct_info = '''
*********************************************************************
dnsexit-ip-update requires DNSExit IP Update login, password, and
domain values provided as environment variables "LOGIN", "PASSWORD,
and "DOMAIN".

Please refer to our github page for more documentation:
https://github.com/christianTragesser/dnsexit-ip-update
*********************************************************************
'''

if login == '' or password == '' or domain == '':
    print(missing_acct_info)
    log.error('DNSExit account info missing. Login, password, and domain must be provided as environment variables.')
    exit(1)

if __name__ == '__main__':
    log.info('Using DNSExit info login:{} domain:{}'.format(login, domain))
    if not utils.validate_credentials(login, password):
        exit(1)
    if not utils.validate_domain(login, domain):
        exit(1)
    update_fqdn = utils.get_update_url(update_data_url)
    while True:
        if not utils.evaluate_ip_sync(domain):
            try:
                utils.update_dns_a_record(update_fqdn, login, password, domain)
            except Exception as e:
                log.error(e)
                exit(1)
        sleep(interval)
