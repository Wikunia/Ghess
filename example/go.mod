module github.com/Wikunia/Ghess/example

go 1.16

replace github.com/Wikunia/Ghess/ghess => ../ghess

require (
	github.com/Wikunia/Ghess/ghess v0.0.0-00010101000000-000000000000
	github.com/gofiber/websocket/v2 v2.0.3 // indirect
)
