import requests
import socket


def get_update_url(data_url):
    update_data = requests.get(data_url)
    update_domain = update_data.text.split()[0].split('//')[-1]
    return 'https://'+update_domain


def evaluate_ip_sync(domain):
    return False