package types

const (
	ERC20EventAbi   = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"tokenOwner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokens","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"}]`
	ERC721EventAbi  = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_owner","type":"address"},{"indexed":true,"internalType":"address","name":"_approved","type":"address"},{"indexed":true,"internalType":"uint256","name":"_tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_owner","type":"address"},{"indexed":true,"internalType":"address","name":"_operator","type":"address"},{"indexed":false,"internalType":"bool","name":"_approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":true,"internalType":"address","name":"_to","type":"address"},{"indexed":true,"internalType":"uint256","name":"_tokenId","type":"uint256"}],"name":"Transfer","type":"event"}]`
	ERC1155EventAbi = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_operator","type":"address"},{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":true,"internalType":"address","name":"_to","type":"address"},{"indexed":false,"internalType":"uint256[]","name":"_ids","type":"uint256[]"},{"indexed":false,"internalType":"uint256[]","name":"_values","type":"uint256[]"}],"name":"TransferBatch","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_operator","type":"address"},{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":true,"internalType":"address","name":"_to","type":"address"},{"indexed":false,"internalType":"uint256","name":"_id","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_value","type":"uint256"}],"name":"TransferSingle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_account","type":"address"},{"indexed":true,"internalType":"address","name":"_operator","type":"address"},{"indexed":false,"internalType":"bool","name":"_approved","type":"bool"}],"name":"ApprovalForAll","type":"event"}]`
)

const (
	ERC20   = "ERC20"
	ERC721  = "ERC721"
	ERC1155 = "ERC1155"
	Unknown = "Unknown"
)

const (
	// ERC20 | ERC721 Events
	EventTransferSignature = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	EventApprovalSignature = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
	// ERC721 | ERC1155 Events
	EventApprovalForAllSignature = "0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31"
	// ERC1155 Events
	EventTransferSingleSignature = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"
	EventTransferBatchSignature  = "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"
)

func GetERC20Sigs() map[string]struct{} {
	return map[string]struct{}{
		"18160ddd": {}, // "totalSupply()",
		"70a08231": {}, // "balanceOf(address)",
		"a9059cbb": {}, // "transfer(address,uint256)",
		"dd62ed3e": {}, // "allowance(address,address)",
		"095ea7b3": {}, // "approve(address,uint256)",
		"23b872dd": {}, // "transferFrom(address,address,uint256)",
	}
}

func GetERC721Sigs() map[string]struct{} {
	return map[string]struct{}{
		"06fdde03": {}, // "name()"
		"95d89b41": {}, // "symbol()"
		"c87b56dd": {}, // "tokenURI(uint256)"
		"70a08231": {}, // "balanceOf(address)",
		"6352211e": {}, // "ownerOf(uint256)",
		"095ea7b3": {}, // "approve(address,uint256)",
		"081812fc": {}, // "getApproved(uint256)",
		"a22cb465": {}, // "setApprovalForAll(address,bool)",
		"e985e9c5": {}, // "isApprovedForAll(address,address)",
		"23b872dd": {}, // "transferFrom(address,address,uint256)",
		"42842e0e": {}, // "safeTransferFrom(address,address,uint256)",
		"b88d4fde": {}, // "safeTransferFrom(address,address,uint256,bytes)",
		"150b7a02": {}, // "onERC721Received(address,address,uint256,bytes)",
	}
}

func GetERC721EnumerableSigs() map[string]struct{} {
	return map[string]struct{}{
		"0x18160ddd": {}, // "totalSupply()"
		"0x2f745c59": {}, // "tokenOfOwnerByIndex(address,uint256)"
		"0x4f6ccce7": {}, // "tokenByIndex(uint256)"
	}
}

func GetERC1155Sigs() map[string]struct{} {
	return map[string]struct{}{
		"0e89341c": {}, // "uri(uint256)"
		// "00fdd58e": {}, // "balanceOf(address,uint256)"
		"4e1273f4": {}, // "balanceOfBatch(address[],uint256[])"
		"a22cb465": {}, // "setApprovalForAll(address,bool)"
		"e985e9c5": {}, // "isApprovedForAll(address,address)"
		"f242432a": {}, // "safeTransferFrom(address,address,uint256,uint256,bytes)"
		"2eb2c2d6": {}, // "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)"
		"f23a6e61": {}, // "onERC1155Received(address,address,uint256,uint256,bytes)"
		"bc197c81": {}, // "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)"
	}
}
