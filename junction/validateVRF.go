package junction

import (
	"context"
	"github.com/airchains-network/decentralized-sequencer/node/shared"

	"fmt"
	"github.com/airchains-network/decentralized-sequencer/junction/types"
	logs "github.com/airchains-network/decentralized-sequencer/log"
	"github.com/airchains-network/decentralized-sequencer/utilis"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func ValidateVRF(serializedRc []byte) bool {
	jsonRpc, stationId, accountPath, accountName, addressPrefix, err := utilis.GetJunctionDetails()
	if err != nil {
		logs.Log.Error("can not get junctionDetails.json data: " + err.Error())
		return false
	}

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error creating account registry: %v", err))
		return false
	}

	newTempAccount, err := registry.GetByName(accountName)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error getting account: %v", err))
		return false
	}

	newTempAddr, err := newTempAccount.Address(addressPrefix)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error getting address: %v", err))
		return false
	}

	ctx := context.Background()
	gas := utilis.GenerateRandomWithFavour(510, 1000, [2]int{520, 700}, 0.7)
	gasFees := fmt.Sprintf("%damf", gas)
	logs.Log.Warn(fmt.Sprintf("Gas Fees Used for validate VRF transaction is: %s\n", gasFees))
	accountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(jsonRpc), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees(gasFees))
	if err != nil {
		logs.Log.Error("Error creating account client")
		return false
	}

	podNumber := shared.GetPodState().LatestPodHeight
	msg := types.MsgValidateVrf{
		Creator:      newTempAddr,
		StationId:    stationId,
		PodNumber:    podNumber,
		SerializedRc: serializedRc,
	}

	txRes, errTxRes := accountClient.BroadcastTx(ctx, newTempAccount, &msg)
	if errTxRes != nil {
		logs.Log.Error("error in transaction" + errTxRes.Error())
		return false
	}

	logs.Log.Info("Transaction Hash For VRF Validation: " + txRes.TxHash)

	return true
}