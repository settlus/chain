import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  defaultNetwork: "settlus",
  solidity: "0.8.8",
  paths: {
    artifacts: "./artifacts",
    cache: "./cache",
    sources: "./contracts",
    tests: "./test",
  },
  networks: {
    settlus: {
      url: "http://localhost:8545",
      accounts: {
        // this is the sample mnemonic from ../config.yml
        mnemonic: "equal broken goose strong twenty upgrade cool pen run opinion gain brick husband repeat magnet foam creek purse alcohol this margin lunch hip birth",
      }
    }
  },
  gasReporter: {
    currency: "SETL",
    enabled: true,
  }
};

export default config;
