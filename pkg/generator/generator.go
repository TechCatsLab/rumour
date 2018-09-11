/*
 * Revision History:
 *     Initial: 2018/07/27        Tong Yuehong
 */

package generator

type Generator struct {
	generator chan uint64
	start uint64
}

func New(s uint64, shutdown chan struct{}) *Generator {
	g := &Generator{
		generator: make(chan uint64, 4096),
		start: s,
	}

	go g.generate(shutdown)

	return g
}

func (g *Generator) generate(shutdown chan struct{}) {
	for {
		g.start++

		select {
		case g.generator <- g.start:
		case <-shutdown:
			return
		}
	}
}

// Get a id from generator.
func (g *Generator) Get() uint64 {
	return <-g.generator
}
