PROJECT_ROOT = "C:/Users/NaveenKumarasinghe/workspace/php-pods"
DB_ROOT_PASSWORD = "rootpass"

# setup
install:
	# create network
	docker network create ppp

	# create database container
	docker create --name ppp-db \
				  --network ppp \
				  --env MARIADB_USER=user \
				  --env MARIADB_PASSWORD=pass \
				  --env MARIADB_ROOT_PASSWORD=${DB_ROOT_PASSWORD} \
				  --publish 127.0.0.1:3306:3306 \
				  --volume ${PROJECT_ROOT}/containers/db/data:/var/lib/mysql \
				  mariadb:lts

	# create wordpress runtime image
	docker build -t ppp-wp ./containers/wp

start:
	# start database container
	docker start ppp-db

stop:
	# kill all containers
	-docker kill `docker ps -q`

clean: stop
	# remove database container & image
	-docker rm ppp-db
	-docker rmi ppp-db

	# remove application image
	-docker rmi ppp-wp

	# remove network
	-docker network rm ppp
	@echo "containers/db/data and sites were intentionally retained"

# shell
shell-wp:
	docker exec -it ppp-wp bash

shell-db:
	docker exec -it ppp-db mariadb -uuser -ppass

# site management
create-site:
	# create site database
	docker exec -it ppp-db mariadb -uroot -p${DB_ROOT_PASSWORD} -e "\
		CREATE DATABASE ${name}; \
		CREATE USER '${name}'@'%' IDENTIFIED BY '${password}'; \
		GRANT ALL PRIVILEGES ON ${name}.* TO '${name}'@'%'; \
		FLUSH PRIVILEGES; \
	"

	# create site directory
	mkdir -p "sites/${name}"
	tar -xzf distros/wordpress-vanilla.tar.gz -C "sites/${name}"
	chmod -R 755 sites/${name}
	mv sites/${name}/wordpress/* sites/${name}
	rm -rf sites/${name}/wordpress

	# generate wp-config.php file
	db_name="${name}" \
	db_user="${name}" \
	db_password="${password}" \
	db_host="wphub-db" \
	envsubst '$$db_name,$$db_user,$$db_password,$$db_host' < "containers/wp/wp-config.php.tmpl" > "sites/${name}/wp-config.php"

delete-site:
	# delete site database
	-docker exec -it ppp-db mariadb -p${DB_ROOT_PASSWORD} -e "\
		DROP DATABASE ${name}; \
		DROP USER '${name}'@'%'; \
	"
	-rm -rf sites/${name}