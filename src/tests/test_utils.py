import responses
import mock
import os
import sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from utils import get_update_url
from utils import evaluate_ip_sync
from utils import update_dns_a_record

current_ip_resource = 'https://api.ipify.org'


@responses.activate
def test_get_update_url():
    # retrieves update data
    # filters update URL from update data
    # returns update FQDN and path
    update_data_url = 'https://test.local/ipupdate/dyndata.txt'
    update_data = '''
    url=http://update.test.local/TestUpdate.sv
    base=http://update.test.local/
    version=2.0
    '''
    responses.add(responses.GET, update_data_url, body=update_data)

    update_addr = get_update_url(update_data_url)
    assert update_addr == 'https://update.test.local/TestUpdate.sv'


@responses.activate
@mock.patch('socket.gethostbyname', return_value='2.2.2.2')
def test_evaluate_ip_synced(mock_dns_lookup):
    # lookup current egress IP address
    # lookup current DNS A record for domain
    # compare current IP with DNS IP
    # return True when matching
    responses.add(responses.GET, current_ip_resource, body='2.2.2.2')

    sync = evaluate_ip_sync('test.local')
    assert sync


@responses.activate
@mock.patch('socket.gethostbyname', return_value='4.4.4.4')
def test_evaluate_ip_unsynced(mock_dns_lookup):
    # lookup current egress IP address
    # lookup current DNS A record for domain
    # compare current IP with DNS IP
    # return False when not matching
    responses.add(responses.GET, current_ip_resource, body='2.2.2.2')

    sync = evaluate_ip_sync('test.local')
    assert not sync


@responses.activate
def test_update_dns_a_record():
    # takes in update url, user, password, and domain
    # performs update query
    # reports update query result
    update_url = 'https://update.test.local/TestUpdate.sv'
    user = 'tester'
    password = 'Hello123'
    domain = 'test.local'
    ip = '2.2.2.2'
    update_query = '{0:s}?login={1:s}&password={2:s}&host={3:s}&myip={4:s}'.format(update_url, user, password, domain, ip)
    responses.add(responses.GET, current_ip_resource, body=ip)
    responses.add(responses.GET, update_query, status=200)

    result = update_dns_a_record(update_url=update_url, user=user, password=password, domain=domain)
    assert result == 200
