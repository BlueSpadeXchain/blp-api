### balancer
- connected to pyth to continuously recheck if order book needs rebalancing and places piority orders in such a case


### escrow
- controls onboarding and offboarding of non-active money
- stake assets cannot be withdrawn unless there is sufficent liquidity in the protocol


### info
- hello world test and version response
- current orderbook data
- current user data


### orders
create, modify, or close positions
protocol orders take priority


### taskqueue
- all the incoming orders likely cannot be handled at once thus we put them in a giant queue


### bindings
- generated for escrow.go?

### ws
- websocket to listen to contract interactions (deposits and stakes)

### db
- orders
  - perps tables for truth tables, user tables users cached data
  - user can prompt to resync data (refresh) from perp tables


### TODO
- [ ] create table structure
- [ ]
- [ ] create of position orders, assume large liquidity 100m BUSD, max order size 1k
  - [x] validate sigs (user or relay)
- [ ] modify orders
- [ ] close orders