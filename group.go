package multiprogress

type renderGroup struct {
	ArrayRenderer

	maxSize int
}

// Add implements RenderGroup.
func (r *renderGroup) Add(child RenderTree) {
	r.ArrayRenderer = append(r.ArrayRenderer, child)
	if len(r.ArrayRenderer) > r.maxSize {
		r.ArrayRenderer = r.ArrayRenderer[1:]
	}
}

var (
	_ RenderGroup = &renderGroup{}
)

func NewRenderGroup(maxSize int) RenderGroup {
	return &renderGroup{maxSize: maxSize}
}
