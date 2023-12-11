# Security Tools - Blazingly Fast Written in GO

### Features
1. Online Port Scanner: Scans 65535 ports
2. Fuzzing tool: Used to fuzz directory of website.

### How to use?
1. Clone the repository
2. Run command `go run main.go`
3. An http server will be stated on http://127.0.0.1:8080

## How to interact?

### Port
Method: GET
URL: http://127.0.0.1:8080/port?host=example.com

### Fuzzing
1. Method: POST
2. URL: http://127.0.0.1:8080/directories
3. Req Body: {url: http://example.com}
4. Wordlist File(optional): Send a wordlist file it is optional else it will use default wordlist. 

### Supported On
1. Linux
2. Windows
3. Termux

I need to highlight these ==very important words==
