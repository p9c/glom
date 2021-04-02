package meta

type (
	// Data is the information about the opt to be used by interface code and other presentations of the data
	Data struct {
		Option        string
		Aliases       []string
		Group         string
		Label         string
		Description   string
		Documentation string
		Type          string
		Widget        string
		Options       []string
		OmitEmpty     bool
		Name          string
	}
)

func (m Data) GetAllOptionStrings() (opts []string) {
	opts = append([]string{m.Option}, m.Aliases...)
	return opts
}
