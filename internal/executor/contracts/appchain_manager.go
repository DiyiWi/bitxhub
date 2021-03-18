package contracts

import (
	"fmt"
	"strconv"

	appchainMgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/boltvm"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
)

type AppchainManager struct {
	boltvm.Stub
	appchainMgr.AppchainManager
}

// Apply: appchain managers apply appchain method and set appchainAdminDID for this new method
// @appchainAdminDID its address should be the same with am.Caller()
// @appchainMethod should be in format like did:bitxhub:appchain1:.
// @sig is the signature of appchain admin for appchainMethod
func (am *AppchainManager) Apply(appchainAdminDID, appchainMethod string, sig []byte) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	res := am.CrossInvoke(constant.MethodRegistryContractAddr.String(), "Apply",
		pb.String(appchainAdminDID), pb.String(appchainMethod), pb.Bytes(sig))
	if !res.Ok {
		return res
	}
	return boltvm.Success(nil)
}

// Register appchain managers registers appchain info caller is the appchain
// manager address return appchain id and error
func (am *AppchainManager) Register(appchainAdminDID, appchainMethod string, sig []byte, docAddr, docHash, validators string, consensusType int32, chainType, name, desc, version, pubkey string) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	res := am.CrossInvoke(constant.InterchainContractAddr.String(), "Register", pb.String(appchainMethod))
	if !res.Ok {
		return res
	}
	res = am.CrossInvoke(constant.MethodRegistryContractAddr.String(), "Register",
		pb.String(appchainAdminDID), pb.String(appchainMethod),
		pb.String(docAddr), pb.Bytes([]byte(docHash)), pb.Bytes(sig))
	if !res.Ok {
		return res
	}

	return responseWrapper(am.AppchainManager.Register(appchainMethod, validators, consensusType, chainType, name, desc, version, pubkey))
}

// UpdateAppchain updates approved appchain
func (am *AppchainManager) UpdateAppchain(appchainMethod, validators string, consensusType int32, chainType, name, desc, version, pubkey string) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.UpdateAppchain(appchainMethod, validators, consensusType, chainType, name, desc, version, pubkey))
}

//FetchAuditRecords fetches audit records by appchain id
func (am *AppchainManager) FetchAuditRecords(id string) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.FetchAuditRecords(id))
}

// CountApprovedAppchains counts all approved appchains
func (am *AppchainManager) CountApprovedAppchains() *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.CountApprovedAppchains())
}

// CountAppchains counts all appchains including approved, rejected or registered
func (am *AppchainManager) CountAppchains() *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.CountAppchains())
}

// Appchains returns all appchains
func (am *AppchainManager) Appchains() *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.Appchains())
}

// GetAppchain returns appchain info by appchain id
func (am *AppchainManager) GetAppchain(id string) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.GetAppchain(id))
}

// GetPubKeyByChainID can get aim chain's public key using aim chain ID
func (am *AppchainManager) GetPubKeyByChainID(id string) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	return responseWrapper(am.AppchainManager.GetPubKeyByChainID(id))
}

// AuditApply bitxhub manager audit appchain method apply
func (am *AppchainManager) AuditApply(relayAdminDID, proposerMethod string, isApproved int32, sig []byte) *boltvm.Response {
	if res := am.IsAdmin(); !res.Ok {
		return res
	}
	res := am.CrossInvoke(constant.MethodRegistryContractAddr.String(), "AuditApply",
		pb.String(relayAdminDID), pb.String(proposerMethod), pb.Int32(isApproved), pb.Bytes(sig))
	return responseWrapper(res.Ok, res.Result)
}

// Audit bitxhub manager audit appchain register info
func (am *AppchainManager) Audit(proposerMethod string, isApproved int32, desc string) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	if res := am.IsAdmin(); !res.Ok {
		return res
	}
	return responseWrapper(am.AppchainManager.Audit(proposerMethod, isApproved, desc))
}

func (am *AppchainManager) DeleteAppchain(relayAdminDID, method string, sig []byte) *boltvm.Response {
	am.AppchainManager.Persister = am.Stub
	if res := am.IsAdmin(); !res.Ok {
		return res
	}
	res := am.CrossInvoke(constant.InterchainContractAddr.String(), "DeleteInterchain", pb.String(method))
	if !res.Ok {
		return res
	}
	res = am.CrossInvoke(constant.MethodRegistryContractAddr.String(), "Delete", pb.String(relayAdminDID), pb.String(method), pb.Bytes(sig))
	if !res.Ok {
		return res
	}
	return responseWrapper(am.AppchainManager.DeleteAppchain(method))
}

func (am *AppchainManager) IsAdmin() *boltvm.Response {
	ret := am.CrossInvoke(constant.RoleContractAddr.String(), "IsAdmin", pb.String(am.Caller()))
	is, err := strconv.ParseBool(string(ret.Result))
	if err != nil {
		return boltvm.Error(fmt.Errorf("judge caller type: %w", err).Error())
	}

	if !is {
		return boltvm.Error("caller is not an admin account")
	}
	return boltvm.Success([]byte("1"))
}

func responseWrapper(ok bool, data []byte) *boltvm.Response {
	if ok {
		return boltvm.Success(data)
	}
	return boltvm.Error(string(data))
}
