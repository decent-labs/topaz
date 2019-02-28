build-linux-api:
	docker build -t topaz-api . -f Dockerfile-build-linux-api
	docker run -v $$PWD:/opt/mount --rm --entrypoint cp topaz-api:latest topaz-api /opt/mount/build/topaz-api-linux

build-linux-batch:
	docker build -t topaz-batch . -f Dockerfile-build-linux-batch
	docker run -v $$PWD:/opt/mount --rm --entrypoint cp topaz-batch:latest topaz-batch /opt/mount/build/topaz-batch-linux

build-linux-migrate:
	docker build -t topaz-migrate . -f Dockerfile-build-linux-migrate
	docker run -v $$PWD:/opt/mount --rm --entrypoint cp topaz-migrate:latest topaz-migrate /opt/mount/build/topaz-migrate-linux

deploy-linux-api:
	rsync -azvv -e ssh build/topaz-api-linux ubuntu@sandbox.topaz.io:~
	ssh ubuntu@sandbox.topaz.io sudo systemctl restart topaz-api

deploy-linux-batch:
	rsync -azvv -e ssh build/topaz-batch-linux ubuntu@sandbox.topaz.io:~
	ssh ubuntu@sandbox.topaz.io sudo systemctl restart topaz-batch

deploy-linux-migrate:
	rsync -azvv -e ssh build/topaz-migrate-linux ubuntu@sandbox.topaz.io:~
	ssh ubuntu@sandbox.topaz.io ./topaz-migrate-linux