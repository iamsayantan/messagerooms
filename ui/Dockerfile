FROM node:12.14-alpine

# dependencies like git
RUN apk add --update \
  python \
  python-dev \
  py-pip \
  build-base \
  git \
  openssh-client \
&& pip install virtualenv \
&& rm -rf /var/cache/apk/*

# create destination directory
WORKDIR /usr/src/messagerooms-ui

# copy package.json and package-lock.json
COPY package*.json ./

RUN npm install

COPY . .

# set app serving to permissive / assigned
ENV NUXT_HOST=0.0.0.0
# set app port
ENV NUXT_PORT=3000

# building the nuxt app. this is a necessary step event if its only
# a static website. this builds the server as well.
RUN npm run build

# the nuxt app would be exposed on port 5050
EXPOSE 3000

# start the app
CMD ["npm", "start"]
