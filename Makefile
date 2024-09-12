easyjs:
	easyjson -no_std_marshalers -all internal/entity

run: 
	rm -rf services/postgres/data
	docker rm tender-platform
	docker rmi tender-workspace-tender-platform
	go mod vendor 
	docker compose up 