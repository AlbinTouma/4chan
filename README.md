# Signal Sifter

Sweden faces a growing crime epidemic, but public access to detailed crime data remains limited despite extensive media coverage, particularly of gang violence. This project seeks to address this gap by using alternative data sources to map crimes and their locations.

One of the key resources explored is Flashback Forum, a popular Swedish platform where users discuss crime, share media links, and analyze court cases. 

These discussions offer raw, unfiltered insights into incidents, potentially revealing patterns and connections not easily accessible through traditional sources.

## Method

![screenshot](./flashback.png)

Flashback is similar to 4Chan and other online forums. The site has boards that contain threads. Each thread is a discussion on a topic and contains posts and replies.

The idea is to crawl crime related boards and parses information from different threads. 
In the future the crawler will collect entities, events, and links from posts and replies on each thread.
