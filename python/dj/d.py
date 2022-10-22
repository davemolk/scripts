import argparse
import json
import random
import time

import requests

def get_joke(url: str, headers: dict) -> str:
    headers["User-Agent"] = get_user_agent()
    try: 
        r = requests.get(url, headers=headers)
    except requests.exceptions.RequestException as e:
        raise SystemExit(e)
    return r.text

def get_user_agent() -> str:
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
    parser = argparse.ArgumentParser(description="joke me")
    parser.add_argument("-j", "--joke", type=str, help="get random joke")
    parser.add_argument("-t", "--term", type=str, default="", help="enter a search term or theme")
    args = parser.parse_args()
    return args.joke, args.term

def prep_format(joke: str):
    joke = joke.replace('"', "")
    if "?" in joke[:-1]:
        setup_punchline(joke, "?")
    elif "." in joke[:-1]:
        setup_punchline(joke, ".")
    else:
        print(joke)

def setup_punchline(text, splitter):
    split_text = text.split(splitter)
    if len(split_text) == 1:
        print(text)
    else:
        print(split_text[0] + splitter + "\n")
        time.sleep(2)
        if splitter == ".": 
            print(split_text[1].strip() + ".")
        else:
            print(split_text[1].strip())


if __name__ == "__main__":
    random_url = "https://icanhazdadjoke.com/"
    theme_url = "https://icanhazdadjoke.com/search?term="
    joke, term = get_user_input()
    if term != "":
        headers = {
            "Accept": "application/json",
        }
        url = f'{theme_url}{term}&limit=30'
        results = json.loads(get_joke(url, headers))
        if results == "":
            print("no jokes for that term")
        else:
            rando = random.randint(0, len(results["results"]) - 1)
            joke = results["results"][rando]["joke"]
            prep_format(joke)
    else:
        headers = {
            "Accept": "text/plain",
        }
        joke = get_joke(random_url, headers)
        prep_format(joke)