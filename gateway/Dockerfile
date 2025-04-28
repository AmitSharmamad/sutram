FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY conf.d /etc/nginx/conf.d

EXPOSE 8010

CMD ["nginx", "-g", "daemon off;"]
