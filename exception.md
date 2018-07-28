### Errors
```
fabric-samples/chaincode/tianchi-cvchain/vendor/github.com/miekg/pkcs11
vendor/github.com/miekg/pkcs11/pkcs11.go:26:18: fatal error: ltdl.h: No such file or directory
 #include <ltdl.h>
                  ^
compilation terminated.
```

Miss libtool for os,

if you are using centos/RHEL 7 you can use this code
```
yum install libtool-ltdl-devel
```
For Mac,
```
brew install libtool
```