package parser

type Any interface{}

type Generator interface {
	Get() Any
	Merge(Generator) Generator
}

type valueGenerator struct {
	value Any
}

func (vg *valueGenerator) Get() Any {
	return vg.value
}

func (vg *valueGenerator) Merge(g Generator) Generator {
	// Values cannot be merged. Return the new one
	return g
}

func mkValueGenerator(v Any) *valueGenerator {
	return &valueGenerator{value: v}
}

type ObjectGenerator struct {
	fields map[string]Generator
}

func mkObjectGenerator() *ObjectGenerator {
	return &ObjectGenerator{
		fields: map[string]Generator{},
	}
}

func (obj *ObjectGenerator) add(field string, value Generator) *ObjectGenerator {
	if gen, ok := obj.fields[field]; ok {
		value = gen.Merge(value)
	}
	obj.fields[field] = value
	return obj
}

func (obj *ObjectGenerator) Get() Any {
	res := map[string]Any{}
	for field, vg := range obj.fields {
		res[field] = vg.Get()
	}
	return res
}

func (obj *ObjectGenerator) Merge(g Generator) Generator {
	switch gt := g.(type) {
	case *ObjectGenerator:
		// Objects can be merged together
		res := mkObjectGenerator()
		for f, v := range obj.fields {
			res.add(f, v)
		}
		for f, v := range gt.fields {
			res.add(f, v)
		}
		return res
	default:
		// other types, less so, return the new one
		return g
	}
}

type arrayGenerator []Generator

func (arr *arrayGenerator) Merge(g Generator) Generator {
	// arrays can' t be merged with other generators
	return g
}

func (arr *arrayGenerator) Get() Any {
	res := make([]Any, len(*arr))
	for idx, el := range *arr {
		res[idx] = el.Get()
	}
	return res
}

func (arr *arrayGenerator) add(g Generator) Generator {
	*arr = append(*arr, g)
	return arr
}
