# TaoIns Protocol

The TaoIns protocol is an innovative asset standard protocol on BitTensor. Since BitTensor doesnʼt have smart contracts or other type running environments for asset issuing, we intend to define an asset standard to embed the asset creations and operations and included the content into BitTensor blockchain.

The TaoIns protocol implements a technical design similar to Ordinals on Bitcoin, by engraving operational information into the receiving addresses of BitTensor transactions.

A BitTensor address consists a mutable 256-bits filed, i.e the PublicKey. In TaoIns protocol, this 256-bits space is divided into four segments: 32bits as the TaoIns indicator(0xffffffff), 4 bits for asset type, 4 bits for content type, and the remaining 216 bits designated for arbitrary data values. Essentially, each address represents unique information or operation. Any user wishing to perform a specific operation or convey a message related to a particular address simply needs to initiate a $TAO transaction to that address.

## About this repo
This repo contains the sample code for Bittensor Data Parsing, TaoIns Data Recognization and Data Persistence.


## Setup
* Please note, these steps require the Go language version 1.21 or above.

1、Download the project to your local machine, modify the configuration files to suit your own settings(`./config/resources/config.toml`), and then run `go install` to install dependencies.

2、Run the `go build main.go` to create an executable file.

3、Run the `go run main` to start the service.
