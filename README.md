# scripts
a collection of odds, ends, and everything in between (added below once stable). see individual directories for individual readmes (maybe?).

# go
## 403
mostly a port of https://github.com/yunemse48/403bypasser into Go

## fof
enter a search term and a file name (list of additional search terms, line-separated) and get back results (blurb plus url) from ask, bing, brave, duckduckgo, yahoo, and yandex. use 'exact' for exact term matching.

full version (renamed searcher) can be found [here](https://github.com/davemolk/searcher)!

## lines
take a bunch of files with line-separated urls/words/whatevers and get one large file containing all those urls/words/whatevers. Or, do the same thing with csv files.

## ngd
google dorking, but with duck duck go (so, not google dorking...)
(larger version under development)

full version (renamed dorking) can be found [here](http://github.com/davemolk/dorking)!

## pu
parse urls in style. read urls from a line-separated file or from stdin. Get domains, keys, a map of keys and values, paths, user info (username and password), or values.

full version (renamed urlbits) can be found [here](http://github.com/davemolk/urlbits)!

## pw
need a kinda sorta safe(ish?) password and too lazy to log into your pw manager or some rando genderating rando site? type something in and get a mostly gibberish password that contains at least one lower-case letter, one upper-case, one special character, and one number. there are, of course, much easier ways to do this -- I was mainly interested in practicing pipelines :)

## rp
replace a key=value parameter in a given URL(s) with your own key=value. Bonus! Pipe in an encoded mess and decode it.

## tas
throw against site (pull down archived links from Wayback Machine, run against the site, see what status codes currently are)

## wl
supply a url, get a wordlist (for that page only). adjust your results by filtering out terms, adding a minimum length, or requiring a certain instance count.

# python
## tas
throw against site (pull down archived links from Wayback Machine, run against the site, see what status codes currently are) (not concurrent...yet)

## dj
dad joke, ftw. use -t to enter a search term, otherwise enjoy a randomly selected dad joke