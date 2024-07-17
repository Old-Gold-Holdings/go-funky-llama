package plugins

type Git struct{}

func (g *Git) New(upstream string) {}

func (g *Git) Add(file string) {}

func (g *Git) AddAll() {}

func (g *Git) Commit(message string) {}

func (g *Git) Push() {}

func (g *Git) Pull() {}

func (g *Git) Checkout(branch string) {}
