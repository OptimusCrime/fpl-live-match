FROM nginx:1.25.4-alpine

ARG ENV=prod

COPY ./_docker/site.conf /etc/nginx/conf.d/site.conf
COPY ./frontend/dist /static

RUN chmod 644 /etc/nginx/conf.d/site.conf \
    && rm /etc/nginx/conf.d/default.conf
