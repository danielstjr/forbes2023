## Simple Levenshtein Golang project

#### DESCRIPTION:
I learned the Gin web framework over a weekend in February of 2023 and put together this little web app in Golang. The premise for this is pretty simple, maintain a textfile dictionary to audit stories for incorrect spellings. If a word doesn't exist in the dictionary, levensthein against the words remaining in the dictionary to find its closest match. This problem was particularly fun to solve because it represents O(M*N) complexity for finding the closest word match where M is the number of words not in the dictionary and N is the list of total items in the dictionary (when run synchronously). To get around this, I asynchronously calculate the levenshtein distance of words in the dictionary using Golang's thread management system. Each word gets its own thread and immutable dictionary instance to compare against. After spwaning all threads I wait for all of them to come back then return a list of closest matches (defering to alphateical order in case of identicial levenshtein distances).

#### ENDPOINTS:
There are two urls available with multiple routes attached:
- GET /dictionary
  - list all items in the dictionary currently
- POST /dictionary
  - adds an item(s) to the dictionary
- DELETE /dictionary
  - deletes an item(s) from the dictionary
  
- POST /story
  - Audit story for mispellings, provide closest suggestion by levenshtein distance of words in dictionary, defering to alphabetical order in the case of multiple words having an identical levenshtein distance
  
#### IMPROVEMENTS:
|The improvements of this project follows adding a database to the dictionary to allow for a larger dictionary list without concurrent file access issues, and adding a NoSQL cache to store levenshtein calculations in the cache based on the current dictionary list. This itself comes with one additional problem in concurrency: what happens if two transactions happen at the same time, one that audits a story and one that edits the dictionary? To solve that problem I would add in a dictionary_transactions table that has a primary id (known forward as version_id), json request body, and timestamp. This version id is then prepended to the NoSQL cache key, and each story audit request pulls its own most recent transaction id at the time of request to gauruntee a consistent levenshtein distance calculation across all words. Each levenshtein call inputs its own results into the cache in the form of {version_number:word} if a calculation is necessary. Each request would also increment/decrement the dictionary\_{version_number}\_requests key in the cache. At the end of the story request it would check if the latest dictionary transaction is at a later version number than the one it loaded at the beginning of the request, if it is and the number of requests on its original version is 0, it deletes all keys of the pattern {version_number:*} and the dictionary\_{version\_number}\_requests key.
