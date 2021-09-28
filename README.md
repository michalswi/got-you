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
[pts2]$ chmod +x ~/Downloads/dw
[pts2]$ ./dw
OR
[pts2]$ make wget
[pts2]$ ./dw

[pts2]$ make get
{
  "data": [
    {
      "id": 1,
      "remoteip": "127.0.0.1:43346",
      "hostname": "dummy",
      "localip": "192.168.1.123",
      "publicip": "12.34.56.78",
      "nmap": [
        22
      ],
      "os": "linux"
    }
  ]
}

[pts2]$ make remove-bins
```
