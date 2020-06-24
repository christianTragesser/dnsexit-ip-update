import pytest
import responses
import mock
import os, sys
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from ip_update import get_update_url

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
    responses.add(responses.GET, update_data, status=200)

    update_addr = get_update_url(update_data_url)
    assert update_addr == 'update.test.local/TestUpdate.sv'
