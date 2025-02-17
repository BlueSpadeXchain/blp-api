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
- [x] fix api structure
- [x] create table structure
  - [x] structure
  - [x] sql calls
- [x] error messages
  - [x] catch errors with backend
  - [x] create panic crash redudancy and checker
  - [x] create internal error wrapper
  - [x] create error sql table generation call
- [ ] order
  - [x] order validation with price feed and balance
  - [x] order math
    - [x] test
  - [x] create order
  - [ ] modify order
  - [x] group merge orders
  - [x] limit order
  - [x] stop loss enabled
  - [x] take profit enabled
- [x] create listener for contract changes
  - [x] create evm contract to be listened to
  - [x] create listener application
    - [x] upon message received call the api to notice log/db log
- [x] create of position orders, assume large liquidity 100m BUSD, max order size 1k
  - [x] validate sigs (user or relay)
  - [x] test liquidation order
- [x] add jwt/auth to supabase queries
- [ ] rebalancer
  - [x] connect to price feed for latest data
  - [ ] rebalancer hosting
    - [x] launch VM
    - [x] activate
    - [x] use railway
    - [ ] shutdown oracle vm
    - [ ] deploy ci/cd for VM
  - [x] get db and initialize liquidations call to api
- [x] organize sql queries
- [x] fix sql files for sqllite testing
- [x] add global balance tracking
- [x] scheduled distributions
- [x] deposit/stake change to use price feed latest price when depositing
- [x] enable fees
  - [x] from opening
  - [x] from closing
  - [x] from liquidations
  - [x] utilization fee
- [ ] api documentation
- [ ] readme sysdoc
- [ ] order queue
- [x] loadbalancer
- [x] ts example script
- [x] poke test db cron job
- [ ] withdrawls
  - [x] table log
  - [ ] kill switch
  - [x] backend functions
  - [x] db functions
- [ ] staking
  - [x] blu distribution
  - [x] "mint" usdc
  - [x] backend function
  - [x] db function
  - [ ] frontend staking
  - [ ] staking metrics (APR etc)
  - [x] seperate staking for blu (token) and blp (protocol liquidity)
  - [x] blp buying both contract and frontend
  - [ ] get distributions history
  - [ ] get user distributions history
- [ ] metrics
  - [x] test snapshot history
  - [ ] gui for historical data
  - [x] create cron job to snapshot data (hourly or daily)
  - [x] create snapshot api
  - [ ] test metrics
  - [x] emphsis of fees and volume
  - [ ] daily active users: frontend telemetry, wallets per order, connected wallets, number of orders
  - [ ] pair specific volume, liquidity, and volume 
- [x] change liquidations from 100% to 95%
  - [x] change math for liqPrice calcualtion
- [x] read replicas for all get calls (especially for rebalancer)
- [x] stop loss price, has no indication of value, this should be precalculated
- [x] handleResponse needs to be revisioned to use utils.Error instead of error
- [x] change fee to be based off of leverage, instead of collateral
- [x] add minimumum order side of 10 usd (handle with backend and show error)
- [ ] log order change history events
- [x] free up any users to be able to use 1250x leverage
- [ ] mainnet hosting/testing
  - [ ] required beta build and prod builds, ontop of the dev build
  - [ ] table refresh for mainnet/beta
- [x] cleanup rebalancer calls
- [x] add liquidation equation function to rebalancer calls
- [ ] solve viem ecdsa issue
- [ ] user data re-sync
  - [ ] db function
  - [ ] api function
  - [ ] test script (ts)
- [ ] local test enviornment auto-init script
- [ ] 10x max leverage on specific position
- [ ] prevent ordrs max payout to exceed global liquidity
  - [ ] create a global used liquidiity value (available/utilized liquidity)
- [ ] need to add base mainnet to rpc const list
- [ ] test escrow withdraw function
- [ ] test escrow transfer function (used by withdraw?)
- [ ] bot to convert funds into stable reserves, need to log this (global state variable)


### Edge cases TODO
- [ ] the use connects a wallet but needs to associate it with a different account
- [ ] pen test
- [ ] stress test
- [ ] cost analysis report cost<>service provider

### Later TODO
- [ ] eddsa key validation
- [ ] solana escrow listener
  - [ ] listen to solana blocks
  - [ ] parse solana block events
  - [ ] listen to eclipse blocks
  - [ ] parse solana block events
- [ ] solana escrow contract
  - [ ] create contract
  - [ ] anchor test
  - [ ] live test
- [ ] create library for generalized utils resources
- [ ] utilize grpc
- [ ] rename get_orders_by_address to get_orders_by_user_address
- [ ] tracking protocol
  - [ ] protocol stats
  - [ ] statis gui
  - [ ] downtime/maintenence feed/history
- [ ] backup
  - [ ] tables snapshot
  - [ ] tables rollback
- [ ] enable protection (DiD)
  - [ ] kill switches
    - [ ] db
    - [ ] smart contracts
    - [ ] backend
    - [ ] bots
  - [ ] failsafe quarantine 
  - [ ]
- [ ] reduce function/table access privalages
- [ ] docker image to run tests
- [ ] metamask snap for liquidity and blu staking