FROM node:16

# Create app directory
WORKDIR /usr/src/app

ENV PORT 11000
ENV ANOTHER_VAR another_value
ENV PORT 11001

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
COPY package*.json ./

RUN npm install

# Bundle app source
COPY . .

EXPOSE 8085
CMD [ "node", "server.js" ]
