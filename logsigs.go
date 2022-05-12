package main

// -----------------------------------------------
// ----------------- METHOD SIGS -----------------
// -----------------------------------------------

const (

	// IERC20
	NameMethod         = `name()`
	SymbolMethod       = `symbol()`
	DecimalsMethod     = `decimals()`
	TotalSupplyMethod  = `totalSupply()`
	BalanceOfMethod    = `balance(address)`
	TransferMethod     = `transfer(address,uint256)`
	AllowanceMethod    = `allowance(address,address)`
	ApproveMethod      = `approve(address,uint256)`
	TransferFromMethod = `transferFrom(address,address,uint256)`

	// Pair
	Token0                 = `token0()`
	Token1                 = `token1()`
	GetReserves            = `getReserves()`
	Price0Cum              = `price0CumulativeLast()`
	Price1Cum              = `price1CumulativeLast()`
	KLastMethod            = `kLast()`
	MinimumLiquidityMethod = `MINIMUM_LIQUIDITY()`

	SetPrice = `setPrice(uint256)`

	DepositEvent = `Deposit(address,uint256)`

	// ERC20
	TransferEvent = `Transfer(address,address,uint256)`
	ApprovalEvent = `Approval(address,address,uint256)`

	// Factory
	PairCreatedEvent = `PairCreated(address,address,address,uint256)`

	// Pair
	MintEvent = `Mint(address,uint256,uint256)`
	BurnEvent = `Burn(address,uint256,uint256,address)`
	SwapEvent = `Swap(address,uint256,uint256,uint256,uint256,address)`
	SyncEvent = `Sync(uint112,uint112)`

	// logs.data = pair & btime. topics = token addresses. testing hashes
	sigPairCreated = `0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9`
	sigTransfer    = `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`
	paircreations  = `0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9`
)

var MethodsPostHash = map[string]string{
	keccakfull(TransferEvent):    TransferEvent,
	keccakfull(ApprovalEvent):    ApprovalEvent,
	keccakfull(MintEvent):        MintEvent,
	keccakfull(BurnEvent):        BurnEvent,
	keccakfull(SwapEvent):        SwapEvent,
	keccakfull(SyncEvent):        SyncEvent,
	keccakfull(PairCreatedEvent): PairCreatedEvent,
	keccakfull(DepositEvent):     DepositEvent,
	keccakfull(Token0):           Token0,
}

/*
Contract Types:
	• Tokens
	• Pairs (LP Tokens)
	• Factory (creator of LP tokens) (MochiFactory)

	Contract analysis is done via method calls on the contract & reading logs.

	Contracts calls looks like:
	(func keccak sig)[:10]  -   0x123232300000...
	(args if any) always pad to 32 bytes.

	Compile lists of relevant events.
*/
