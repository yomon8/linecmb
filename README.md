# linecmb
Read multiple input streams and Write them into stdout 


## Usage
multiple file descripters passed in an argument, Detect EOL(\n), and Merge and Output to file

```
linecmb <(ssh server1 -C "tail -f /path/to/file") <(ssh server2 -C "tail -f /path/to/file") <(ssh server3 -C "tail -f /path/to/file")
```

## Install


```
go get github.com/yomon8/linecmb
go install github.com/yomon8/linecmb/...
```

or 
 
Download from [released file](https://github.com/yomon8/linecmb/releases)


## Licence

[MIT](https://github.com/yomon8/linecmb/blob/master/LICENSE)

## Author

[yomon8](https://github.com/yomon8)

