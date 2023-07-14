# Use Alpine Linux as the base image
FROM alpine:3.18

# INSTALL DEPENDENCIES
RUN apk add --no-cache \
	# ESSENTIALS
	lighttpd \
	php82-fpm \
	php82-mysqli \
	# PHP EXTRAS
	#performs remote request operations
	php82-curl \
	# validate Text Widget content
	php82-dom \
	# properly handle UTF8 text
	php82-mbstring \
	# permits SSL-based connections to other hosts.
	php82-openssl \
	# decompressing plugins, themes, update packages
	php82-zip \
	# high-performance object caching
	php82-pecl-memcached \
	# xml parsing
	php82-xml \
	# implement access to files servers via FTP
	php82-ftp

# CONFIGURE PHP
# enable php-fpm socket
COPY php-fpm.conf /etc/php82/php-fpm.d/www.conf
RUN mkdir /run/php

# CONFIGURE LIGHTTPD
COPY lighttpd.conf /etc/lighttpd/lighttpd.conf

# EXPOSE PORTS
EXPOSE 80

# START PHP AND NGINX
CMD php-fpm82 && lighttpd -D -f /etc/lighttpd/lighttpd.conf;