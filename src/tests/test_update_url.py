import pytest
import responses
import mock
import os, sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from ip_update import get_update_url, evaluate_ip_sync

update_data_url = 'https://test.local/ipupdate/dyndata.txt'
update_data = '''
url=http://update.test.local/TestUpdate.sv
base=http://update.test.local/
version=2.0
'''
@responses.activate
def test_get_update_url():
    #retrieves update data
    #filters update URL from update data
    #returns update FQDN and path
    responses.add(responses.GET, update_data_url, body=update_data)

    update_addr = get_update_url(update_data_url)
    assert update_addr == 'https://update.test.local/TestUpdate.sv'


current_ip_resource = 'https://api.ipify.org'
@responses.activate
@mock.patch('socket.gethostbyname', return_value='2.2.2.2')
def test_evaluate_ip_synced(mock_dns_lookup):
    #lookup current egress IP address
    #lookup current DNS A record for domain
    #compare current IP with DNS IP
    #return True when matching
    responses.add(responses.GET, current_ip_resource, body='2.2.2.2')

    sync = evaluate_ip_sync('test.local')
    assert sync


current_ip_resource = 'https://api.ipify.org'
@responses.activate
@mock.patch('socket.gethostbyname', return_value='4.4.4.4')
def test_evaluate_ip_unsynced(mock_dns_lookup):
    #lookup current egress IP address
    #lookup current DNS A record for domain
    #compare current IP with DNS IP
    #return False when not matching
    responses.add(responses.GET, current_ip_resource, body='2.2.2.2')

    sync = evaluate_ip_sync('test.local')
    assert sync == False