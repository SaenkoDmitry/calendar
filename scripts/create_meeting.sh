curl --request POST \
  --url http://localhost:8080/meetings \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Important meeting",
	"organizer_id": 1,
	"participants": [2,3,4],
	"from": "2022-01-02T18:00",
	"to": "2022-01-02T19:00"
}'
