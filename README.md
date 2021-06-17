# Fibonacci Demo

Simple demo of a fibonacci generator

### Dependencies
*Developed on a Mac, not tested on Linux but I would expect it to work too*
- Docker
- docker-compose (I realize it's baked into Docker now)
- Make

### Building
```
cd fibonacci
make
```
Point browser at http://localhost:8080/fib

### API
Simple GET API with query strings, most importantly the "type" query, which can be:
- "fibati"
- "getmemoized"
- "clear"

*Examples*
  - http://localhost:8080/fib?type=fibati&i=12 # get the 12th fibonacci number
  - http://localhost:8080/fib?type=getmemoized&value=120 # retrieves the count of memoized values less than or equal to 120
  - http://localhost:8080/fib?type=clear # clears the database
