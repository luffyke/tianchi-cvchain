### Errors
1. Missing libtool
```
fabric-samples/chaincode/tianchi-cvchain/vendor/github.com/miekg/pkcs11
vendor/github.com/miekg/pkcs11/pkcs11.go:26:18: fatal error: ltdl.h: No such file or directory
 #include <ltdl.h>
                  ^
compilation terminated.
```

If you are using centos/RHEL 7, you can use below command.
```
yum install libtool-ltdl-devel
```
For Mac, use below command.
```
brew install libtool
```

2. zip with MACOS folder
```
Chaincode package failed Error getting chaincode code chaincode: Error getting chaincode package bytes: Error obtaining imports: <go, [list -f {{ join .Imports n}} github.com/hyperledger/fabric/examples/chaincode/go/cvChain/__MACOSX]>: failed with error: exit status 1">
```

```
zip -d cvChain.zip __MACOSX/\*
zip -d cvChain.zip \*/.DS_Store
```