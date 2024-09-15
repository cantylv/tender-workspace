easyjs:
	easyjson -no_std_marshalers -all internal/entity

run: 
	rm -rf services/postgres/data
	go mod vendor 
	docker compose up 