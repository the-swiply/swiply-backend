package edmonds

type Graph[T any] struct {
	g      map[string][]string
	values map[string]T
	key    func(T) string
}

type Pair[T any] struct {
	First, Second T
}

func NewGraph[T any](key func(T) string) *Graph[T] {
	return &Graph[T]{
		g:      make(map[string][]string),
		values: make(map[string]T),
		key:    key,
	}
}

func (g *Graph[T]) AddEdge(from, to T) {
	fromKey := g.key(from)
	toKey := g.key(to)
	g.g[fromKey] = append(g.g[fromKey], toKey)

	if _, ok := g.values[fromKey]; !ok {
		g.values[fromKey] = from
	}

	if _, ok := g.values[toKey]; !ok {
		g.values[toKey] = to
	}
}

func (g *Graph[T]) lca(a, b string, base, matched, parent map[string]string) string {
	used := make(map[string]bool)
	for {
		a = base[a]
		used[a] = true
		if matched[a] == "" {
			break
		}
		a = parent[matched[a]]
	}

	for {
		b = base[b]
		if used[b] {
			return b
		}
		b = parent[matched[b]]
	}
}

func (g *Graph[T]) markPath(v, b, children string, blossom map[string]bool, base, matched, parent map[string]string) {
	for base[v] != b {
		blossom[base[v]] = true
		blossom[base[matched[v]]] = true
		parent[v] = children
		children = matched[v]
		v = parent[matched[v]]
	}
}

func (g *Graph[T]) findPath(root string, matched, parent map[string]string) string {
	for k := range parent {
		delete(parent, k)
	}

	used := make(map[string]bool)
	used[root] = true

	base := make(map[string]string)
	for key := range g.g {
		base[key] = key
	}

	q := make([]string, 0, len(base))
	qh, qt := 0, 0
	q[qt] = root
	qt++

	for qh < qt {
		v := q[qh]
		qh++

		for _, to := range g.g[v] {
			if base[v] == base[to] || matched[v] == to {
				continue
			}

			if to == root || matched[to] != "" && parent[matched[to]] != "" {
				curbase := g.lca(v, to, base, matched, parent)
				blossom := make(map[string]bool)

				g.markPath(v, curbase, to, blossom, base, matched, parent)
				g.markPath(to, curbase, v, blossom, base, matched, parent)

				for key := range g.g {
					if blossom[base[key]] {
						base[key] = curbase
						if !used[key] {
							used[key] = true
							q[qt] = key
							qt++
						}
					}
				}
			} else if parent[to] == "" {
				parent[to] = v
				if matched[to] == "" {
					return to
				}
				to = matched[to]
				used[to] = true
				q[qt] = to
				qt++
			}
		}
	}

	return ""
}

func (g *Graph[T]) MatchPairs() []Pair[T] {
	matched := make(map[string]string)
	parent := make(map[string]string)
	for key := range g.g {
		if val := matched[key]; val == "" {
			for lastVertex := g.findPath(key, matched, parent); lastVertex != ""; {
				previousVertex := parent[lastVertex]
				matchedPreviousVertex := matched[previousVertex]
				matched[lastVertex] = previousVertex
				matched[previousVertex] = lastVertex
				lastVertex = matchedPreviousVertex
			}
		}
	}

	var ans []Pair[T]
	for key, value := range matched {
		if key < value && key != "" && value != "" {
			ans = append(ans, Pair[T]{First: g.values[key], Second: g.values[value]})
		}
	}

	return ans
}
