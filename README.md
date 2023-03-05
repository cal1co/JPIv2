# JPIv2

JPiv2 is a two-way Japanese-English API that allows users to search for roughly 280,000 Japanese dictionary entries from JMDict. It provides two endpoints: ***/generatekey*** and ***/search***. The ***/generatekey*** endpoint generates an API key for users to use other API functions, while the ***/search*** endpoint is a fuzzysearch method that indexes the JMDict entries.

## Endpoints 

### /generatekey
The ***/generatekey*** endpoint generates an API key that users can use to access other API functions.

**Request**
```
GET /generatekey
```
**Response**
```json 
{
  "apikey": "a1b2c3d4e5f6g7h8i9j10k11l12m13n14o15p16q"
}
```

### /search
The /search endpoint is a GET method that searches roughly 280,000 Japanese dictionary entries from JMDict. This method is a fuzzysearch and supports two-way indexing. To use the search function, queries are formatted like ***/search?query=word***.  

#### Example 1:
**Request**
```/search?query=鴨```

**Response**
```json
    {
        "Word": "鴨",
        "Alternate": "かも",
        "Freq": "611",
        "Def": [
            "duck"
        ],
        "Pitch": "かも｛鴨｝\n発音図：カ↓モ [1]\n助詞付：カ↓モオ\n"
    },
    {
        "Word": "鴨",
        "Alternate": "かも",
        "Freq": "610",
        "Def": [
            "easy mark",
            "sucker",
            "sitting duck"
        ],
        "Pitch": "かも｛鴨｝\n発音図：カ↓モ [1]\n助詞付：カ↓モオ\n"
    },
```

#### Example 2:
**Request**
```/search?query=duck```

**Response**
```json
[
    {
        "Word": "鴨",
        "Alternate": "かも",
        "Freq": "611",
        "Def": [
            "duck"
        ],
        "Pitch": "かも｛鴨｝\n発音図：カ↓モ [1]\n助詞付：カ↓モオ\n"
    },
    {
        "Word": "鶩",
        "Alternate": "あひる",
        "Freq": "6",
        "Def": [
            "domestic duck"
        ],
        "Pitch": ""
    },
]
```

NOTE: Spaces are allowed when searching.

## Technologies Used
JPiv2 is written in Go and uses Mux, a rate limiter, API keys, and OpenSearch.

### Data
All queries include a kanji reading, hiragana reading, English definition, and pitch accent information. Eventually, audio pronunciation will be included for each entry.

### Rate Limiting
JPiv2 includes rate limiting to prevent abuse of the API. Each API key is limited to a certain number of requests per minute.

## Conclusion
JPiv2 is a powerful two-way Japanese-English API that provides users with access to roughly 280,000 Japanese dictionary entries from JMDict. With a simple API key generation endpoint and easy-to-use search functionality, JPiv2 is the perfect tool for anyone looking to incorporate Japanese language support into their applications.
