package packet

type Port struct {
	start uint16
	end   uint16
}

func Range(start, end uint16) Port {
	return Port{start, end}
}

func Point(port uint16) Port {
	return Port{port, port}
}

func (p Port) Contains(port uint16) bool {
	return port >= p.start && port <= p.end
}
