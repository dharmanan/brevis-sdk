package tradingvolume

import (
	"os"
	"fmt"
	"testing"

	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/brevis-network/brevis-sdk/test"
	"github.com/ethereum/go-ethereum/common"
)

// In this example, we want to analyze the `Swap` events emitted by Uniswap's
// UniversalRouter contract. Let's declare the fields we want to use:

func TestCircuit(t *testing.T) {
	rpc := os.Getenv("INFURA_RPC_URL")
	outDir := "$HOME/circuitOut/myBrevisApp"
	app, err := sdk.NewBrevisApp(1, rpc, outDir)
	check(err)

	// Adding a receipt query into the querier
	// In this tx, the user sold USDC and took native ETH out
	app.AddReceipt(sdk.ReceiptData{
		TxHash: common.HexToHash("0x53b37ec7975d217295f4bdadf8043b261fc49dccc16da9b9fc8b9530845a5794"),
		Fields: []sdk.LogFieldData{
			// USDC.Transfer.from (LogPos: 2, topic, index 1, contract: USDC, event: Transfer, value: userAddr)
			{LogPos: 2, IsTopic: true, FieldIndex: 1, Contract: common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), EventID: common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), Value: common.HexToHash("0x000000000000000000000000aefb31e9eeee2822f4c1cbc13b70948b0b5c0b3c")},
			// USDCPool.Swap.amount0 (LogPos: 3, not topic, index 0, contract: USDCPool, event: Swap, value: 123456)
			{LogPos: 3, IsTopic: false, FieldIndex: 0, Contract: common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"), EventID: common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"), Value: common.HexToHash("0x000000000000000000000000000000000000000000000000000000000001e240")},
			// USDCPool.Swap.recipient (LogPos: 3, topic, index 2, contract: USDCPool, event: Swap, value: userAddr)
			{LogPos: 3, IsTopic: true, FieldIndex: 2, Contract: common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"), EventID: common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"), Value: common.HexToHash("0x000000000000000000000000aefb31e9eeee2822f4c1cbc13b70948b0b5c0b3c")},
		},
	})
	// More receipts can be added, but in this example we only add one to keep it simple
	// app.AddReceipt(...)
	// app.AddReceipt(...)

	// Initialize our AppCircuit and prepare the circuit assignment
	userAddr := common.HexToAddress("0xaefB31e9EEee2822f4C1cBC13B70948b0B5C0b3c")
	appCircuit := &AppCircuit{
		UserAddr: sdk.ConstUint248(userAddr),
	}
	appCircuitAssignment := &AppCircuit{
		UserAddr: sdk.ConstUint248(userAddr),
	}

	// Execute the added queries and package the query results into circuit inputs
	in, err := app.BuildCircuitInput(appCircuit)
	check(err)

	///////////////////////////////////////////////////////////////////////////////
	// Testing
	///////////////////////////////////////////////////////////////////////////////

	// Use the test package to check if the circuit can be solved using the given
	// assignment
	test.ProverSucceeded(t, appCircuit, appCircuitAssignment, in)
}

func TestE2E(t *testing.T) {
	// The compiled circuit, proving key, and verifying key are saved to outDir,,
	// query data will be stored under outDir/input and
	// the downloaded SRS in the process is saved to srsDir
	outDir := "$HOME/circuitOut/tradingvolume"
	srsDir := "$HOME/kzgsrs"
	rpc := "https://mainnet.infura.io/v3/392b6fec32744b34a4850eb2ce3cea2c"
	app, err := sdk.NewBrevisApp(1, rpc, outDir)
	check(err)

	// Adding a receipt query into the querier
	// In this tx, the user sold USDC and took native ETH out
	app.AddReceipt(sdk.ReceiptData{
		TxHash: common.HexToHash("53b37ec7975d217295f4bdadf8043b261fc49dccc16da9b9fc8b9530845a5794"),
		Fields: []sdk.LogFieldData{
			// LogPos: 2 önce, sonra LogPos: 3'ler
			{LogPos: 2, IsTopic: true, FieldIndex: 1},  // field: USDC.Transfer.from
			{LogPos: 3, IsTopic: false, FieldIndex: 0}, // field: USDCPool.Swap.amount0
			{LogPos: 3, IsTopic: true, FieldIndex: 2},  // field: USDCPool.Swap.recipient (topic field)
		},
	})
	// More receipts can be added, but in this example we only add one to keep it simple
	// app.AddReceipt(...)
	// app.AddReceipt(...)

	// Initialize our AppCircuit and prepare the circuit assignment
	appCircuit := &AppCircuit{
		// you need to give every custom input an assignment or otherwise the circuit won't compile
		UserAddr: sdk.ConstUint248(0),
	}
	appCircuitAssignment := &AppCircuit{
		UserAddr: sdk.ConstUint248(common.HexToAddress("0xaefB31e9EEee2822f4C1cBC13B70948b0B5C0b3c")),
	}

	// Execute the added queries and package the query results into circuit inputs
	in, err := app.BuildCircuitInput(appCircuitAssignment)
	check(err)

	///////////////////////////////////////////////////////////////////////////////
	// Testing
	///////////////////////////////////////////////////////////////////////////////

	// Use the test package to check if the input can be proved with the given
	// circuit
	test.ProverSucceeded(t, appCircuit, appCircuitAssignment, in)

	///////////////////////////////////////////////////////////////////////////////
	// Compiling and Setup
	///////////////////////////////////////////////////////////////////////////////

	compiledCircuit, pk, vk, _, err := sdk.Compile(appCircuit, outDir, srsDir, app)
	check(err)

	// Once you saved your ccs, pk, and vk files, you can read them back into memory
	// for use with the provided utils
	compiledCircuit, pk, vk, _, err = sdk.ReadSetupFrom(appCircuit, outDir, app)
	check(err)

	///////////////////////////////////////////////////////////////////////////////
	// Proving
	///////////////////////////////////////////////////////////////////////////////

	fmt.Println(">> prove")
	witness, publicWitness, err := sdk.NewFullWitness(appCircuitAssignment, in)
	check(err)
	proof, err := sdk.Prove(compiledCircuit, pk, witness)
	check(err)

	///////////////////////////////////////////////////////////////////////////////
	// Verifying
	///////////////////////////////////////////////////////////////////////////////

	// The verification of the proof generated by you is done on Brevis' side. But
	// you can also verify your own proof to make sure everything works fine and
	// pk/vk are serialized/deserialized properly
	err = sdk.Verify(vk, publicWitness, proof)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
