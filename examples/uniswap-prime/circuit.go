// circuit.go

package main

import (
    "github.com/brevis-network/brevis-sdk/sdk"
)

type AppCircuit struct {
    WETHAddress      sdk.Bytes32
    DEAIAddress      sdk.Bytes32
    UniswapPool      sdk.Bytes32
    TransferEventID  sdk.Bytes32
    SwapEventID      sdk.Bytes32
}

var _ sdk.AppCircuit = &AppCircuit{}

func (c *AppCircuit) Allocate() (maxReceipts, maxStorage, maxTxs int) {
    return 32, 0, 0
}

func (c *AppCircuit) Define(api *sdk.CircuitAPI, in sdk.DataInput) error {
    receipts := sdk.NewDataStream(api, in.Receipts)

    transferLogs := sdk.Filter(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
        contractB32 := api.ToBytes32(receipt.Fields[0].Contract)
        eventIDB32 := api.ToBytes32(receipt.Fields[0].EventID)
        isTransfer := api.Uint248.And(
            api.Uint248.Or(
                api.Bytes32.IsEqual(contractB32, c.WETHAddress),
                api.Bytes32.IsEqual(contractB32, c.DEAIAddress),
            ),
            api.Bytes32.IsEqual(eventIDB32, c.TransferEventID),
        )
        return isTransfer
    })

    sdk.AssertEach(transferLogs, func(receipt sdk.Receipt) sdk.Uint248 {
        return sdk.ConstUint248(1)
    })

    swapLogs := sdk.Filter(receipts, func(receipt sdk.Receipt) sdk.Uint248 {
        contractB32 := api.ToBytes32(receipt.Fields[0].Contract)
        eventIDB32 := api.ToBytes32(receipt.Fields[0].EventID)
        isSwap := api.Uint248.And(
            api.Bytes32.IsEqual(contractB32, c.UniswapPool),
            api.Bytes32.IsEqual(eventIDB32, c.SwapEventID),
        )
        return isSwap
    })

    sdk.AssertEach(swapLogs, func(receipt sdk.Receipt) sdk.Uint248 {
        return sdk.ConstUint248(1)
    })

    return nil
}
