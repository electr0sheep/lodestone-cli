# A Lodestone CLI
This project started with me wanting to sync data to ffxivcollect.com.
However, as I was writing, it became clear that I was essentially writing a CLI
for Lodestone.
## Prerequisties
[GO](https://golang.org/)
## Install
```
go get github.com/electr0sheep/lodestone-cli
```
### Usage
```
lodestone-cli help
```

### How does it work?
To read private data from lodestone, you will have to log in and copy the authentication cookie
![Alt text](/screenshots/lodestone.png?raw=true "Optional Title")
To sync data to ffxivcollect.com, you will again have to log in to the website and copy the authentication cookie
![Alt text](/screenshots/ffxivcollect.png?raw=true "Optional Title")