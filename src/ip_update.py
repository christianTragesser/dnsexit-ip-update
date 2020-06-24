import requests
import socket


def get_update_url(data_url):
    update_data = requests.get(data_url)
    update_domain = update_data.text.split()[0].split('//')[-1]
    return 'https://'+update_domain


def evaluate_ip_sync(domain):
    egress_ip = requests.get('https://api.ipify.org').text
    dns_ip = socket.gethostbyname(domain)
    if egress_ip == dns_ip:
        return True
    else:
        return False