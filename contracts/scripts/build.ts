import hre from "hardhat"
import fs from "fs"

async function main() {
    const paths = await hre.artifacts.getArtifactPaths()

    paths.forEach((path) => {
        const artifactName = path.split("/").pop()?.split(".")[0]
        if (artifactName === undefined) {
            throw new Error("could not parse artifact name")
        }

        if (!['ERC721', 'IIERC721', 'ERC20MinterPauserBurnerDecimals', 'ERC20NonTransferable'].includes(artifactName)) return

        const fullyQualifiedName = `contracts/${artifactName}.sol:${artifactName}`
        const artifact = hre.artifacts.readArtifactSync(fullyQualifiedName)
        const bin = artifact.bytecode.split('0x')[1]
        if (bin === undefined) {
            throw new Error("could not parse bytecode")
        }

        const compiledContract = {
            abi: JSON.stringify(artifact.abi),
            bin,
            contractName: artifactName,
        }

        fs.writeFileSync(`compiled_contracts/${artifactName}.json`, JSON.stringify(compiledContract, null, 2))
    })
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
