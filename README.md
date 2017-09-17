# linecmb
Read multiple input streams and Write them into stdout 


## Usage
multiple file descripters passed in an argument, Detect EOL(\n), and Merge and Output to file

```
linecmb <(ssh server1 -C "tail -f /path/to/file") <(ssh server2 -C "tail -f /path/to/file") <(ssh server3 -C "tail -f /path/to/file")
```

## Description & Example

if you want to combine these 3 simple streams.

```
1 (sleep 1;echo SlowSlow)
2 (sleep 0;for i in $(seq 1 25);do printf aaa;echo;done) 
3 (sleep 0;for i in $(seq 1 25);do printf bbb;echo;done)
```

Expected output 

```
aaa
bbb
bbb
aaa
...
...
bbb
aaa
SlowSlow
```


handle this with simple bash like this.

```
((sleep 1;echo SlowSlow) & \
 (sleep 0;for i in $(seq 1 25);do printf aaa;echo;done) & \
 (sleep 0;for i in $(seq 1 25);do printf bbb;echo;done))
```

A line `a` and `b` and `\n` mixed is printed.

```sh
aaa
baaa
aaabaaa
\n
aaabaaa
\n
baaaaaa
\n
baaa
...
...
aaaSlowSlow
```


if you use `cat`.

```
cat <(sleep 1;echo SlowSlow) \
 <(sleep 0;for i in $(seq 1 25);do printf aaa;echo;done) \
 <(sleep 0;for i in $(seq 1 25);do printf bbb/;echo;done)  
```

```sh
SlowSlow # 1sec delay
aaa      # blocked. wait for 1 sec
aaa
...
...
bbb      # blocked. wait for all aaa outputed
bbb
```


if you use `paste`.

```
paste -d \\n <(sleep 1;echo SlowSlow) \
 <(sleep 0;for i in $(seq 1 25);do printf aaa;echo;done) \
 <(sleep 0;for i in $(seq 1 25);do printf bbb;echo;done)  
```

```sh 
SlowSlow # 1sec delay
aaa      # blocked. wait for 1 sec. aaa\b and bbb\n are outputed one by one.
bbb
aaa
bbb
...
...
aaa
bbb
```



you can get nonblocked and expected output with linecmb.

```
linecmb <(sleep 1;echo SlowSlow) \
 <(sleep 0;for i in $(seq 1 25);do printf aaa;echo;done) \
 <(sleep 0;for i in $(seq 1 25);do printf bbb;echo;done)  
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

