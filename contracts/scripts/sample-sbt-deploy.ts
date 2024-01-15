import hre from 'hardhat'

const sampleTreasuryAddress = '0x0e74F59A0DA956309A2Fb17E87067ACc8b597117'

const deploySBT = async (contractName: string) => {
    console.log(`Deploying ${contractName}...`)
    const factory = await hre.ethers.getContractFactory(contractName)
    const contract = await factory.deploy("earned BLUC", "eBLUC", sampleTreasuryAddress)

    const receipt = await contract.deploymentTransaction()?.wait(2)
    return receipt?.contractAddress
}

async function main() {
    await hre.run('compile')
    console.log(`Compiling...`)
    const contractName = 'ERC20NonTransferable'
    const contractAddress = await deploySBT(contractName)
    if (!contractAddress) {
        throw Error('Contract address is not set')
    }
    console.log('Soul Bound Token deployed at: ' + contractAddress)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })
