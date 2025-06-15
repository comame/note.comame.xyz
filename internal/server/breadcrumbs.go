package server

type breadcrumb struct {
	Label    string
	Location string
}

func joinBreadcrumbs(breadcrumbs ...breadcrumb) []breadcrumb {
	result := []breadcrumb{{Label: "Top", Location: "/"}}
	result = append(result, breadcrumbs...)
	return result
}

func postPageBreadcrumb(p *post) []breadcrumb {
	return []breadcrumb{
		{Label: p.Title, Location: p.getURL()},
	}
}
