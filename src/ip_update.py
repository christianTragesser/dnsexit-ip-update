import requests

update_data_url = 'https://www.dnsexit.com/ipupdate/dyndata.txt'

def get_update_url(data_url):
    update_data = requests.get(data_url)
    update_domain = update_data.text.split()[0].split('//')[-1]
    return 'https://'+update_domain
