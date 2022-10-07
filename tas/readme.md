# throw against site
Get all archived URLs for a given prefix URL (thanks Wayback Machine).
Throw against current site.
Maybe uncover a something or other that was once crawlable and no longer is.

## command-line options
Usage of tas:
  -g int
    	Number of goroutines (default is 10).
  -json bool
    	Output results as json (default true).
  -t int
    	Request timeout (in milliseconds). Default is 5000.
  -txt bool
    	Output results as txt (default false)
  -u string
    	URL to get.

## notes
* don't be greedy -- pick either txt or json for your output (default json)