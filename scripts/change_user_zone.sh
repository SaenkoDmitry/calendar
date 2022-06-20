curl --request PUT \
  --url http://localhost:8080/users/1 \
  --header 'Content-Type: application/json' \
  --data '{
	"zone": "Europe/Copenhagen"
}'
