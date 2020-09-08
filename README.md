[![Build Status](https://travis-ci.com/viveknangal/awsjack.svg?branch=master)](https://travis-ci.com/viveknangal/awsjack)
![GitHub](https://img.shields.io/github/license/viveknangal/awsjack?style=plastic)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/viveknangal/awsjack?style=plastic)
![picture](static/images/awsjack.png)
## How does it work
![picture](static/images/aws-jack.gif)

## Run `awsjack` using Docker as below:-

- Run `docker build` command
```
docker build -t awsjack:latest .
```
- Run `docker run` command
   ```
docker run -v $(echo ~/.aws):/root/.aws -p 8080:8080 awsjack:latest
```
- Access application via below endpoint 
```
     http://localhost:8080/
```
## License
Licensed under the [MIT License](LICENSE)
