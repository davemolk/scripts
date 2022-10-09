import argparse
from collections import defaultdict
import json
import random
import time

import requests

def create_url(url):    
    begin = "https://web.archive.org/web/timemap/json?url="
    mid = "&matchType=prefix&collapse=urlkey&output=json&fl=original%2Cmimetype%2Ctimestamp%2Cendtimestamp%2Cgroupcount%2Cuniqcount&filter=!statuscode%3A%5B45%5D..&limit=10000&_="
    milli = time.time() * 1000
    return f'{begin}{url}{mid}{milli}'

def make_request(url, timeout):
    u = create_url(url)
    headers = {
        'user-agent': get_user_agent(),
    }
    try:
        r = requests.get(u, headers=headers, timeout=timeout)
    except requests.exceptions.RequestException as e:
        raise SystemExit(e)

    return r

def response_to_json(response):
    return response.json()

def get_user_agent():
    agents = [
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4692.56 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4889.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/601.7.7 (KHTML, like Gecko) Version/9.1.2 Safari/601.7.7",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:99.0) Gecko/20100101 Firefox/99.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
    ]
    rando = random.randint(0, len(agents) - 1)
    return agents[rando]

def get_user_input():
    parser = argparse.ArgumentParser(description="submit url plz")
    parser.add_argument('-u', '--url', type=check_input, default="", help='provide full url')
    parser.add_argument('-j', '--json', type=bool, default=True, help='output results as json')
    parser.add_argument('-t', '--timeout', type=int, default=5000, help='timeout in ms')
    parser.add_argument('-txt', type=bool, default=False, help='output results as txt')
    args = parser.parse_args()
    return args.url, args.timeout, args.txt

def check_input(u):
    s = str(u)
    if s == "":
        raise argparse.ArgumentTypeError('must provide a URL')
    return s

def write_json(results):
    with open('results.json', 'w') as f:
            json.dump(results, f, indent=2, ensure_ascii = False)

def write_txt(name, results, key=None):
    if key is None:
        with open(name, 'w') as f:
            for k, v in results.items():
                f.write(f'{k}: {v}\n')
    else:
        with open(name, 'w') as f:
            for v in results[key]:
                f.write(f'{v}\n')


if __name__ == "__main__":
    url, timeout, txt = get_user_input()
    timeout = timeout / 1000
    response = make_request(url, timeout)
    data = response_to_json(response)

    results = defaultdict(list)
    with requests.Session() as s:
        for u in data[1:]:
            s.headers['User-Agent'] = get_user_agent()
            r = s.get(u[0], headers=s.headers, timeout=timeout)
            results[r.status_code].append(u[0])

    if not txt:
        write_json(results)
    else:
        write_txt("results.txt", results)
    for k in results.keys():
        write_txt(f'{k}.txt', results, k)