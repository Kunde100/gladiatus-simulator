FROM php:8.2-fpm

RUN apt-get update && apt-get install -y \
	nginx

RUN rm /etc/nginx/sites-enabled/default

COPY nginx.conf /etc/nginx/sites-available/default

RUN ln -s /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default

COPY api/ /var/www/html/

RUN chown -R www-data:www-data /var/www/html

COPY start.sh /usr/bin/start.sh

RUN chmod +x /usr/bin/start.sh

EXPOSE 8080

CMD ["/usr/bin/start.sh"]
