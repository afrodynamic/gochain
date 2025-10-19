# Gochain <!-- omit in toc -->

This project is an application showcasing a bridge adapter for interfacing with multiple blockchains.

## Table of Contents<!-- omit in toc -->

- [Overview](#overview)
- [Usage](#usage)
  - [Requirements](#requirements)
  - [Setup and Running](#setup-and-running)

## Overview

Gochain features a modular architecture that allows for easy integration with other blockchain networks by providing a bridge adapter that facilitates communication between different blockchains. The abstracted design streamlines the process of adding support for new networks, making it easier for developers to expand the application's capabilities.

The project integrates and leverages Metamask for wallet management and supports multiple testnet networks.

## Usage

To use this project please follow the instructions provided below.

### Requirements

You will need the following dependencies to run this project:

- `git`, for cloning the project ([download](https://git-scm.com/downloads))
- `pnpm`, for managing frontend application packages ([installation guide](https://pnpm.io/installation))
- An Infura API key, for utilizing the bridge adapter functionality ([sign up](https://infura.io/))
- A web3-compatible wallet using Metamask ([download](https://metamask.io/download/))

### Setup and Running

With the requirements installed, follow the steps below to run the project:

1. Clone this repository to your local machine

   ```shell
   git clone https://github.com/afrodynamic/gochain.git
   ```

2. Navigate to the project directory

   ```shell
   cd gochain/web
   ```

3. Obtain an Infura API key by signing up at [Infura](https://infura.io/). You will need the key to utilize the bridge adapter functionality.
4. Set the `NEXT_PUBLIC_INFURA_API_KEY` environment variable with your Infura API key in the project's `.env` file.

   ```shell
   NEXT_PUBLIC_INFURA_API_KEY=your_infura_api_key_here
   ```

5. Install the required frontend packages

   ```shell
   pnpm install
   ```

6. Run the frontend

   ```shell
   pnpm run dev
   ```
