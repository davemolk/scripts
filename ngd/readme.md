# ngd
google dorking, but for duck duck go (so...not google dorking?)

heads up — this is under active development in a new repo, so stay tuned and reach out if there's something you want to see!

edit: the full version (renamed dorking) can be found [here](http://github.com/davemolk/dorking)


### example usage
```
./ndg -qe "dave molk music" -no davemolkmusic.com
```

### flags
```
-q string
    search query
-qe string
    search query with exact matching
-ft string
    search for filetypes (the duck supports pdf, doc(x), xls(x), ppt(x), and html)
-inurl bool
    return search results with query in the url
-site string
    a site/domain to search
-nosite string
    a site/domain to exclude from search results
-t int
    request timeout (in milliseconds)
-intitle bool
    return search results with query in the title
```

#### a note
there are a few query limitations here, yes (e.g. something like intitle currently doesn't allow for its own query). they're all being updated in the soon(ish) to be released sparkly new version.