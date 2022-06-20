curl --request GET \
  --url http://localhost:8080/users/2435/meetings \
  --header 'Content-Type: application/json' \
  --data '{
	"from": "2022-08-26T18:00",
	"to": "2022-11-21T19:00"
}'
