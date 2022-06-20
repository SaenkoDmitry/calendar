curl --request POST \
  --url http://localhost:8080/meetings/suggest \
  --header 'Content-Type: application/json' \
  --data '{
	"participants": [2438, 2440, 2441, 2442],
	"min_duration_in_minutes": 60
}'
