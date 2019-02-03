# Edgar

A crawler to get company filing data from XBRL filings. The fetcher parses through the HTML pages and extracts data based on the XBRL tags that it finds and collects it into filing data arranged by filing date. 

# Goals
The primary requirement of any stock analysis is the data that the comapny files with the SEC. The filings, starting from 2010-11 required to be done using XBRL tags which required the filer to tag each data point with a tag. This made identifying and classifying the data a lot easier for machines.

The goal of this package is to get publically available filings on sec.gov and parse the filings to get data origanized into queriable data points. This package is being ddeveloped with the goal to provide an interface for other packages to get filing data for a company and use that data to create insights such as valuations and trends. 

# Design

The package is primarily organized as multiple interfaces into the provided functionality. The user is expected to make use of these interfaces to gather the data and query it as needed

# Interfaces 

All interfaces needed to use this package is defined in edgar.go and described below.

# FilingFetcher
This is the starting point for use of this package. The package is initialized with a fetcher. The user will use the fetcher interface to provide a ticker and filing type to startup a company folder. The user has an additional API in the interface to initialize a company folder with a saved folder. 

# CompanyFolder
A user will be given a company folder with the filings (retrieved ones) for every company (ticker). The user uses the folder to get any filing information related to that company. The filings are indexed internally based on filing type and the date of filing. When a user of the package requests a filing, the filing is looked up in the cache and if not available, will be retrieved from edgar and populated into the folder.

# Filing
Filing is an interface to get filing data related to a specific filing. The user uses this interface to extract required data. The Filing is retrieved from the company folder as needed. An error is returned if the data was unavailable.
 
