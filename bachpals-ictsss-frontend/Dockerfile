# Build
FROM node:alpine AS build-env

ARG REACT_APP_AUTH_PATH
ARG REACT_APP_API_PATH

ADD . /app/src

WORKDIR /app/src

RUN yarn

RUN REACT_APP_AUTH_PATH=${REACT_APP_AUTH_PATH} \
    REACT_APP_API_PATH=${REACT_APP_API_PATH} \
    yarn build

FROM nginx:alpine

ADD ./nginx.conf /etc/nginx/nginx.conf

WORKDIR /app

COPY --from=build-env /app/src/build .

RUN rm -rf src
