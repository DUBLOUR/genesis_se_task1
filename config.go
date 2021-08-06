package main

const LogFile string = "/var/log/btc-requester/responses.log"
const UsersFile string = "/etc/btc-requester/users.csv"
const ServerPort string = ":9990"

const MarketEndpoint string = "https://api3.binance.com/api/v3/ticker/price?"
const PasswordSalt string = "Yeeh_zMVk3"
const TokenLength int = 12
