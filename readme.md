# Harmony One Chain Explorer.
Indexes Account/Contract/Token info from Genesis Block into SQL tables.\
Serves explorer to :1234 via Websockets, as it fills.\
I mainly use certain functions on their own to compile Token Contracts info and their transactions, harmony node api allows filtering txs based on contract address, otherwise sync time increases with height, because tx density increases with height, albeit plateauing with Harmony internal block/max tx constraints.\
When running later in the chain, block processing is essentially rate limiter.\
Run over section you want to inspect, go from recent to monitor current chain activity.\
The table sizes become quite large and may need to be split depending on your system.\
Useful for compiling token holder amounts, reserves ratios, and other details on liquidity pools.\
Method Sigs are common to all ERC20 chains.\
Alter solidity contract to deploy token. According to IERC20 interface standards i.e. Mintable, Burnable, Ownable etc... Ref: https://github.com/OpenZeppelin/openzeppelin-contracts \
Tokens are created by compiling .sol file, marshalling the bytes into tx data, and sending tx.\
Also deploy here remix.ethereum.org to Ethereum, Binance or Harmony Chains.

This does not spin up a Harmony Node, rather it uses https://api.s0.t.hmny.io. There is also wss://ws.s0.t.hmny.io.
Official Explorer: https://explorer.harmony.one/

Notes/
* First Pair a48797a204b5e02b4a5daca80cbc4aa09212225f (Crowdsale: beneficiary is the zero address). Block 38e8ea22a0773dde754b323ef38d99bd4361abd229b67dd9c9c52ef38dc4e443
* First Block dfeff1fba1aeed89fb75ef4ee9bf9e0fca1ff9b26d78d471565bf151f965274b
* Second Block 61ce03ef5efa374b0d0d527ea38c3d13cb05cf765a4f898e91a5de1f6b224cdd
* First Token 95eb8075f7f7afb37f5f4cc90663f709040088e1 (OneBTC, 1BTC)

Token begin from 5481185