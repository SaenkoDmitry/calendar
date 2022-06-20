curl --request POST \
  --url http://localhost:8080/users \
  --header 'Content-Type: application/json' \
  --data '{
	"first_name": "Ilon",
	"second_name": "Mask",
	"email": "ilon.mask@mail.ru",
	"zone": "Europe/Moscow"
}'
