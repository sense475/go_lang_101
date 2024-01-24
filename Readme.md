
Hello this is how to set up this
install go ให้ได้
```
go mod init project_name
```

mongoDB

```
# This is used for local development
version: '3.8'
services:
  mongodb:
    image: mongo:latest # use the latest image.
    container_name: mongodb-go-sense
    restart: unless-stopped
    environment: # set required env variables to access mongo
      MONGO_INITDB_ROOT_USERNAME: username
      MONGO_INITDB_ROOT_PASSWORD: password

    ports:
      - 27019:27017
    volumes: # optional to preserve database after container is deleted.
      - ./database-data:/data/db
    networks:
      - customer-network

networks:
  customer-network:
    driver: bridge
    name: customer-network

```
```
go get -u github.com/gin-gonic/gin
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
```

สร้างไฟล์ main ของโปรเจค 

```
touch main.go
```



##How to run

```
go run main.go
```