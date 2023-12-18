module chwojkofrank.com/17

require (
	chwojkofrank.com/astar v0.0.0-20190116155806-2b7a7f8c9d4a
	chwojkofrank.com/dijkstra v0.0.0-00010101000000-000000000000
)

require chwojkofrank.com/cursor v0.0.0-20190116155806-2b7a7f8c9d4a // indirect

replace chwojkofrank.com/astar => ../astar

replace chwojkofrank.com/dijkstra => ../dijkstra

replace chwojkofrank.com/cursor => ../cursor

go 1.21.5
