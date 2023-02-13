## Hello! This is my submission for the Forbes Take-Home Test for the role of Senior Software Engineer on the Systems Team

#### STACK DECISION:
After reviewing the technical requirements of this challenge, I decided to develop this in Go for a few reasons. The
primary reason being that to do an accurate comparison of closest words to a given word would require the levenshtein
algorithm, and the comparison itself would be comparable to an array lookup with significantly heavier lift on comparison.
That combined with the fact that the number of full array lookup comparisons would be equal to the number of items that
were mispelled, it presented the perfect opportunity to introduce some GO concurrency.

#### LOCAL DEVELOPMENT STACK:
This project was developed on a Golang stack on Windows 10 with Gin as a webserver, the command to run is:
go run main.go (and any development or production qualifiers you want to add)

#### OPTIMIZATION THOUGHT:
This project doesn't utilize a database currently but I believe that would be much more efficient with a growing
dictionary size because there are several databases that include very efficient levenshtein comparisons (there are also
several in golang, python, and java that would be available on a true production version of this application)

#### DEPLOYMENT PLAN:
With Golang being created by Google and used by a lot of the development community today, the deployment of this
webserver would be relatively non-complicated. After submitting this on 02/13/2023 I will be hosting it on my personal
website at forbes.danielstone.dev behind an nginx proxy forwarding calls to the local Gin Golang webserver running on a
Ubuntu Gnome 22.02 server. This server is hosted on an AWS S3 instance, has domain hosting through Namecheap, and has
nameserver settings configured with Route 53 through AWS and instanced on an AWS EC2 instance. Serving the webserver
would include setting up a systemd service that is set to always restart and points to the production compiled binary
version of the Gin Golang webserver.

