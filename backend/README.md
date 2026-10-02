# memebrary

A library for your memes

## Dev notes

```shell
# shell 1
./run-env.sh

# shell 2
./run-for-dev.sh

# shell 3
websocat ws://localhost:7070/api/__stream --exit-on-eof | jq

# shell 4
curl -X POST http://localhost:7070/api/memes -d '[{"filename": "blah", "mime_type": "blah", "original_name": "blah", "size": 1}]' | jq
curl -X POST http://localhost:7070/api/tags -d '[{"name": "foo"}]' | jq
curl -X POST http://localhost:7070/api/meme-tags -d "[{\"meme_id\": \"$(curl -s -X GET http://localhost:7070/api/memes | jq -r .objects[0].id)\", \"tag_id\": \"$(curl -s -X GET http://localhost:7070/api/tags | jq -r .objects[0].id)\"}]" | jq
```
