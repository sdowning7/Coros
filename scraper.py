from lxml import html
from urllib.parse import urlparse
from argparse import ArgumentParser
import requests

class Scraper():
    
    # url: points to an html document for this item to query and parse
    def __init__(self, url):
        self.url = url
        try:
            self.resp = requests.get(self.url)
            self.tree = html.fromstring(self.resp.content)
        except(Exception):
            raise ValueError("unable to get response from URL")
        

    # returns a list of all unique URLs from the html document pointed to by self.url
    def listURLs(self):
        URLs = [url for element, attribute, url, pos in self.tree.iterlinks()]
        uniqueURLs = list(set(URLs))
        return uniqueURLs
    
    # determines if the given url is external from the location of the parent url 
    # url: the URL to be checked to see if it is external
    # returns: True if the url is external, otherwise returns false
    def isExternal(self, url):
        return urlparse(url).netloc != ''

    # returns a list of unique External URLs from the html document pointed to by self.url
    def listExternalURLs(self):
        uniqueURLs = self.listURLs()
        externalURLs  = []

        for url in uniqueURLs:
            if self.isExternal(url):
                externalURLs.append(url)

        return externalURLs

if __name__ == "__main__":

    #set up argument parser
    parser = ArgumentParser()
    parser.add_argument('urls', nargs='*')

    #get our arguments
    args = parser.parse_args()
    urls = args.urls

    #scrape each url for external links
    for url in urls:
        try:
            s = Scraper(url)
            print(f'{url} {len(s.listExternalURLs())}')
        except ValueError:
            print(f'{url} -1')