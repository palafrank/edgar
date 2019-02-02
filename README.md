# Edgar

A crawler to get company filing data from XBRL filings. The fetcher parses through the HTML pages and extracts data based on the XBRL tags that it finds and collects it into filing data arranged by filing date. 

# FilingFetcher
This is the main package handle used by the user. The package is initialized with a fetcher. The user will use the fetcher interface to provide a ticker and filing type to startup a company folder. The user has an additional API in the interface to initialize a company folder with a saved folder based on a writer that was saved at an earlier point.

# CompanyFolder
A user will be given a company folder with the filings (retrieved ones) for every company (ticker). The user uses the folder to get any filing information related to that company. The filings are indexed internally based on filing type and the date of filing. When a user of the package requests a filing, the filing is looked up in the cache and if not available, will be retrieved from edgar.

# Filing
Filing is an interface to get filing data related to a specific filing. The user uses this interface to extract required data. The Filing is retrieved from the company folder as needed. An error is returned if the data was unavailable.

All attempt is made to generate essential data that was not found/tagged in the filing. 
