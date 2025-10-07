//

package main

import (
    "math/big"
    "testing"

    "github.com/brevis-network/brevis-sdk/sdk"
    "github.com/brevis-network/brevis-sdk/test"
    "github.com/ethereum/go-ethereum/common"
)

func TestAppCircuit(t *testing.T) {
    chainId := uint64(1)
    rpcUrl := "https://mainnet.infura.io/v3/392b6fec32744b34a4850eb2ce3cea2c"
    outDir := "/tmp/brevisapp_out"
    app, err := sdk.NewBrevisApp(chainId, rpcUrl, outDir)
    if err != nil {
        t.Fatalf("Failed to create BrevisApp: %v", err)
    }

    txHash := common.HexToHash("0xcd88108ce4961294b4544342946f4f9696fd734fbc1e1452834c4923423e514a")
    wethAddress := common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
    deaiAddress := common.HexToAddress("0x1495bc9e44af1f8bcb62278d2bec4540cf0c05ea")
    uniswapPoolAddress := common.HexToAddress("0x1385fc1fe0418ea0b4fcf7adc61fc7535ab7f80d")
    transferEventID := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
    swapEventID := common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67")

    app.AddReceipt(sdk.ReceiptData{
        TxHash:   txHash,
        BlockNum: big.NewInt(18446788),
        Fields: []sdk.LogFieldData{
            {
                Contract:   wethAddress,
                EventID:    transferEventID,
                LogPos:     0,
                IsTopic:    true,
                FieldIndex: 1,
                Value:      common.HexToHash("0x0000000000000000000000001385Fc1Fe0418ea0B4Fcf7Adc61FC7535AB7F80d"),
            },
            {
                Contract:   wethAddress,
                EventID:    transferEventID,
                LogPos:     0,
                IsTopic:    true,
                FieldIndex: 2,
                Value:      common.HexToHash("0x000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564"),
            },
        },
    })

    app.AddReceipt(sdk.ReceiptData{
        TxHash:   txHash,
        BlockNum: big.NewInt(18446788),
        Fields: []sdk.LogFieldData{
            {
                Contract:   deaiAddress,
                EventID:    transferEventID,
                LogPos:     1,
                IsTopic:    true,
                FieldIndex: 1,
                Value:      common.HexToHash("0x00000000000000000000000074BED4Ce9D183F69dBb51a394FfF48ED861523E1"),
            },
            {
                Contract:   deaiAddress,
                EventID:    transferEventID,
                LogPos:     1,
                IsTopic:    true,
                FieldIndex: 2,
                Value:      common.HexToHash("0x0000000000000000000000001385Fc1Fe0418ea0B4Fcf7Adc61FC7535AB7F80d"),
            },
        },
    })

    app.AddReceipt(sdk.ReceiptData{
        TxHash:   txHash,
        BlockNum: big.NewInt(18446788),
        Fields: []sdk.LogFieldData{
            {
                Contract:   uniswapPoolAddress,
                EventID:    swapEventID,
                LogPos:     2,
                IsTopic:    true,
                FieldIndex: 2,
                Value:      common.HexToHash("0x000000000000000000000000E592427A0AEce92De3Edee1F18E0157C05861564"),
            },
        },
    })

    appCircuit := &AppCircuit{}
    appCircuitAssignment := &AppCircuit{
        WETHAddress:      sdk.ConstFromBigEndianBytes(wethAddress.Bytes()),
        DEAIAddress:      sdk.ConstFromBigEndianBytes(deaiAddress.Bytes()),
        UniswapPool:      sdk.ConstFromBigEndianBytes(uniswapPoolAddress.Bytes()),
        TransferEventID:  sdk.ConstFromBigEndianBytes(transferEventID.Bytes()),
        SwapEventID:      sdk.ConstFromBigEndianBytes(swapEventID.Bytes()),
    }

    circuitInput, err := app.BuildCircuitInput(appCircuitAssignment)
    if err != nil {
        t.Fatalf("Failed to build circuit input: %v", err)
    }

    test.ProverSucceeded(t, appCircuit, appCircuitAssignment, circuitInput)
}
