import responses
import mock
import os
import sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from utils import get_update_url
from utils import evaluate_ip_sync
from utils import update_dns_a_record
from utils import validate_credentials 

current_ip_resource = 'https://api.ipify.org'


@responses.activate
def test_get_update_url(caplog):
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
    assert 'DNSExit update URL is https://update.test.local/TestUpdate.sv' in caplog.text


@responses.activate
@mock.patch('socket.gethostbyname', return_value='2.2.2.2')
def test_evaluate_ip_synced(mock_dns_lookup, caplog):
    # lookup current egress IP address
    # lookup current DNS A record for domain
    # compare current IP with DNS IP
    # return True when matching
    responses.add(responses.GET, current_ip_resource, body='2.2.2.2')

    sync = evaluate_ip_sync('test.local')
    assert sync
    assert 'Evaluating DNS A record for test.local' in caplog.text
    assert 'egress 2.2.2.2 - dns 2.2.2.2' in caplog.text
    assert 'DNS A record for test.local is up to date.' in caplog.text


@responses.activate
@mock.patch('socket.gethostbyname', return_value='4.4.4.4')
def test_evaluate_ip_unsynced(mock_dns_lookup, caplog):
    # lookup current egress IP address
    # lookup current DNS A record for domain
    # compare current IP with DNS IP
    # return False when not matching
    responses.add(responses.GET, current_ip_resource, body='2.2.2.2')

    sync = evaluate_ip_sync('test.local')
    assert not sync
    assert 'Evaluating DNS A record for test.local:' in caplog.text
    assert 'egress 2.2.2.2 - dns 4.4.4.4' in caplog.text
    assert 'Updating test.local DNS A record.' in caplog.text


@responses.activate
def test_update_dns_a_record(caplog):
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

    result = update_dns_a_record(update_fqdn=update_url, user=user, password=password, domain=domain)
    assert result == 200
    assert 'test.local DNS A record has been updated to 2.2.2.2.' in caplog.text


@responses.activate
def test_valid_credentials(caplog):
    # takes in login and password
    # validates credentials with DNSExit
    # parse text response
    # log credentials are valid
    # return true if credentials are valid
    login = 'tester'
    password = 'Hello123'
    creds_validation_url = 'https://update.dnsexit.com/ipupdate/account_validate.jsp?login={}&password={}'.format(login, password)
    responses.add(responses.GET, creds_validation_url, body='\r\n\r\n\r\n\r\n\r\n\r\n\r\n0=OK\r\n')

    result = validate_credentials(login=login, password=password)
    assert result
    assert 'DNSExit IP Update credentials are valid.' in caplog.text


@responses.activate
def test_invalid_credentials(caplog):
    # takes in login and password
    # validates credentials with DNSExit
    # parse text response
    # log credentials are invalid
    # return false if credentials are invalid
    login = 'tester'
    password = 'Hello123'
    creds_validation_url = 'https://update.dnsexit.com/ipupdate/account_validate.jsp?login={}&password={}'.format(login, password)
    responses.add(responses.GET, creds_validation_url, body='\r\n\r\n\r\n\r\n\r\n\r\n\r\n1=Password Invalid\r\n')

    result = validate_credentials(login=login, password=password)
    assert not result
    assert 'The provided DNSExit IP Update credentials are not valid, exiting.' in caplog.text
