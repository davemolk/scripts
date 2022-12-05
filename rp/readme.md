# replace params
Replace a key=value parameter in a given URL(s) with your own key=value. Bonus! Pipe in an encoded mess and decode it.

## Flags
```
-q: query parameter as key=value 
-d: decode mode 
```

## Test URLs
```
https://example.com/foo/bar?first=1&second=2
https://example.com/bar?foo=me&baz=barz
https://www.google.com/search?name=golang&language=en&mascot=gopher&foo=bar
google.com
```

## Usage
(replace a value)
```
cat urls.txt | rp -v foo=bar
https://example.com/foo/bar?first=1&foo=bar&second=2 // adds &foo=bar
https://example.com/bar?baz=barz&foo=bar // overwrites value for preexisting key
https://www.example.com/rp?foo=bar&language=en&mascot=gopher&name=golang // kv already exists
2022/12/05 05:52:44 parse "google.com": invalid URI for request // no thanks
```

(decode query string(s))
```
echo "cart=%5B%5B0%2C+%7B%22logo%22%3A+%22kitten.jpg%22%2C+%22price%22%3A+0%2C+%22name%22%3A+%22Kitten%22%2C+%22desc%22%3A+%228%5C%22x10%5C%22+color+glossy+photograph+of+a+kitten.%22%7D%5D%2C+%5B0%2C+%7B%22logo%22%3A+%22kitten.jpg%22%2C+%22price%22%3A+0%2C+%22name%22%3A+%22Kitten%22%2C+%22desc%22%3A+%228%5C%22x10%5C%22+color+glossy+photograph+of+a+kitten.%22%7D%5D%2C+%5B1%2C+%7B%22logo%22%3A+%22puppy.jpg%22%2C+%22price%22%3A+0%2C+%22name%22%3A+%22Puppy%22%2C+%22desc%22%3A+%228%5C%22x10%5C%22+color+glossy+photograph+of+a+puppy.%22%7D%5D%5D" | rp -d

cart=[[0, {"logo": "kitten.jpg", "price": 0, "name": "Kitten", "desc": "8\"x10\" color glossy photograph of a kitten."}], [0, {"logo": "kitten.jpg", "price": 0, "name": "Kitten", "desc": "8\"x10\" color glossy photograph of a kitten."}], [1, {"logo": "puppy.jpg", "price": 0, "name": "Puppy", "desc": "8\"x10\" color glossy photograph of a puppy."}]]
```