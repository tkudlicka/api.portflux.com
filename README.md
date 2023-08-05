# api.portflux.com
api.portflux.com: A Go service for stock tracking using PostgreSQL and Hexagonal Architecture.

## Table of Contents

<!-- TOC -->

- [api.portflux.com](#api.portflux.com)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Usage](#usage)
  - [Requirements](#requirements)
  - [Installation](#installation)
  - [Author Information](#author-information)
  <!-- TOC -->

## Overview
The `api.portflux.com` project is a Go-based service designed to track the gains of a stock portfolio. Utilizing the PostgreSQL database and adhering to the principles of Hexagonal Architecture, it ensures the application is maintainable, scalable, and loosely coupled.


## Usage
The `api.portflux.com` service operates primarily behind the scenes, continuously monitoring various stock sources and synchronizing data with the portflux.com API's PostgreSQL database. While it doesn't expose an API for direct interaction, the service ensures that the stock portfolio data remains accurate and up-to-date.


## Requirements
To run the api.portflux.com backend, you will need the following:

- Go programming language (version 1.20.5)
- PostgreSQL database
- Dependencies specified in the project's go.mod file

## Installation
Use docker-compose to start requirements resources

```bash
docker-compose up -d
```

Create a .env file with this default envs in env.example

```bash
make up
```

## Author Information

This module is maintained by the contributors listed on [GitHub](https://github.com/tkudlicka/api.portflux.com/graphs/contributors).

