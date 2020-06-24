import os
from time import sleep
import utils

login = os.environ['LOGIN'] if 'LOGIN' in os.environ else ''
password = os.environ['PASSWORD'] if 'PASSWORD' in os.environ else ''
host = os.environ['HOST'] if 'HOST' in os.environ else ''
interval = int(os.environ['CHECK_INTERVAL']) if 'CHECK_INTERVAL' in os.environ else 600
update_data_url = os.environ['UPDATE_DATA_URL'] if 'UPDATE_DATA_URL' in os.environ else 'https://www.dnsexit.com/ipupdate/dyndata.txt'


if __name__ == '__main__':
    print('Initializing...')
    update_fqdn = utils.get_update_url(update_data_url)
    while True:
        if not utils.evaluate_ip_sync(host):
            try:
                utils.update_dns_a_record(update_fqdn, login, password, host)
            except Exception as e:
                print(e)
                exit(1)
        sleep(interval)
