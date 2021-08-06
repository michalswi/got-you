\# to be **installed**:
- make
- upx

\# **make** your friend

```
[pts1]$ make build-linux
[pts1]$ tree -ah bin/
bin/
├── [1.8K]  checker.go
└── [1.7M]  msconfig
[pts1]$ make web

[pts2]$ firefox http://localhost:8080/x/

[pts2]$ make wget
[pts2]$ ./dw
[pts2]$ make get

[pts2]$ make remove-bins
```