API Spec utilized for construction of this project:
```
openapi: 3.1.0
info:
  title: API Team Take Home
  contact: {}
  version: '1.0'
jsonSchemaDialect: https://json-schema.org/draft/2020-12/schema
servers:
- url: http://localhost:8080
  variables: {}
paths:
  /dictionary:
    get:
      tags:
      - Misc
      summary: Fetech All Dicionary
      description: This endpoint should return a list of all items in the dictionary
      operationId: FetechAllDicionary
      parameters: []
      responses:
        '200':
          description: OK
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/200OK'
              example:
                dictionary:
                  list:
                  - '20'
                  - a
                  - and
                  - apple
                  - are
                  - believes
                  - big
                  - biggest
                  - brand
                  - brands
                  - business
                  - capitalism
                  - celebrate
                  - change
                  - committed
                  - company
                  - consistently
                  - conversations
                  - culture
                  - drive
                  - entrepreneurial
        '500':
          description: Internal Server Error
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error1'
              example:
                error: Unable to reteive data
      deprecated: false
    post:
      tags:
      - Misc
      summary: Add Dictionary Entry
      description: This endpoint should be used to add words to the dictionary
      operationId: AddDictionaryEntry
      parameters: []
      requestBody:
        description: ''
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AddDictionaryEntryRequest'
            example:
              dictionary:
                add:
                - example
        required: true
      responses:
        '202':
          description: Accepted
          headers: {}
          content: {}
        '200':
          description: OK
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/DuplicateEntry'
              example:
                dictionary:
                  add:
                  - example
                error: duplicate entry
        '400':
          description: Bad Request
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BadRequest1'
              example:
                error: bad request
      deprecated: false
    delete:
      tags:
      - Misc
      summary: Delete Dictionary Entry
      description: This endpoint should be used to remove words from the dictionary
      operationId: DeleteDictionaryEntry
      parameters: []
      requestBody:
        description: ''
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DeleteDictionaryEntryRequest'
            example:
              dictionary:
                remove:
                - example
        required: true
      responses:
        '202':
          description: Accepted
          headers: {}
          content: {}
        '400':
          description: Bad Request
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BadRequest1'
              example:
                error: bad request
        '404':
          description: Not Found
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/NotFound1'
              example:
                dictionary:
                  add:
                  - example
                error: entry not found
      deprecated: false
    parameters: []
  /story:
    post:
      tags:
      - Misc
      summary: Story
      description: Using the posted body (story), the API should return a list of words that are not found or added to the dictionary of the system. For each word provide what the closest match for it would be from the dictionary. If more than one item is the same level of close match to a misspelled entity return the suggestion from the dictionary that comes first.
      operationId: Story
      parameters: []
      requestBody:
        description: ''
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/StoryRequest'
            example:
              story: Forbes belives in the power of entreprenuerial capitalizm and uses it various platforms to ignight
        required: true
      responses:
        '200':
          description: OK
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GoodRequest'
              example:
                story: Forbes belives in the power of entreprenuerial capitalizm
                unmachedWords:
                - word: belives
                  closeMatch: belives
                - word: entreprenuerial
                  closeMatch: entrepreneurial
                - word: capitalizm
                  closeMatch: capitalism
        '400':
          description: Bad Request
          headers: {}
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BadRequest1'
              example:
                error: bad request
      deprecated: false
    parameters: []
components:
  schemas:
    200OK:
      title: 200OK
      required:
      - dictionary
      type: object
      properties:
        dictionary:
          $ref: '#/components/schemas/Dictionary'
      examples:
      - dictionary:
          list:
          - '20'
          - a
          - and
          - apple
          - are
          - believes
          - big
          - biggest
          - brand
          - brands
          - business
          - capitalism
          - celebrate
          - change
          - committed
          - company
          - consistently
          - conversations
          - culture
          - drive
          - entrepreneurial
    Dictionary:
      title: Dictionary
      required:
      - list
      type: object
      properties:
        list:
          type: array
          items:
            type: string
          description: ''
      examples:
      - list:
        - '20'
        - a
        - and
        - apple
        - are
        - believes
        - big
        - biggest
        - brand
        - brands
        - business
        - capitalism
        - celebrate
        - change
        - committed
        - company
        - consistently
        - conversations
        - culture
        - drive
        - entrepreneurial
    Error1:
      title: Error1
      required:
      - error
      type: object
      properties:
        error:
          type: string
      examples:
      - error: Unable to reteive data
    AddDictionaryEntryRequest:
      title: AddDictionaryEntryRequest
      required:
      - dictionary
      type: object
      properties:
        dictionary:
          $ref: '#/components/schemas/Dictionary1'
      examples:
      - dictionary:
          add:
          - example
    Dictionary1:
      title: Dictionary1
      required:
      - add
      type: object
      properties:
        add:
          type: array
          items:
            type: string
          description: ''
      examples:
      - add:
        - example
    DuplicateEntry:
      title: DuplicateEntry
      required:
      - dictionary
      - error
      type: object
      properties:
        dictionary:
          $ref: '#/components/schemas/Dictionary1'
        error:
          type: string
      examples:
      - dictionary:
          add:
          - example
        error: duplicate entry
    BadRequest1:
      title: BadRequest1
      required:
      - error
      type: object
      properties:
        error:
          type: string
      examples:
      - error: bad request
    DeleteDictionaryEntryRequest:
      title: DeleteDictionaryEntryRequest
      required:
      - dictionary
      type: object
      properties:
        dictionary:
          $ref: '#/components/schemas/Dictionary3'
      examples:
      - dictionary:
          remove:
          - example
    Dictionary3:
      title: Dictionary3
      required:
      - remove
      type: object
      properties:
        remove:
          type: array
          items:
            type: string
          description: ''
      examples:
      - remove:
        - example
    NotFound1:
      title: NotFound1
      required:
      - dictionary
      - error
      type: object
      properties:
        dictionary:
          $ref: '#/components/schemas/Dictionary1'
        error:
          type: string
      examples:
      - dictionary:
          add:
          - example
        error: entry not found
    StoryRequest:
      title: StoryRequest
      required:
      - story
      type: object
      properties:
        story:
          type: string
      examples:
      - story: Forbes belives in the power of entreprenuerial capitalizm and uses it various platforms to ignight
    GoodRequest:
      title: GoodRequest
      required:
      - story
      - unmachedWords
      type: object
      properties:
        story:
          type: string
        unmachedWords:
          type: array
          items:
            $ref: '#/components/schemas/UnmachedWord'
          description: ''
      examples:
      - story: Forbes belives in the power of entreprenuerial capitalizm
        unmachedWords:
        - word: belives
          closeMatch: belives
        - word: entreprenuerial
          closeMatch: entrepreneurial
        - word: capitalizm
          closeMatch: capitalism
    UnmachedWord:
      title: UnmachedWord
      required:
      - word
      - closeMatch
      type: object
      properties:
        word:
          type: string
        closeMatch:
          type: string
      examples:
      - word: belives
        closeMatch: belives
tags:
- name: Misc
  description: ''
```