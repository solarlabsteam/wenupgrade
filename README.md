# wenupgrade

![Latest release](https://img.shields.io/github/v/release/solarlabsteam/wenupgrade)
[![Actions Status](https://github.com/solarlabsteam/wenupgrade/workflows/test/badge.svg)](https://github.com/solarlabsteam/wenupgrade/actions)

wenupgrade is a tool that displays an approximate date on when some block would be generated. Can be useful if you're running a fullnode and want to check how much time left till the specific block (for example, for calculating time till the chain halts to upgrade).

## How can I set it up?

Download the latest release from [the releases page](https://github.com/solarlabsteam/missed-blocks-checker/releases/). After that, you should unzip it and you are ready to go:

```sh
./wenupgrade 1337
```

## How does it work?

It queries the latest block and the block generated 100 blocks earlier, calculates the average block time, then displays the approximate date the specified block would be generated. 

## How can I configure it?

You can pass the artuments to the executable file to configure it. Here is the parameters list:

- `--tendermint-rpc` - the tendermint node URL. Defaults to `localhost:26657`
- `--log-devel` - logger level. Defaults to `info`. You can set it to `debug` to make it more verbose.
