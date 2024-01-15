import hre from 'hardhat'
import { Contract, ContractTransaction } from 'ethers'

type MintContract = Contract & {
    safeMint: (to: string, uri: string) => Promise<ContractTransaction>
}

const deployNFT = async (contractName: string) => {
    console.log(`Deploying ${contractName}...`)
    const factory = await hre.ethers.getContractFactory(contractName)
    const contract = await factory.deploy()

    const receipt = await contract.deploymentTransaction()?.wait(2)
    console.log(`${contractName} deployed to ${receipt?.contractAddress}`)

    return receipt?.contractAddress
}

const mintNFT = async (contractName: string, contractAddress: string, receiverAddress: string) => {
    for (let i = 0; i < 5; i++) {
        try {
            const sampleContract = (await hre.ethers.getContractAt(
                contractName,
                contractAddress,
            )) as MintContract
            await sampleContract.safeMint(
              receiverAddress,
              `https://example.com/${i}.json`,
            )
            console.log(`Mint ${contractName} ${i} finished`)
        } catch (error: unknown) {
            console.log(`Failed: ${contractName} ${i}: ${error}`)
        }
    }
    console.log(`\ncheck nft ownership with \n$ chaind q nftownership get-nft-owner settlus_5371-1 ${contractAddress} 0x0`)
}

async function main() {
    await hre.run('compile')
    console.log(`Compiling...`)
    const contractName = 'SampleNFT'
    const receiverAddress = '0x626ff58bce1d71d3e9466bfa2973a1aaf8172082'
    const contractAddress = await deployNFT(contractName)
    if (!contractAddress) {
        throw Error('Contract address is not set')
    }
    await mintNFT(contractName, contractAddress, receiverAddress)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })
