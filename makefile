build-linux:
	docker build -t topaz-linux . -f Dockerfile-build-linux
	docker run -v $$PWD:/opt/mount --rm --entrypoint cp topaz-linux:latest topaz-api /opt/mount/build/topaz-api-linux
	docker run -v $$PWD:/opt/mount --rm --entrypoint cp topaz-linux:latest topaz-batch /opt/mount/build/topaz-batch-linux
	docker run -v $$PWD:/opt/mount --rm --entrypoint cp topaz-linux:latest topaz-migrate /opt/mount/build/topaz-migrate-linux
	rsync -azvv -e ssh build/topaz-api-linux ubuntu@sandbox.topaz.io:~
	rsync -azvv -e ssh build/topaz-batch-linux ubuntu@sandbox.topaz.io:~
	rsync -azvv -e ssh build/topaz-migrate-linux ubuntu@sandbox.topaz.io:~