# Crime-Sifter

Sweden faces a growing crime epidemic, but public access to detailed crime data remains limited despite extensive media coverage, particularly of gang violence. Crime-Sifter seeks to address this gap by using alternative data sources to map crimes and their locations.

One of the key resources explored is Flashback Forum, a popular Swedish platform where users discuss crime, share media links, and analyze court cases. 

Discussions on Flashback offer raw, unfiltered insights into incidents, potentially revealing patterns and connections not easily accessible through traditional sources.

## Method


<image src="./flashback.png" style='margin: 10px;'>

Flashback is similar to forums like 4Chan. There are boards with different threads. Each thread is a discussion containing posts and replies. Our method collects threads and pass them through a locally hosted LLM that dentifies crime related data.

We create a corpus of data using the go colly scraper and preserve threads in a sqlite database.  Threads are passed through our locally hosted LLM. Our prompt asks the model to identify the crime mentioned in a thread and the LLM's response is stored in the crimes table. We use this table to filter for threads with a mentioned crime. These threads are slated for further analysis.

## Tables & Schema

Data collected while scraping boards include metadata such as link, number of replies, and views. We use this data to rank crimes with the assumption that threads with a higher count of views and replies are more likely to contain a crime. Data points at this stage of collection that are useful include:

Crime related data:

- Crime
- Location(s) mentioned
- Date(s) mentioned

Metadata:

- Public interest in crime (replies/views)
- Sentiment analysis



**Threads table**


|                     title                     |       date       |    author    |          link          | replies | views |
|-----------------------------------------------|------------------|--------------|------------------------|---------|-------|
| 24åring knivskuren i Lund 11 mars             | 2007-07-02 01:58 | malmomannen  | /t456961               | 19      | 2790  |
| Gruppvåldtäkt på 13-åring....                 | 2007-06-30 11:32 | Carcass      | /t508734               | 22      | 2762  |
| Kvinna rånad och dödad i Malmö                | 2007-06-29 22:11 | Trinnit      | /t506873               | 40      | 5006  |
| Stenkastning i Rinkeby mot polisen.           | 2007-06-29 20:13 | Kingfisher   | /t508300               | 9       | 1059  |
| bilbomb i centrala London                     | 2007-06-29 17:42 | Gefundenes F | /t508495               | 2       | 452   |
| Vem är dörrvakten?                            | 2007-06-29 14:14 | RedCircle    | /t508563               | 0       | 2159  |
| Skyddsgruppen                                 | 2007-06-28 17:53 | arbetarklass | /t507528               | 17      | 1946  |
| Vem vill spränga Allsången i luften?          | 2007-06-27 23:35 | Mr Magoo     | /t507197               | 12      | 1376  |
| Narkotikaliga på väg att sprängas i Västerås. | 2007-06-27 20:21 | Maeglin      | /t507672               | 3       | 1408  |
> 


Our locally hosted LLM parses the crimes. 

**Crimes table**

 promp                       |         crime         
-------------------------------------------------|-----------------------
♦·♦·♦ Läs detta innan du postar! ♦·♦·♦           |  No crime.            
24åring knivskuren i Lund 11 mars                |  Assault.             
Gruppvåldtäkt på 13-åring....                    |  Gruppvåldtäkt (Group 
Kvinna rånad och dödad i Malmö                   |  Infanticide.         
Stenkastning i Rinkeby mot polisen.              |  Arson.               
bilbomb i centrala London                        |  Bomb threat.         
Vem är dörrvakten?                               |  No crime.            
Skyddsgruppen                                    |  No crime.            
Vem vill spränga Allsången i luften?             |  Bomb threat.         
Narkotikaliga på väg att sprängas i Västerås.    |  Narcoterrorism.      
Vad kan det här ha varit för folk?               |  No crime.            
Mord i Malmö                                     |  Mord (Murder)        
Landsortskriminalitet                            |  No crime.            
Skottlossning i Luleå i natt. Mellan vilka gäng? |  Bomb threat.         
Mord i hallstahammar                             |  Mord (Murder)        
Poliser beskjutna i Nyköping                     |  Assault.             
Elitidrottare dömd för våldtäkt                  |  Assault.             
Samuel i Värnesborg                              |  No crime.            
Ny specialpolis jagar efterlysta                 |  No crime.            
Hur ta reda på om någon är häktad?               |  No crime.            



## Roadmap

So far the project has crawled one board and collected metadata on 66,000 threads posted between 2007-2005. Here are some possible future features:

- Rank threads by probability of crime related data based on thread title.
- Collect crime data from posts and replies on threads.
- Scrale links present in threads and enrich threads data.
- Geolocate incidents and pictures
- Plot incidents on an interactive map.

