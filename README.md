<p align="center">
  <a href="https://topaz.io">
    <img alt="Topaz" src="topaz.svg" width="240" />
  </a>
</p>
<h3 align="center">
  💎 ⛓️ 🔒
</h3>
<h3 align="center">
  Placing trust in technology
</h3>
<p align="center">
  Topaz is a free and open source "Layer 2" solution that helps developers
  secure data from any application using Ethereum.
</p>

<h3 align="center">
  <a href="https://topaz.io/docs/">Documentation</a>
  <span> · </span>
  <a href="https://topaz.io/tutorial/">Tutorial</a>
  <span> · </span>
  <a href="https://join.slack.com/t/topaz-developers/shared_invite/zt-7bxno80m-nGrysu2fid_vh0iFFr5hUg">Community Slack</a>
</h3>

Topaz is a self-hosted web service that secures data at scale on Ethereum.

- **Simplify Blockchain Software Development.** Integrate any existing
  application with Ethereum, enabling blockchain use cases without smart
  contract development or prior experience in blockchain development.

- **Bypass Transaction Constraints.** Topaz batches transactions into secure
  data structures, anchoring many transactions on chain at once avoiding the
  scaling constraints of transaction throughput.

- **Outscale Any Blockchain.** Our hybrid development pattern using data from
  centralized applications secured on decentralized networks enables the usage
  of blockchains for data integrity at the massive scale the web demands.

- **Trial with Ease.** Our SDKs enable Topaz integration with just three lines
  of code - now, anyone can experiment with using blockchains to create
  transparent, trustless applications with ease.

- **Go to Mainnet for Pennies.** Topaz enables applications that simply could
  not exist due to scaling constraints brought on by the cost of Ethereum
  transactions.

[**Learn how to use Topaz for your dApp project.**](https://topaz.io/docs)

## What’s In This Document

- [Getting Started](#-getting-started)
- [Learning Topaz](#-learning-topaz)
- [How to Contribute](#-how-to-contribute)
- [License](#-license)
- [Thanks](#-thanks)

## Getting Started

You'll need [`go`](https://golang.org/) installed.

You can get Topaz running on your local dev environment in a few minutes with
these four steps:

1. **Configure your dev environment.**

    This project uses golang modules, so simple building or running the
    project(s) will download and install the correct version of dependencies,
    based on `go.mod` and `go.sum`.

    Configure your environment:

    ```shell
    cp .env.example .env
    ```

2. **Start the Docker containers.**

    Next, spin up your containers:

    ```shell
    docker-compose up -d
    ```

3. **Start the API.**

    Launch the API with `make`.

    ```shell
    cd api
    make run
    ```

4. **Start the Batch process.**

    In a new terminal window,

    ```shell
    cd batch
    make run
    ```

### Use the API

Check out `api.md` to learn how the API functions.

## Where you at?

Made with :heart: in:
* British Airways Flight 1575
* Miami, Florida
* Positano, Italy
* Santa Cruz, California
* Cleveland, Ohio
* Boston, Massachusetts
* Copenhagen, Denmark
* Prague, Czech Republic
* Stockholm, Sweden
* Berlin, Germany
* Amsterdam, Netherlands
